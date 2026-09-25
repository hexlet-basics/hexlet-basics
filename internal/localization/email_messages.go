package localization

import "github.com/nicksnyder/go-i18n/v2/i18n"

// Account email copy, ported from legacy config/locales/*.mailers.yml. Each
// email has the same shape — why it was sent, a button, the bare link for
// clients that drop buttons, and what to do if it was not asked for — so the
// template stays one and only these strings differ.
var (
	// EmailLinkInvalid answers every refused Magic Link or Password Reset link
	// alike: expired, forged and outlived links are not told apart.
	EmailLinkInvalid = Message{value: i18n.Message{
		ID:    "emails.link_invalid",
		Other: "This link is no longer valid. Request a new one.",
	}}

	MagicLinkSubject = Message{value: i18n.Message{
		ID:    "emails.magic_link.subject",
		Other: "Your sign-in link",
	}}
	MagicLinkIntro = Message{value: i18n.Message{
		ID:    "emails.magic_link.intro",
		Other: "You requested a sign-in link for code-basics.com. Follow this link or use the button below to sign in.",
	}}
	MagicLinkAction = Message{value: i18n.Message{
		ID:    "emails.magic_link.action",
		Other: "Sign in",
	}}
	MagicLinkLinkLabel = Message{value: i18n.Message{
		ID:    "emails.magic_link.link_label",
		Other: "Sign-in link for code-basics.com",
	}}
	MagicLinkIgnore = Message{value: i18n.Message{
		ID:    "emails.magic_link.ignore",
		Other: "If you did not request this, just ignore this email.",
	}}

	PasswordResetSubject = Message{value: i18n.Message{
		ID:    "emails.password_reset.subject",
		Other: "Reset password",
	}}
	PasswordResetIntro = Message{value: i18n.Message{
		ID:    "emails.password_reset.intro",
		Other: "You requested a password reset link for code-basics.com. Follow this link or use the button below to set a new password.",
	}}
	PasswordResetAction = Message{value: i18n.Message{
		ID:    "emails.password_reset.action",
		Other: "Change password",
	}}
	PasswordResetLinkLabel = Message{value: i18n.Message{
		ID:    "emails.password_reset.link_label",
		Other: "Password reset link for code-basics.com",
	}}
	PasswordResetIgnore = Message{value: i18n.Message{
		ID:    "emails.password_reset.ignore",
		Other: "If you did not request a password reset, ignore this email. Your password will stay the same until you open the link and choose a new one.",
	}}
)
