import { HttpResponse, http } from "msw";
import { describe, expect, it } from "vitest";
import { page } from "vitest/browser";
import type { AuditRecord, AuditRecordPage } from "@/client/types.gen";
import { Route } from "@/routes/{-$locale}/admin/audit/index";
import { admin } from "@/test/fixtures";
import { worker } from "@/test/msw";
import { validationProblem } from "@/test/problem";
import { renderRoute } from "@/test/renderRoute";

const THIRTY_DAYS_MS = 2_592_000_000;

// The audit trail page through its real route and the admin layout, with the API
// faked at the HTTP boundary. Asserted is what a staff member sees: the records
// with channel and outcome, the details of a change and of a read, the filters
// the page sends, the menu entry that only the holder of audit.read gets, and
// the refusals.

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

describe("audit trail", () => {
  it("lists the records with their channel and outcome", async () => {
    expect.hasAssertions();
    worker.use(http.get("*/admin/audit_records", () => HttpResponse.json(trail)));

    await open();

    await expect
      .element(page.getByRole("heading", { name: "Audit trail", level: 2 }))
      .toBeVisible();
    await expect.element(page.getByRole("cell", { name: "AdminUpdateBanner" })).toBeVisible();
    await expect.element(page.getByRole("cell", { name: "AdminListUsers" })).toBeVisible();
    await expect.element(page.getByRole("cell", { name: "MCP", exact: true })).toBeVisible();
    await expect
      .element(page.getByRole("cell", { name: "Website", exact: true }).first())
      .toBeVisible();
    await expect.element(page.getByRole("cell", { name: "Refused", exact: true })).toBeVisible();
    await expect
      .element(page.getByRole("cell", { name: "Done", exact: true }).first())
      .toBeVisible();
  }, 60_000);

  it("a change shows the old and new value of each field and hides a secret", async () => {
    expect.hasAssertions();
    worker.use(http.get("*/admin/audit_records", () => HttpResponse.json(trail)));
    await open();

    await page.getByRole("button", { name: "Details" }).first().click();

    await expect.element(page.getByText('"Old text"')).toBeVisible();
    await expect.element(page.getByText('"New text"')).toBeVisible();
    await expect.element(page.getByText("Secret value, not recorded")).toBeVisible();
    await expect.element(page.getByText("req-1")).toBeVisible();
  }, 60_000);

  it("a read shows how many rows it touched, their ids and the Connected App", async () => {
    expect.hasAssertions();
    worker.use(http.get("*/admin/audit_records", () => HttpResponse.json(trail)));
    await open();

    await page.getByRole("button", { name: "Details" }).nth(1).click();

    await expect.element(page.getByText("Rows read: 2")).toBeVisible();
    await expect.element(page.getByText("41, 42")).toBeVisible();
    await expect.element(page.getByText("Connected App")).toBeVisible();
  }, 60_000);

  it("applying the filters sends them with the window to the API", async () => {
    expect.hasAssertions();
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
    const {
      from: fromParam,
      to: toParam,
      ...filters
    } = Object.fromEntries(new URLSearchParams(sent));
    expect(filters).toMatchObject({
      actorUserId: "7",
      channel: "mcp",
      entityType: "Banner",
      entityId: "3",
    });
    const from = new Date(String(fromParam));
    const to = new Date(String(toParam));
    // The last day picked is included: the window ends at the next midnight.
    expect(to.getTime() - from.getTime()).toBe(THIRTY_DAYS_MS);
    // The first request, before any filter, carried none.
    expect(queries[0]?.has("operation")).toBe(false);
  }, 60_000);

  it("a window the API refuses is explained", async () => {
    expect.hasAssertions();
    worker.use(
      http.get("*/admin/audit_records", () =>
        validationProblem({ from: ["the period may not be longer than 92 days"] }),
      ),
    );

    await open();

    await expect.element(page.getByText("the period may not be longer than 92 days")).toBeVisible();
  }, 60_000);

  it("a caller without audit.read is told so", async () => {
    expect.hasAssertions();
    worker.use(
      http.get("*/admin/audit_records", () =>
        HttpResponse.json({ title: "Forbidden", status: 403 }, { status: 403 }),
      ),
    );

    await open();

    await expect.element(page.getByText("No access to the audit trail")).toBeVisible();
  }, 60_000);

  it("the menu offers the audit trail only to a holder of audit.read", async () => {
    expect.hasAssertions();
    worker.use(
      http.get("*/admin/audit_records", () => HttpResponse.json(trail)),
      http.get("*/admin/capabilities", () => HttpResponse.json({ held: ["audit.read"] })),
    );

    await open();

    await expect.element(page.getByRole("link", { name: "Audit trail" })).toBeVisible();
  }, 60_000);

  it("the menu does not offer the audit trail to an administrator without the grant", async () => {
    expect.hasAssertions();
    worker.use(
      http.get("*/admin/audit_records", () => HttpResponse.json(trail)),
      http.get("*/admin/capabilities", () => HttpResponse.json({ held: [] })),
    );

    await open();

    await expect.element(page.getByRole("link", { name: "Users", exact: true })).toBeVisible();
    await expect.element(page.getByRole("link", { name: "Audit trail" })).not.toBeInTheDocument();
  }, 60_000);
});
