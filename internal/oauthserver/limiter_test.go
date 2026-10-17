package oauthserver

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTheLimiterCountsAnIPv6PrefixAsOneAddress(t *testing.T) {
	now := time.Now()
	l := newLimiter(2, time.Hour, func() time.Time { return now })

	ok1, _ := l.allow("2001:db8:1:2::1")
	ok2, _ := l.allow("2001:db8:1:2:ffff::2")
	ok3, wait := l.allow("2001:db8:1:2:abcd::3")
	other, _ := l.allow("2001:db8:1:3::1")

	assert.True(t, ok1)
	assert.True(t, ok2)
	assert.False(t, ok3, "one machine must not mint a budget per address inside its /64")
	assert.Positive(t, wait)
	assert.True(t, other)
}

func TestTheLimiterSharesOneBucketAmongUnresolvedAddresses(t *testing.T) {
	now := time.Now()
	l := newLimiter(1, time.Hour, func() time.Time { return now })

	first, _ := l.allow("")
	second, _ := l.allow("not an address")

	assert.True(t, first)
	assert.False(t, second)
}

func TestTheLimiterStopsTrackingNewAddressesAtItsCapUntilWindowsExpire(t *testing.T) {
	now := time.Now()
	l := newLimiter(1, time.Hour, func() time.Time { return now })
	for i := range maxTrackedAddresses {
		ok, _ := l.allow(fmt.Sprintf("10.%d.%d.%d", i>>16, (i>>8)&255, i&255))
		require.True(t, ok)
	}

	full, _ := l.allow("203.0.113.1")
	now = now.Add(time.Hour + time.Second)
	later, _ := l.allow("203.0.113.1")

	assert.False(t, full)
	assert.True(t, later, "expired windows are pruned to make room")
	assert.Len(t, l.windows, 1)
}
