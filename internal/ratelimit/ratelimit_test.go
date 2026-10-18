package ratelimit

import (
	"strconv"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAKeyHasItsBudgetBackWhenTheWindowPasses(t *testing.T) {
	t.Parallel()

	now := time.Now()
	limiter := New(time.Minute, func() time.Time { return now })

	first, _ := limiter.Allow(Charge{"a", 1})
	second, wait := limiter.Allow(Charge{"a", 1})
	other, _ := limiter.Allow(Charge{"b", 1})
	now = now.Add(time.Minute)
	later, _ := limiter.Allow(Charge{"a", 1})

	assert.True(t, first)
	assert.False(t, second)
	assert.Equal(t, time.Minute, wait)
	assert.True(t, other, "keys do not share a budget")
	assert.True(t, later)
}

func TestARefusedAttemptSpendsNoBudgetOnTheOtherKeys(t *testing.T) {
	t.Parallel()

	now := time.Now()
	limiter := New(time.Hour, func() time.Time { return now })
	ok, _ := limiter.Allow(Charge{"client", 1})
	require.True(t, ok)

	refused, wait := limiter.Allow(Charge{"user", 2}, Charge{"client", 1})
	first, _ := limiter.Allow(Charge{"user", 2})
	second, _ := limiter.Allow(Charge{"user", 2})
	third, _ := limiter.Allow(Charge{"user", 2})

	assert.False(t, refused)
	assert.Equal(t, time.Hour, wait)
	assert.True(t, first)
	assert.True(t, second, "the refused attempt did not spend the user's budget")
	assert.False(t, third)
}

func TestTheLimiterStopsTrackingNewKeysAtItsCapUntilWindowsExpire(t *testing.T) {
	t.Parallel()

	now := time.Now()

	limiter := New(time.Hour, func() time.Time { return now })
	for i := range maxTracked {
		ok, _ := limiter.Allow(Charge{strconv.Itoa(i), 1})
		require.True(t, ok)
	}

	full, _ := limiter.Allow(Charge{"new", 1})
	now = now.Add(time.Hour + time.Second)
	later, _ := limiter.Allow(Charge{"new", 1})

	assert.False(t, full)
	assert.True(t, later, "expired windows are pruned to make room")
	assert.Len(t, limiter.windows, 1)
}

// The refusal's wait is the real time left in the window: with the wall clock
// inside a synctest bubble, sleeping that long is what frees the budget, and
// sleeping a moment less is not enough.
func TestTheRealClockFreesTheBudgetAfterTheWaitItReported(t *testing.T) {
	t.Parallel()

	for name, early := range map[string]time.Duration{
		"just before the wait": time.Nanosecond,
		"exactly the wait":     0,
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			synctest.Test(t, func(t *testing.T) {
				limiter := New(time.Minute, time.Now)
				first, _ := limiter.Allow(Charge{"a", 1})
				refused, wait := limiter.Allow(Charge{"a", 1})

				require.True(t, first)
				require.False(t, refused)

				time.Sleep(wait - early)

				again, _ := limiter.Allow(Charge{"a", 1})

				assert.Equal(t, early == 0, again)
			})
		})
	}
}
