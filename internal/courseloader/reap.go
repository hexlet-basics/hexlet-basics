package courseloader

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/oops"

	"hexletbasics/ent"
	"hexletbasics/ent/courseversion"
)

// StuckBuildAfter is how long a version may sit in `building` before the reaper
// declares its worker dead. It mirrors legacy Language::Version::STUCK_BUILD_AFTER:
// a real build finishes in minutes, so two hours without a state change means
// the process was hard-killed (OOM/SIGKILL) and never reached Loader.fail.
const StuckBuildAfter = 2 * time.Hour

// ReapStuck marks every version stuck in `building` past StuckBuildAfter as
// `failed` with an explicit result, and returns the versions it reaped. It is
// the Go port of the legacy ReapStuckVersionBuildsJob: without it a killed
// build stays `building` forever with an empty result and no trace in the admin
// versions table. The previously-live version is never touched — only the
// state/result of the dead build change.
func (l *Loader) ReapStuck(ctx context.Context) ([]*ent.CourseVersion, error) {
	stuck, err := l.db.CourseVersion.Query().
		Where(
			courseversion.State(stateBuilding),
			// Inclusive bound, matching the legacy `updated_at: ..2.hours.ago` range.
			courseversion.UpdatedAtLTE(time.Now().Add(-StuckBuildAfter)),
		).
		All(ctx)
	if err != nil {
		return nil, oops.Wrapf(err, "query stuck course versions")
	}

	reaped := make([]*ent.CourseVersion, 0, len(stuck))
	for _, version := range stuck {
		// The result quotes the timestamp read above: the write below bumps
		// updated_at, so it must be captured first. The state/updated_at guard
		// skips a build that finished between the select and this update instead
		// of clobbering its outcome — the same predicate-as-lock shape as claim.
		// Rails' Time#iso8601 in the default UTC zone is RFC 3339 at second
		// precision with a `Z` suffix.
		updated, err := l.db.CourseVersion.UpdateOneID(version.ID).
			Where(
				courseversion.State(stateBuilding),
				courseversion.UpdatedAt(version.UpdatedAt),
			).
			SetState(stateFailed).
			SetResult(fmt.Sprintf(
				"Build reaped: stuck in 'building' since %s (worker likely killed, e.g. OOM)",
				version.UpdatedAt.UTC().Format(time.RFC3339),
			)).
			Save(ctx)
		if ent.IsNotFound(err) {
			continue
		}
		if err != nil {
			return reaped, oops.Wrapf(err, "reap course version %d", version.ID)
		}
		reaped = append(reaped, updated)
	}
	return reaped, nil
}
