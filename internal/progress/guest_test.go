package progress_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/progress"
)

const testSecret = "test-secret"

func TestGuestCookieRoundTripsProgress(t *testing.T) {
	codec := progress.NewGuestCodec(testSecret)
	stored := progress.GuestProgress{}.
		Record("javascript", "variables").
		Record("python", "hello-world")

	value, err := codec.Encode(stored)
	require.NoError(t, err)

	decoded, err := codec.Decode(value)
	require.NoError(t, err)

	lesson, ok := decoded.Furthest("javascript")
	require.True(t, ok)
	assert.Equal(t, "variables", lesson)

	_, ok = decoded.Furthest("ruby")
	assert.False(t, ok, "a course the visitor never touched has no entry")
}

// Recording a course again replaces its entry rather than accumulating one, and
// moves it to the most-recently-touched end.
func TestGuestProgressKeepsOneEntryPerCourse(t *testing.T) {
	stored := progress.GuestProgress{}.
		Record("javascript", "hello-world").
		Record("python", "intro").
		Record("javascript", "variables")

	require.Len(t, stored.Entries, 2)
	assert.Equal(t, "python", stored.Entries[0].CourseSlug, "least recently touched first")
	assert.Equal(t, "javascript", stored.Entries[1].CourseSlug)

	lesson, _ := stored.Furthest("javascript")
	assert.Equal(t, "variables", lesson)
}

// A tampered payload is rejected: this state becomes rows and events on sign-up,
// so an unsigned edit would be a way to mint completions.
func TestGuestCookieRejectsATamperedPayload(t *testing.T) {
	codec := progress.NewGuestCodec(testSecret)
	value, err := codec.Encode(progress.GuestProgress{}.Record("javascript", "hello-world"))
	require.NoError(t, err)

	payload, signature, _ := strings.Cut(value, ".")
	forged := payload[:len(payload)-4] + "AAAA" + "." + signature

	_, err = codec.Decode(forged)
	assert.ErrorIs(t, err, progress.ErrGuestCookieInvalid)
}

// A cookie signed with another secret is not ours.
func TestGuestCookieRejectsAForeignSignature(t *testing.T) {
	value, err := progress.NewGuestCodec("someone-elses-secret").
		Encode(progress.GuestProgress{}.Record("javascript", "hello-world"))
	require.NoError(t, err)

	_, err = progress.NewGuestCodec(testSecret).Decode(value)
	assert.ErrorIs(t, err, progress.ErrGuestCookieInvalid)
}

func TestGuestCookieRejectsGarbage(t *testing.T) {
	codec := progress.NewGuestCodec(testSecret)

	for _, raw := range []string{"", "not-a-cookie", "!!!.???", "eyJ9."} {
		_, err := codec.Decode(raw)
		assert.ErrorIsf(t, err, progress.ErrGuestCookieInvalid, "value %q", raw)
	}
}

// Past the size limit the least recently touched courses go first, so the
// visitor loses the courses they stopped using rather than the whole cookie.
func TestGuestCookieEvictsLeastRecentlyTouchedCourses(t *testing.T) {
	codec := progress.NewGuestCodec(testSecret)

	stored := progress.GuestProgress{}
	for i := range 200 {
		stored = stored.Record(
			fmt.Sprintf("course-with-a-fairly-long-slug-%03d", i),
			fmt.Sprintf("lesson-with-a-fairly-long-slug-%03d", i),
		)
	}

	value, err := codec.Encode(stored)
	require.NoError(t, err)
	assert.LessOrEqual(t, len(value), 3500, "stays under the browser's cookie limit")

	decoded, err := codec.Decode(value)
	require.NoError(t, err)
	assert.Less(t, len(decoded.Entries), 200, "something was evicted")

	_, ok := decoded.Furthest("course-with-a-fairly-long-slug-199")
	assert.True(t, ok, "the most recently touched course survives")
	_, ok = decoded.Furthest("course-with-a-fairly-long-slug-000")
	assert.False(t, ok, "the least recently touched one is dropped first")
}

func TestGuestCookieIsHttpOnlyAndLongLived(t *testing.T) {
	codec := progress.NewGuestCodec(testSecret)

	cookie, err := codec.Cookie(progress.GuestProgress{}.Record("javascript", "hello-world"), true)
	require.NoError(t, err)

	assert.Equal(t, progress.GuestCookieName, cookie.Name)
	assert.True(t, cookie.HttpOnly, "the client never reads guest progress")
	assert.True(t, cookie.Secure)
	assert.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
	assert.Equal(t, "/", cookie.Path)
	assert.Equal(t, 365*24*60*60, cookie.MaxAge, "a year")
}
