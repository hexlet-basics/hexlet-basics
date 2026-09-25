// Package accountemails sends the emails that let a User back into their
// account: a Magic Link and a Password Reset. The HTTP process only enqueues the
// intent (Enqueuer); the worker renders and delivers it (Sender).
package accountemails

import (
	"bytes"
	"context"
	"database/sql"
	"embed"
	"html/template"
	"net/url"

	"github.com/riverqueue/river"
	"github.com/samber/oops"

	"hexletbasics/ent"
	"hexletbasics/internal/emailtokens"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/mailer"
)

//go:embed templates/email.html
var templateFiles embed.FS

var emailTemplate = template.Must(template.ParseFS(templateFiles, "templates/email.html"))

// Locale picks the language of an account email from the request locale. Only
// ru and es are emailed; anything else — including en, which the rewrite is
// dropping — gets ru.
func Locale(requestLocale string) string {
	if requestLocale == "es" {
		return "es"
	}
	return "ru"
}

// Enqueuer schedules account emails from the HTTP process through the
// insert-only River client.
type Enqueuer struct {
	river *river.Client[*sql.Tx]
}

// NewEnqueuer wires the shared insert-only River client.
func NewEnqueuer(riverClient *river.Client[*sql.Tx]) *Enqueuer {
	return &Enqueuer{river: riverClient}
}

// EnqueueAccountEmail inserts the job. A duplicate within the cooldown is
// skipped by River and is not an error: the caller answers the same either way.
func (e *Enqueuer) EnqueueAccountEmail(ctx context.Context, purpose emailtokens.Purpose, userID int, locale string) error {
	_, err := e.river.Insert(ctx, jobs.AccountEmailArgs{Purpose: purpose, UserID: userID, Locale: locale}, nil)
	if err != nil {
		return oops.Wrapf(err, "enqueue %s email", purpose)
	}
	return nil
}

// Sender renders and delivers one account email in the worker.
type Sender struct {
	db      *ent.Client
	tokens  *emailtokens.Tokens
	mail    mailer.Mailer
	i18n    *localization.Translator
	siteURL string
}

// NewSender builds the worker-side sender. siteURL is the origin the emailed
// links open (config.SiteURL).
func NewSender(
	db *ent.Client,
	tokens *emailtokens.Tokens,
	mail mailer.Mailer,
	translator *localization.Translator,
	siteURL string,
) *Sender {
	return &Sender{db: db, tokens: tokens, mail: mail, i18n: translator, siteURL: siteURL}
}

type copyText struct {
	subject, intro, action, linkLabel, ignore localization.Message
}

var copies = map[emailtokens.Purpose]copyText{
	emailtokens.MagicLink: {
		subject:   localization.MagicLinkSubject,
		intro:     localization.MagicLinkIntro,
		action:    localization.MagicLinkAction,
		linkLabel: localization.MagicLinkLinkLabel,
		ignore:    localization.MagicLinkIgnore,
	},
	emailtokens.PasswordReset: {
		subject:   localization.PasswordResetSubject,
		intro:     localization.PasswordResetIntro,
		action:    localization.PasswordResetAction,
		linkLabel: localization.PasswordResetLinkLabel,
		ignore:    localization.PasswordResetIgnore,
	},
}

// Send issues a fresh token and mails the link. A User deleted, or left without
// an email, between the request and the job has nobody to mail, so the job
// completes instead of retrying.
func (s *Sender) Send(ctx context.Context, purpose emailtokens.Purpose, userID int, locale string) error {
	text, ok := copies[purpose]
	if !ok {
		return oops.Errorf("unknown account email purpose %q", purpose)
	}
	u, err := s.db.User.Get(ctx, userID)
	if ent.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return oops.Wrapf(err, "load user for %s email", purpose)
	}
	if u.Email == nil {
		return nil
	}
	token, err := s.tokens.Issue(purpose, u)
	if err != nil {
		return err
	}

	link, err := s.link(purpose, token, locale)
	if err != nil {
		return err
	}

	subject := s.i18n.TextIn(locale, text.subject)
	var body bytes.Buffer
	err = emailTemplate.Execute(&body, map[string]string{
		"Locale":    locale,
		"Subject":   subject,
		"Intro":     s.i18n.TextIn(locale, text.intro),
		"Action":    s.i18n.TextIn(locale, text.action),
		"LinkLabel": s.i18n.TextIn(locale, text.linkLabel),
		"Ignore":    s.i18n.TextIn(locale, text.ignore),
		"URL":       link,
	})
	if err != nil {
		return oops.Wrapf(err, "render %s email", purpose)
	}
	return s.mail.Send(ctx, mailer.Message{To: *u.Email, Subject: subject, HTML: body.String()})
}

// link builds the legacy URL the email opens, under the email's locale prefix
// (both ru and es are prefixed on the site).
func (s *Sender) link(purpose emailtokens.Purpose, token, locale string) (string, error) {
	segments := []string{locale, "magic_links", token}
	if purpose == emailtokens.PasswordReset {
		segments = []string{locale, "password", token, "edit"}
	}
	link, err := url.JoinPath(s.siteURL, segments...)
	if err != nil {
		return "", oops.Wrapf(err, "build %s link", purpose)
	}
	return link, nil
}
