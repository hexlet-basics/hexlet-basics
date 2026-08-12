package progress

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// GuestCookieName carries a visitor's progress. It is httpOnly: the client
// renders what the server derives from it and never reads it, so there is one
// place that interprets this state.
const GuestCookieName = "guest_progress"

// guestCookieMaxBytes keeps the encoded value clear of the 4096-byte browser
// limit. A visitor who exceeds it loses their least recently touched courses
// rather than the whole cookie — silently dropping everything would be the one
// failure mode a learner would notice.
const guestCookieMaxBytes = 3500

// guestCookieTTL matches the legacy guest window.
const guestCookieTTL = 365 * 24 * time.Hour

// ErrGuestCookieInvalid reports a cookie whose signature does not verify, or
// which is not decodable at all. Both are treated the same way: the visitor is
// handled as if they had no cookie.
var ErrGuestCookieInvalid = errors.New("progress: guest cookie is not valid")

// GuestProgress is a visitor's whole state: per Course, the slug of the
// furthest Lesson they finished.
//
// A slug rather than a position, because a position is meaningless across
// Versions — promoting a Version that inserts or removes Lessons renumbers
// everything, and a stored number would come to denote a different Lesson. The
// Position is resolved against the current Version at read time.
//
// One entry per Course is exact rather than lossy: sequential progression makes
// a guest's gaps unrepresentable, so the furthest Lesson implies the rest.
// Entries are ordered least-recently-touched first, which is the order they are
// evicted in.
type GuestProgress struct {
	Entries []GuestEntry `json:"entries"`
}

// GuestEntry is one Course's furthest finished Lesson.
type GuestEntry struct {
	CourseSlug string `json:"course"`
	LessonSlug string `json:"lesson"`
}

// Furthest reports the furthest Lesson the visitor finished in a Course.
func (g GuestProgress) Furthest(courseSlug string) (string, bool) {
	for _, entry := range g.Entries {
		if entry.CourseSlug == courseSlug {
			return entry.LessonSlug, true
		}
	}
	return "", false
}

// Record stores a Course's furthest finished Lesson and marks it most recently
// touched. The caller decides what "furthest" means — this type stores, it does
// not compare positions, because only the current Version can order them.
func (g GuestProgress) Record(courseSlug, lessonSlug string) GuestProgress {
	kept := make([]GuestEntry, 0, len(g.Entries)+1)
	for _, entry := range g.Entries {
		if entry.CourseSlug != courseSlug {
			kept = append(kept, entry)
		}
	}
	return GuestProgress{Entries: append(kept, GuestEntry{CourseSlug: courseSlug, LessonSlug: lessonSlug})}
}

// GuestCodec signs and verifies the cookie with the authentication secret.
//
// Signing is not optional: this state becomes database rows and domain events
// when the visitor signs up, so an unsigned cookie would let anyone mint
// fabricated completions on their own account.
type GuestCodec struct {
	secret []byte
}

// NewGuestCodec builds the codec around the existing authentication secret.
func NewGuestCodec(secret string) *GuestCodec {
	return &GuestCodec{secret: []byte(secret)}
}

// Encode renders the progress as a signed cookie value, evicting the least
// recently touched Courses until it fits.
func (c *GuestCodec) Encode(progress GuestProgress) (string, error) {
	for {
		value, err := c.sign(progress)
		if err != nil {
			return "", err
		}
		if len(value) <= guestCookieMaxBytes || len(progress.Entries) == 0 {
			return value, nil
		}
		progress = GuestProgress{Entries: progress.Entries[1:]}
	}
}

// Decode verifies the signature and returns the progress it carries.
func (c *GuestCodec) Decode(raw string) (GuestProgress, error) {
	payload, signature, found := strings.Cut(raw, ".")
	if !found {
		return GuestProgress{}, ErrGuestCookieInvalid
	}

	body, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return GuestProgress{}, ErrGuestCookieInvalid
	}
	got, err := base64.RawURLEncoding.DecodeString(signature)
	if err != nil {
		return GuestProgress{}, ErrGuestCookieInvalid
	}
	if !hmac.Equal(got, c.mac(body)) {
		return GuestProgress{}, ErrGuestCookieInvalid
	}

	var progress GuestProgress
	if err := json.Unmarshal(body, &progress); err != nil {
		return GuestProgress{}, ErrGuestCookieInvalid
	}
	return progress, nil
}

// Cookie is the response cookie carrying the progress: httpOnly so only the
// server interprets it, Lax so it survives the navigation back from an OAuth
// provider, and a year long.
func (c *GuestCodec) Cookie(progress GuestProgress, secure bool) (*http.Cookie, error) {
	value, err := c.Encode(progress)
	if err != nil {
		return nil, err
	}
	return &http.Cookie{
		Name:     GuestCookieName,
		Value:    value,
		Path:     "/",
		MaxAge:   int(guestCookieTTL.Seconds()),
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}, nil
}

func (c *GuestCodec) sign(progress GuestProgress) (string, error) {
	body, err := json.Marshal(progress)
	if err != nil {
		return "", fmt.Errorf("encode guest progress: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(body) + "." +
		base64.RawURLEncoding.EncodeToString(c.mac(body)), nil
}

func (c *GuestCodec) mac(body []byte) []byte {
	sum := hmac.New(sha256.New, c.secret)
	sum.Write(body)
	return sum.Sum(nil)
}
