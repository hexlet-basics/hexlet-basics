// Package accountemail sends the emails that let a User back into their
// account: a Magic Link and a Password Reset. The HTTP process only enqueues the
// intent (Enqueuer); the worker renders and delivers it (Sender).
package accountemail

import (
	"bytes"
	"context"
	"embed"
	"html/template"
	"net/url"

	"github.com/samber/oops"

	"hexletbasics/ent"
	"hexletbasics/internal/emailtoken"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/mailer"
)

//go:embed templates/email.html
var templateFiles embed.FS

func emailTemplate() *template.Template {
	return template.Must(template.ParseFS(templateFiles, "templates/email.html"))
}

// Locale picks the language of an account email from the request locale. Only
// ru and es are emailed; anything else — including en, which the rewrite is
// dropping — gets ru.
func Locale(requestLocale string) string {
	if requestLocale == "es" {
		return "es"
	}

	return "ru"
}

// Sender renders and delivers one account email in the worker.
type Sender struct {
	db      *ent.Client
	tokens  *emailtoken.Tokens
	mail    mailer.Mailer
	i18n    *localization.Translator
	siteURL string
}

// NewSender builds the worker-side sender. siteURL is the origin the emailed
// links open (config.SiteURL).
func NewSender(
	db *ent.Client,
	tokens *emailtoken.Tokens,
	mail mailer.Mailer,
	translator *localization.Translator,
	siteURL string,
) *Sender {
	return &Sender{db: db, tokens: tokens, mail: mail, i18n: translator, siteURL: siteURL}
}

type copyText struct {
	subject, intro, action, linkLabel, ignore localization.Message
}

func copies() map[emailtoken.Purpose]copyText {
	return map[emailtoken.Purpose]copyText{
		emailtoken.MagicLink: {
			subject:   localization.MagicLinkSubject,
			intro:     localization.MagicLinkIntro,
			action:    localization.MagicLinkAction,
			linkLabel: localization.MagicLinkLinkLabel,
			ignore:    localization.MagicLinkIgnore,
		},
		emailtoken.PasswordReset: {
			subject:   localization.PasswordResetSubject,
			intro:     localization.PasswordResetIntro,
			action:    localization.PasswordResetAction,
			linkLabel: localization.PasswordResetLinkLabel,
			ignore:    localization.PasswordResetIgnore,
		},
	}
}

// Send issues a fresh token and mails the link. A User deleted, or left without
// an email, between the request and the job has nobody to mail, so the job
// completes instead of retrying.
func (s *Sender) Send(ctx context.Context, purpose emailtoken.Purpose, userID int, locale string) error {
	text, ok := copies()[purpose]
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

	err = emailTemplate().Execute(&body, map[string]string{
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
func (s *Sender) link(purpose emailtoken.Purpose, token, locale string) (string, error) {
	segments := []string{locale, "magic_links", token}
	if purpose == emailtoken.PasswordReset {
		segments = []string{locale, "password", token, "edit"}
	}

	link, err := url.JoinPath(s.siteURL, segments...)
	if err != nil {
		return "", oops.Wrapf(err, "build %s link", purpose)
	}

	return link, nil
}
