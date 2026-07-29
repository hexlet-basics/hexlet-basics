# Frontend: SSR via TanStack Start, generation-first data and forms

The frontend is a **server-rendered** React app built with **TanStack Start**
(SSR on a Node runtime) talking to the Go API. SSR is required for SEO parity
with the legacy Inertia app — Hexlet Basics is a public education/marketing
site whose catalog, course landings, blog, and `/pages` depend on it. The Go
server stays API-only (ADR-0001); a Node SSR process renders HTML and calls Go
server-side over the internal network, the browser hydrates and then talks to
Go directly.

The stack stays generation-first (ADR-0001): the OpenAPI contract drives the
TS types, TanStack Query hooks, and Zod validators (all via **hey-api**), so
data-fetching and form validation are generated, not hand-written. Admin CRUD
screens are generated from the same schema rather than authored per resource.

## Considered Options

- **TanStack Start** (chosen) — SSR framework built on the TanStack Router the
  app already uses; `router.tsx` and the hey-api client carry over, Query SSR
  dehydration is built in, stays on Vite. Least churn from the current stack.
- **React Router v7 / Remix** (rejected) — mature, but replaces TanStack Router
  and its loader model competes with TanStack Query; discards work already done.
- **Next.js** (rejected) — largest swap: abandons TanStack Router, RSC data
  model, non-Vite build. Too much off-stack change for "just SSR".
- **Manual Vite SSR** (rejected) — keeps the router but hand-rolls hydration,
  head, streaming, and Query dehydration; against the least-hand-written-code
  rule.
- **Build-time prerender / pure CSR** (rejected) — no second process, but not
  real SSR for dynamic/authenticated pages; SEO regression vs legacy.

## Decisions

- **Data:** hey-api generates `xxxOptions()`/mutations/query-keys. Start
  `loader` calls `queryClient.ensureQueryData(xxxOptions())` for SSR prefetch +
  dehydration via one generic adapter; components read the same options with
  `useQuery` — no double fetch, no hand-written loaders.
- **Forms:** TanStack Form (headless) + Mantine input wrappers; validators are
  Zod schemas generated from OpenAPI by the hey-api zod plugin. Chosen over
  `@mantine/form` because a schema-driven field layer is built regardless, which
  neutralizes Mantine-form's zero-wiring edge, and it keeps the stack on
  TanStack.
- **Admin CRUD:** generated from OpenAPI/JSON-Schema — TanStack Table columns
  from the resource model, form fields from the create/update schema. One CRUD
  engine (`list/new/edit/delete`) parameterized per resource; non-standard
  screens (dashboard, role_permissions) stay hand-built.
- **Auth:** JWT in an httpOnly cookie (not localStorage) so the SSR server can
  render authenticated pages; route guards via TanStack Router `beforeLoad`.
- **Locale:** production scheme kept — default `en` unprefixed, `ru`/`es` as a
  path prefix (`/ru/...`, `/es/...`); root auto-detects (session →
  Accept-Language → country-by-IP). Router optional param + resolve on root.
- **State:** zustand stays (lesson-player store); TanStack Store not adopted.

## Consequences

- Production runs **three processes**: Go API, Go async worker, and Node SSR
  (thin BFF). The Node process is the cost of runtime SSR; the Go split keeps
  HTTP availability independent from background consumption.
- Auth must be cookie-based end-to-end; the API sets/reads the JWT cookie.
- Admin CRUD depends on **regular** TypeSpec resources (`list/get/create/
  update/delete`, consistent models) — that contract shape is the first work.
- Adopted incrementally: **TanStack Table** (admin), **TanStack Form** (forms),
  hey-api **zod plugin**. Deferred until needed: TanStack Virtual, TanStack
  Pacer. Not adopted: TanStack Store, TanStack DB.
- Build/rollout order is tracked in [`docs/FRONTEND_PLAN.md`](../FRONTEND_PLAN.md).
