// Package mailer delivers transactional email (ADR-0006). Callers build a
// finished Message; which adapter carries it — Postbox in production, the log in
// development, memory in tests — is chosen once in DI, so no flow ever branches
// on the environment.
package mailer

import (
	"context"
	"log/slog"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/samber/oops"

	"hexletbasics/internal/config"
)

// Message is one rendered email. It carries HTML only: the text/plain part is a
// post-parity improvement (ADR-0006), as it was absent in legacy too.
type Message struct {
	To      string
	Subject string
	HTML    string
}

// Mailer delivers a Message. An error means the provider did not accept it, so
// the caller (a River job) retries.
type Mailer interface {
	Send(ctx context.Context, msg Message) error
}

// New picks the adapter the configuration asks for: Postbox when its access key
// is set, the log otherwise. Production refuses to start without the key
// (config.validateProduction), so the log adapter can never swallow real mail.
func New(cfg config.MailConfig, logger *slog.Logger) Mailer {
	if cfg.PostboxEnabled() {
		return NewPostbox(cfg)
	}
	return NewLog(logger)
}

// Postbox sends through Yandex Cloud Postbox's SES v2-compatible API. Only the
// endpoint and credentials differ from AWS SES; the SDK does the signing.
type Postbox struct {
	client *sesv2.Client
	from   string
}

// NewPostbox builds the SES v2 client pointed at Postbox.
func NewPostbox(cfg config.MailConfig) *Postbox {
	client := sesv2.New(sesv2.Options{
		Region:       cfg.PostboxRegion,
		BaseEndpoint: aws.String(cfg.PostboxEndpoint),
		Credentials: credentials.NewStaticCredentialsProvider(
			cfg.PostboxAccessKeyID, cfg.PostboxSecretAccessKey, "",
		),
	})
	return &Postbox{client: client, from: cfg.From}
}

// Send submits one simple (subject + HTML body) message.
func (p *Postbox) Send(ctx context.Context, msg Message) error {
	_, err := p.client.SendEmail(ctx, &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(p.from),
		Destination:      &types.Destination{ToAddresses: []string{msg.To}},
		Content: &types.EmailContent{Simple: &types.Message{
			Subject: &types.Content{Data: aws.String(msg.Subject), Charset: aws.String("UTF-8")},
			Body: &types.Body{Html: &types.Content{
				Data: aws.String(msg.HTML), Charset: aws.String("UTF-8"),
			}},
		}},
	})
	if err != nil {
		return oops.Wrapf(err, "send email via postbox")
	}
	return nil
}

// Log writes each message to the log instead of sending it, so a developer can
// follow the link in it without any mail setup.
type Log struct {
	logger *slog.Logger
}

// NewLog builds the development adapter.
func NewLog(logger *slog.Logger) *Log {
	return &Log{logger: logger}
}

// Send logs the message and never fails.
func (l *Log) Send(ctx context.Context, msg Message) error {
	l.logger.InfoContext(ctx, "email not sent: no mail provider configured",
		"to", msg.To, "subject", msg.Subject, "html", msg.HTML)
	return nil
}

// Recorder keeps sent messages in memory for tests to assert on. It is safe
// for concurrent use because River works jobs on several goroutines.
type Recorder struct {
	mu   sync.Mutex
	sent []Message
}

// Send records the message.
func (r *Recorder) Send(_ context.Context, msg Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sent = append(r.sent, msg)
	return nil
}

// Sent returns a copy of every recorded message, oldest first.
func (r *Recorder) Sent() []Message {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]Message(nil), r.sent...)
}
