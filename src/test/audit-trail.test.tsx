import { http, HttpResponse } from "msw";
import { expect, test } from "vitest";
import { page } from "vitest/browser";
import type { AuditRecord, AuditRecordPage } from "@/client/types.gen";
import { Route } from "@/routes/{-$locale}/admin/audit/index";
import { admin } from "@/test/fixtures";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// The audit trail page through its real route and the admin layout, with the API
// Faked at the HTTP boundary. Asserted is what a staff member sees: the records
// With channel and outcome, the details of a change and of a read, the filters
// The page sends, the menu entry that only the holder of audit.read gets, and
// The refusals.

function record(overrides: Partial<AuditRecord>): AuditRecord {
  return {
    id: 1,
    occurredAt: "2026-10-01T10:00:00Z",
    actorUserId: 7,
    channel: "web",
    connectedAppId: null,
    tokenId: null,
    requestId: "req-1",
    address: "203.0.113.5",
    userAgent: "Browser/1.0",
    operation: "AdminUpdateBanner",
    entityType: "Banner",
    entityId: "3",
    locale: "en",
    outcome: "succeeded",
    changes: [],
    reads: null,
    ...overrides,
  };
}

const changed = record({
  changes: [
    {
      entity: "Banner",
      id: "3",
      op: "update",
      fields: [
        { name: "body", old: "Old text", new: "New text" },
        { name: "token", secret: true },
      ],
    },
  ],
});

const read = record({
  id: 2,
  channel: "mcp",
  connectedAppId: 11,
  tokenId: "tok-12",
  operation: "AdminListUsers",
  entityType: "User",
  entityId: null,
  reads: { count: 2, ids: ["41", "42"] },
});

const refused = record({ id: 3, operation: "AdminDeleteUser", outcome: "forbidden" });

const trail: AuditRecordPage = { items: [changed, read, refused], total: 3, page: 1, perPage: 25 };

function open() {
  return renderRoute(Route, {
    path: "/{-$locale}/admin/audit/",
    initialPath: "/admin/audit",
    user: admin,
  });
}

test("lists the records with their channel and outcome", async () => {
  worker.use(http.get("*/admin/audit_records", () => HttpResponse.json(trail)));

  await open();

  await expect.element(page.getByRole("heading", { name: "Audit trail", level: 2 })).toBeVisible();
  await expect.element(page.getByRole("cell", { name: "AdminUpdateBanner" })).toBeVisible();
  await expect.element(page.getByRole("cell", { name: "AdminListUsers" })).toBeVisible();
  await expect.element(page.getByRole("cell", { name: "MCP", exact: true })).toBeVisible();
  await expect
    .element(page.getByRole("cell", { name: "Website", exact: true }).first())
    .toBeVisible();
  await expect.element(page.getByRole("cell", { name: "Refused", exact: true })).toBeVisible();
  await expect.element(page.getByRole("cell", { name: "Done", exact: true }).first()).toBeVisible();
});

test("a change shows the old and new value of each field and hides a secret", async () => {
  worker.use(http.get("*/admin/audit_records", () => HttpResponse.json(trail)));
  await open();

  await page.getByRole("button", { name: "Details" }).first().click();

  await expect.element(page.getByText('"Old text"')).toBeVisible();
  await expect.element(page.getByText('"New text"')).toBeVisible();
  await expect.element(page.getByText("Secret value, not recorded")).toBeVisible();
  await expect.element(page.getByText("req-1")).toBeVisible();
});

test("a read shows how many rows it touched, their ids and the Connected App", async () => {
  worker.use(http.get("*/admin/audit_records", () => HttpResponse.json(trail)));
  await open();

  await page.getByRole("button", { name: "Details" }).nth(1).click();

  await expect.element(page.getByText("Rows read: 2")).toBeVisible();
  await expect.element(page.getByText("41, 42")).toBeVisible();
  await expect.element(page.getByText("Connected App")).toBeVisible();
});

test("applying the filters sends them with the window to the API", async () => {
  const queries: URLSearchParams[] = [];
  worker.use(
    http.get("*/admin/audit_records", ({ request }) => {
      queries.push(new URL(request.url).searchParams);
      return HttpResponse.json(trail);
    }),
  );
  await open();

  await page.getByLabelText("User id").fill("7");
  await page.getByLabelText("Channel").selectOptions("mcp");
  await page.getByLabelText("Operation").fill("AdminUpdateBanner");
  await page.getByLabelText("Entity type").fill("Banner");
  await page.getByLabelText("Entity id").fill("3");
  await page.getByLabelText("From").fill("2026-09-01");
  await page.getByLabelText("To").fill("2026-09-30");
  await page.getByRole("button", { name: "Apply" }).click();

  await expect.poll(() => queries.at(-1)?.get("operation")).toBe("AdminUpdateBanner");
  const sent = queries.at(-1);
  expect(sent?.get("actorUserId")).toBe("7");
  expect(sent?.get("channel")).toBe("mcp");
  expect(sent?.get("entityType")).toBe("Banner");
  expect(sent?.get("entityId")).toBe("3");
  const from = new Date(sent?.get("from") ?? "");
  const to = new Date(sent?.get("to") ?? "");
  // The last day picked is included: the window ends at the next midnight.
  expect(to.getTime() - from.getTime()).toBe(30 * 24 * 3600 * 1000);
  // The first request, before any filter, carried none.
  expect(queries[0]?.has("operation")).toBe(false);
});

test("a window the API refuses is explained", async () => {
  worker.use(
    http.get("*/admin/audit_records", () =>
      HttpResponse.json(
        { errors: { from: ["the period may not be longer than 92 days"] } },
        { status: 422 },
      ),
    ),
  );

  await open();

  await expect.element(page.getByText("the period may not be longer than 92 days")).toBeVisible();
});

test("a caller without audit.read is told so", async () => {
  worker.use(
    http.get("*/admin/audit_records", () =>
      HttpResponse.json({ title: "Forbidden", status: 403 }, { status: 403 }),
    ),
  );

  await open();

  await expect.element(page.getByText("No access to the audit trail")).toBeVisible();
});

test("the menu offers the audit trail only to a holder of audit.read", async () => {
  worker.use(
    http.get("*/admin/audit_records", () => HttpResponse.json(trail)),
    http.get("*/admin/capabilities", () => HttpResponse.json({ held: ["audit.read"] })),
  );

  await open();

  await expect.element(page.getByRole("link", { name: "Audit trail" })).toBeVisible();
});

test("the menu does not offer the audit trail to an administrator without the grant", async () => {
  worker.use(
    http.get("*/admin/audit_records", () => HttpResponse.json(trail)),
    http.get("*/admin/capabilities", () => HttpResponse.json({ held: [] })),
  );

  await open();

  await expect.element(page.getByRole("link", { name: "Users", exact: true })).toBeVisible();
  await expect.element(page.getByRole("link", { name: "Audit trail" })).not.toBeInTheDocument();
});
