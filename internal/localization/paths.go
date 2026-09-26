package localization

import "golang.org/x/text/language"

// PathPrefix is the URL path segment a site page carries for locale: none for
// English, the default, and `/<locale>` for every other language — legacy
// AppHost.locale_for_url and the frontend's optional `{-$locale}` segment. It
// is the one place the backend knows that rule, for the absolute URLs it
// builds itself (the Yandex feed).
func PathPrefix(locale string) string {
	if locale == "" || locale == language.English.String() {
		return ""
	}
	return "/" + locale
}
