# ADR-0015: Cutover scope, English locale kept until after it, and the `/api` prefix

**Status:** Accepted

## Context

ADR-0002 commits to a hard cutover at parity, but "parity" was never pinned
down against what production actually does. A pass over the legacy app turned
up surfaces that exist in code and not in use, and one that exists in the UI
and has never worked. Porting them all would be parity with the codebase, not
with the product, and every extra surface is one more possible cause of a bad
cutover night.

Separately, the contract's paths start at the root (`/account/profile`) and
collide with the frontend's page paths, and the browser reaches the API
through a `VITE_API_URL` that production has never had a value for.

## Decision

**Done means deployable.** The Go stack is ready for cutover when everything
production users rely on works end to end and the stack deploys to the
existing cluster; not when the contract is implemented.

**Parity is with production, not with the legacy codebase.** These are not
ported, and their operations leave the contract:

- **Phone sign-in** (`createPhoneAuth`, `confirmPhoneAuth`) — the button is
  live but every environment, production included, wires `SmsSenderStub`, so
  no SMS has ever been sent.
- **Passkeys** (the six passkey operations) — never shipped to production
  users. They come back after cutover as a new feature on go-webauthn
  (ADR-0003), with no compatibility burden.
- **Google one-tap / OAuth** — routes already commented out, no UI entry.
- **Surveys** — models and tables with no routes, controllers or UI.
- **Flipper** — installed, no flag is read anywhere.
- **Ahoy** — its only product use is lead attribution; the frontend keeps the
  first visit's UTM parameters and landing page in a cookie and sends them
  with `createLead`, and the IP comes from the request.
- **`getPage`** — the static pages (about, authors, privacy, tos, cookie
  policy) are hardcoded markup with no data; they become frontend routes.

Their tables stay in the database until the post-cutover cleanup.

**English stays until after cutover.** The rewrite serves `en` exactly as
legacy does — unprefixed URLs, `ru` and `es` prefixed — so cutover changes no
URL and no stored locale. Dropping `en` (moving `users.locale` to `ru`,
redirecting the unprefixed URLs, retiring English content) is a separate change
made after the rollback window.

**Two-week rollback window.** For two weeks after cutover the legacy image
stays deployable and the Go stack runs only additive migrations, so rolling
back is pointing the ingress at legacy again. Everything destructive — the
`en` removal, the dead tables above — waits for the window to close.

**Everything Go serves lives under `/api` on the site's own host.** The
ingress routes `/api/*` to Go and the rest to the SSR server. One origin means
no CORS and no parent-domain cookies: the httpOnly JWT and the XSRF cookie
(ADR-0003) work as they do in development. The legacy Yandex feed URL,
`/api/feeds/yandex_courses`, already sits under the prefix and keeps its
address. The GitHub webhook keeps `/webhooks/github` for now through an
ingress rule.

## Considered options

- **An `api.` subdomain** — rejected: it buys nothing a path prefix does not,
  and costs CORS plus cookies scoped to the parent domain.
- **Dropping `en` at cutover** — rejected: it would change URLs and user data
  in the same moment as the stack, so a traffic drop would have two possible
  causes.
