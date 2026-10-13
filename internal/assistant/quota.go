package assistant

import (
	"context"
	"errors"
	"time"

	"github.com/samber/oops"

	"hexletbasics/ent"
	"hexletbasics/ent/aimessage"
)

// DailyQuota is how many questions a learner may ask per UTC day. The panel
// promises seven; legacy's policy compared `count <= 7` before storing and so
// actually admitted an eighth, which is the one place this port does not copy
// it.
const DailyQuota = 7

// ErrQuotaExceeded refuses a question once today's quota is used up.
var ErrQuotaExceeded = errors.New("assistant daily quota exceeded")

// Quota decides whether a learner may ask another question today. The history
// read (to show the panel as used up) and the ask (to refuse) put the same
// question, so both go through it.
//
// The check is not atomic with storing the question: two questions started at
// once at the limit can both pass, as they could on legacy, and the cost of one
// extra answer does not justify a lock.
type Quota struct {
	db  *ent.Client
	now func() time.Time
}

// NewQuota builds the quota over an ent client.
func NewQuota(db *ent.Client) *Quota {
	return &Quota{db: db, now: time.Now}
}

// Exhausted reports whether the learner has used up today's questions.
//
// Today is counted from the start of the UTC day — legacy's Date.current under
// the default UTC time zone. The columns are `timestamp` without a zone,
// holding UTC wall-clock time as Rails wrote it, which is why the rows are
// stamped in UTC explicitly: the mixin's time.Now default would store the
// process's local wall clock, and a question asked late in the evening west of
// Greenwich would land on yesterday.
func (q *Quota) Exhausted(ctx context.Context, userID int) (bool, error) {
	dayStart := q.now().UTC().Truncate(24 * time.Hour)
	asked, err := q.db.AiMessage.Query().
		Where(
			aimessage.UserID(userID),
			aimessage.RoleEQ(RoleUser),
			aimessage.CreatedAtGTE(dayStart),
		).
		Count(ctx)
	if err != nil {
		return false, oops.Wrapf(err, "count today's assistant questions")
	}
	return asked >= DailyQuota, nil
}
