# Frontend build plan

Rollout plan for the React frontend of the Go rewrite. Architecture decisions
live in [`docs/adr/0008-frontend-ssr-with-tanstack-start.md`](./adr/0008-frontend-ssr-with-tanstack-start.md).
Stack: **TanStack Start** (SSR) + Router + Query + Table + Form, Mantine UI,
i18next, zustand, hey-api generated client.

## Principles

- **Contract is the gate.** A page can be built as soon as its operation exists
  in TypeSpec (`api-spec/main.tsp`); the Go handler can land in parallel — the
  page develops against a prism mock of the OpenAPI in the meantime.
- **Generation-first.** Types, Query hooks, and Zod validators come from
  hey-api; admin CRUD is generated from the schema. Hand-write only what the
  generators can't (custom screens, the lesson player).
- **Legacy URLs preserved** (ADR-0002): keep `/languages`, locale prefixes, etc.

## Waves

Order is deliberate: **admin first** (large but self-contained, no SEO, pure
generation), then the value spine, then the public/SEO surface last.

### Wave 0 — SSR foundation
- Migrate the current Vite SPA to TanStack Start; wire SSR + Query dehydration.
- Locale routing: optional prefix (`en` bare, `ru`/`es` prefixed) + root
  auto-detect (session → Accept-Language → country-by-IP).
- Generic Start `loader` adapter over hey-api `Options` (`ensureQueryData`).
- httpOnly-cookie auth plumbing (stub until the auth backend lands).
- Run the existing catalog/course slice through SSR as the proof.

### Wave 1 — admin contract + CRUD engine
- Author admin resources in TypeSpec **regularly**: `list/get/create/update/
  delete` with consistent models per resource.
- Enable the hey-api **zod plugin**; regenerate the client.
- Build the engine: `<CrudList>` (TanStack Table) + `<CrudForm>` (TanStack Form
  + generated Zod + Mantine field wrappers) + admin route guard in `beforeLoad`.

### Wave 2 — admin screens
- Drive the ~40 CRUD screens through the engine: banners, blog_posts,
  language_categories, language_landing_pages, languages/lessons,
  lesson_members, reviews, roles, staff_members, users, leads, messages.
- Hand-build the non-standard ones: admin home/dashboard, role_permissions.

### Wave 3 — auth
- login, signup, passwords/edit, remind_passwords, magic_links, phone_auth.
- Cookie-JWT end-to-end; guards enforced via `beforeLoad`.

### Wave 4 — lesson player (spike + phase)
- Its own contract design done early even though built here: run code / stream
  output / submit solution is the hardest part of the API.
- Editor, code-run (river Docker jobs), SSE output streaming (generated
  `serverSentEvents.gen.ts`), multi-tab, zustand store.

### Wave 5 — user area
- `my` dashboard, account profile, course progress.

### Wave 6 — marketing / public (SEO)
- home, language_categories, blog_posts, cases, books, reviews, `/pages`
  (about/authors/privacy/tos/cookie), sitemap, error pages.
- SSR/SEO exercised fully here: hreflang across locales, canonical URLs, meta.

## Adoption notes

- Add when needed: **TanStack Virtual** (long tables/lists), **TanStack Pacer**
  (search debounce).
- Not adopted: TanStack Store (zustand stays), TanStack DB.
