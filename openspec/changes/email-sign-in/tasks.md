# Tasks

## 1. Contract

- [x] 1.1 In `api-spec/auth.tsp`, make `updatePassword` return `User` (with the cookie header, as `createSession` does) or `NotFoundError`, and set `@minLength(6)` on `ResetPasswordInput.password` and `SignUpInput.password`; run `make gen` and verify `go build ./...` and `pnpm check` pass
- [x] 1.2 Add a harness test that signing up with a five-character password is rejected with 400 and creates no User; verify with `go test ./internal/handlers/ -run SignUp`

## 2. Mail foundation

- [x] 2.1 Add config for `EMAIL_TOKEN_SECRET` (required, production must differ from the dev default, like `JWT_SECRET`), Postbox endpoint/credentials and `MAIL_FROM`; verify `go test ./internal/config/`
- [x] 2.2 Add the mailer package: `Mailer` port, Postbox adapter on `aws-sdk-go-v2/service/sesv2` (`go get`), `slog` adapter, in-memory recorder; wire the adapter choice in `internal/di`; verify `go build ./...` and that the worker starts with no Postbox credentials using the log adapter
- [x] 2.3 Add email-token issue/verify on `golang-jwt/jwt/v5` with purpose, user id, fifteen-minute expiry and fingerprint (email hash / password-digest hash); verified through the flow tests in groups 3 and 4
- [x] 2.4 Add ru/es subject and body strings for both emails to `internal/localization/locales` and the `html/template` templates, porting the legacy `ru.mailers.yml` wording; account emails fall back to `ru` for any other locale; verify the go-i18n bundle loads in `go test ./internal/localization/`
- [x] 2.5 Add the `AccountEmailArgs{Kind, UserID, Locale}` River job (unique by kind+user for one minute, locale excluded) and its worker that loads the User, issues the token, renders and sends; register it in the worker graph; verify with a jobs test that the recorder receives a Russian email with a `/magic_links/<token>` link and that a second insert within the minute is skipped

## 3. Magic Link

- [ ] 3.1 Extract the shared `signIn` helper from `CreateSession` (event, cookies, guest merge) and switch `CreateSession` to it; verify the existing session and guest-merge tests still pass
- [ ] 3.2 Implement `CreateMagicLink` (normalised lookup, enqueue for known users, always 204) and `ConsumeMagicLink` (verify token, `signIn`, 404 otherwise); add harness tests for every `magic-link` spec scenario: registered/unknown/case-insensitive email, guest merge on consume, repeat follow, expired, email changed, forged/legacy token; verify `go test ./internal/handlers/ -run MagicLink`

## 4. Password Reset

- [ ] 4.1 Implement `CreatePasswordReminder` (always 204), `CheckPasswordResetToken` and `UpdatePassword` (verify token, set bcrypt password, `signIn`, 404 on a refused token); add harness tests for every `password-reset` spec scenario: unknown email, valid/forged check, successful reset then old-password login fails, reused link, expired link, short password leaves the link valid; verify `go test ./internal/handlers/ -run Password`

## 5. Frontend

- [ ] 5.1 Add `magic_links/new` and `magic_links/$token` routes under `src/routes/{-$locale}/` with the generated Query hooks and ru/es strings (success → home, 404 → `new` with a message); add Vitest browser tests with MSW for both outcomes; verify `pnpm check` and `pnpm test`
- [ ] 5.2 Add `remind_password/new` and `password/$token/edit` routes (loader checks the token, 404 → `remind_password/new` with a message; success → signed in, home) with ru/es strings and Vitest tests for valid, refused and short-password cases; verify `pnpm check` and `pnpm test`
- [ ] 5.3 Link the sign-in form to `magic_links/new` and `remind_password/new`; verify both links in the `session/new` test

## 6. Integration and docs

- [ ] 6.1 Run both flows end to end with `make dev` in `ru` and `es`: request, copy the link from the worker log, follow it, and confirm you are signed in and guest progress is kept
- [ ] 6.2 Mark blocker #3 done in `docs/PARITY.md` and record the new env vars in the deploy notes; run `make lint` and `make test` clean
