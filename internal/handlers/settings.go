package handlers

// Settings are the few configuration values the handlers read, handed over as
// plain values: the handlers package depends on what it uses, not on the
// process's whole configuration.
type Settings struct {
	// PublicURL is the origin this API is reached at; self-served asset URLs are
	// built on it, since there is no *http.Request at the ogen handler boundary.
	PublicURL string
	// SiteURL is the origin of the public site that redirects and emailed links
	// open.
	SiteURL string
	// SecureCookies marks cookies Secure; true when the site is served over https.
	SecureCookies bool
	// GuestCookieSecret signs the visitor-progress cookie.
	GuestCookieSecret string
	// EmailTokenSecret verifies Magic Link and Password Reset tokens.
	EmailTokenSecret string
	// CourseRepoBaseURL is the GitHub org base of the exercises repositories.
	CourseRepoBaseURL string
	// BookBlobKey is the bucket key of the free book's PDF.
	BookBlobKey string
}
