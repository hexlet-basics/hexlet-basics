package handlers

import (
	"context"
	"net/http"
	"time"

	"hexletbasics/internal/api"
)

// localeCookie remembers the locale a visitor switched to. It stands in for
// legacy `session[:locale]`: the SSR root reads it to decide whether `/`
// redirects to a prefixed locale, so a Russian-speaking browser that picked
// English stays on English. Signed-in users also get the choice on their row.
const (
	localeCookie    = "locale"
	localeCookieTTL = 365 * 24 * time.Hour
)

// GetProfile returns the signed-in user's editable profile. The contract
// guarantees a session, so the user is always in the context here.
func (s *Server) GetProfile(ctx context.Context) (api.GetProfileRes, error) {
	u := MustAuthenticatedUser(ctx)
	profile := s.conv.ToUser(u)
	return &profile, nil
}

// UpdateProfile saves the first and last name. Their length and character
// rules are contract constraints, so the generated server has already refused
// anything legacy's validations would have. A null clears the column, as
// legacy's `update(struct.attributes)` did with a nil attribute.
func (s *Server) UpdateProfile(ctx context.Context, req *api.ProfileInput) (api.UpdateProfileRes, error) {
	u := MustAuthenticatedUser(ctx)

	update := u.Update()
	if value, ok := req.FirstName.Get(); ok {
		update.SetFirstName(string(value))
	} else {
		update.ClearFirstName()
	}
	if value, ok := req.LastName.Get(); ok {
		update.SetLastName(string(value))
	} else {
		update.ClearLastName()
	}
	u, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}
	profile := s.conv.ToUser(u)
	return &profile, nil
}

// DeleteAccount removes the signed-in user's account and signs them out. The
// row stays, in the removed state, as in legacy. Legacy ended only the current
// device's session; here every other device's session stops working too,
// because a session is only as good as the account it names and a removed
// account is refused whatever session it carries.
func (s *Server) DeleteAccount(ctx context.Context) (api.DeleteAccountRes, error) {
	u := MustAuthenticatedUser(ctx)
	if err := s.remover.Remove(ctx, u.ID); err != nil {
		return nil, lastAdminError(err)
	}
	if err := s.sessions.End(ctx); err != nil {
		return nil, err
	}
	return &api.DeleteAccountNoContent{}, nil
}

// SwitchLocale remembers the chosen locale (legacy LocalesController#switch).
// Legacy answered with a redirect to the same page in the new locale; the page
// knows its own URL, so here it navigates itself once the choice is stored.
// The contract's Locale enum already refuses anything but en, ru and es.
func (s *Server) SwitchLocale(ctx context.Context, params api.SwitchLocaleParams) (*api.SwitchLocaleNoContent, error) {
	locale := string(params.Locale)
	if u, ok := AuthenticatedUser(ctx); ok {
		if err := u.Update().SetLocale(locale).Exec(ctx); err != nil {
			return nil, err
		}
	}

	cookie := &http.Cookie{
		Name:     localeCookie,
		Value:    locale,
		Path:     "/",
		MaxAge:   int(localeCookieTTL.Seconds()),
		HttpOnly: true,
		Secure:   s.secure,
		SameSite: http.SameSiteLaxMode,
	}
	return &api.SwitchLocaleNoContent{SetCookie: []string{cookie.String()}}, nil
}
