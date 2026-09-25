package jobs

import (
	"context"
	"time"

	"github.com/riverqueue/river"

	"hexletbasics/internal/emailtokens"
)

// accountEmailCooldown is how long a second request for the same email to the
// same User is a no-op. Legacy kept the same one-minute cooldown in its cache.
const accountEmailCooldown = time.Minute

// AccountEmailSender renders and delivers one account email
// (accountemails.Sender in production).
type AccountEmailSender interface {
	Send(ctx context.Context, purpose emailtokens.Purpose, userID int, locale string) error
}

// AccountEmailArgs is the intent to email a User a Magic Link or a Password
// Reset. It deliberately carries no token or HTML: the worker issues the token
// when it sends, so a retried job never mails a link that already expired and
// the job table never holds a usable credential.
type AccountEmailArgs struct {
	Purpose emailtokens.Purpose `json:"purpose" river:"unique"`
	UserID  int                 `json:"user_id" river:"unique"`
	// Locale is left out of uniqueness: switching language between two clicks
	// is still the same request within the cooldown.
	Locale string `json:"locale"`
}

// Kind is River's stable job discriminator; do not rename once jobs are enqueued.
func (AccountEmailArgs) Kind() string { return "account_email" }

// InsertOpts makes the cooldown River uniqueness over (purpose, user) for one
// minute, and bounds retries: five attempts with River's backoff span about
// six minutes, past which a sign-in link is no longer what the visitor waits
// for, and the failure is already in Sentry.
func (AccountEmailArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		MaxAttempts: 5,
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: accountEmailCooldown,
		},
	}
}

type accountEmailWorker struct {
	river.WorkerDefaults[AccountEmailArgs]
	sender AccountEmailSender
}

func (w *accountEmailWorker) Work(ctx context.Context, job *river.Job[AccountEmailArgs]) error {
	return w.sender.Send(ctx, job.Args.Purpose, job.Args.UserID, job.Args.Locale)
}
