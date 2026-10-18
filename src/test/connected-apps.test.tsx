import { HttpResponse, http } from "msw";
import { expect, test } from "vitest";
import { page } from "vitest/browser";
import type { ConnectedApp, User } from "@/client/types.gen";
import { Route as connectedAppsRoute } from "@/routes/{-$locale}/account/connected_apps";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// The Connected Apps page, driven through its real route with the API faked at
// The HTTP boundary: the list, the confirm step, the revoke and its refresh.

const user: User = {
  id: 1,
  firstName: "Dora",
  lastName: null,
  name: "Dora",
  email: "dora@example.com",
  admin: false,
  canAccessAdmin: false,
  assistantMessagesCount: 0,
  createdAt: "2026-09-01T00:00:00Z",
  createdAtAsTimestamp: null,
  type: "user",
};

const claude: ConnectedApp = {
  id: 7,
  name: "Claude Code",
  scopes: ["catalog:read", "progress:read", "offline_access"],
  connectedAt: "2026-10-01T12:00:00Z",
};

const documentApp: ConnectedApp = {
  id: 8,
  name: "Editor Helper",
  scopes: ["catalog:read", "something:new"],
  connectedAt: "2026-10-05T12:00:00Z",
  clientHost: "helper.example.com",
};

function openPage(apps: ConnectedApp[], visitor: User | null = user) {
  let current = apps;
  worker.use(http.get("*/api/oauth/connected_apps", () => HttpResponse.json({ items: current })));
  const result = renderRoute(connectedAppsRoute, {
    path: "/{-$locale}/account/connected_apps",
    initialPath: "/account/connected_apps",
    user: visitor,
  });
  return { result, setApps: (next: ConnectedApp[]) => (current = next) };
}

test("a guest is sent to sign in first", async () => {
  const { result } = openPage([], null);
  const { router } = await result;

  await expect.poll(() => router.state.location.pathname).toBe("/session/new");
});

test("each app shows its name, what it may do and, for a document client, its host", async () => {
  await openPage([claude, documentApp]).result;

  await expect.element(page.getByText("Claude Code")).toBeVisible();
  await expect.element(page.getByText("Read the course catalog").first()).toBeVisible();
  await expect.element(page.getByText("Read your progress in courses")).toBeVisible();
  await expect.element(page.getByText("Keep access until you revoke it")).toBeVisible();
  await expect.element(page.getByText("Editor Helper")).toBeVisible();
  await expect.element(page.getByText(/helper\.example\.com/)).toBeVisible();
  // A scope the page has no words for is shown as it is.
  await expect.element(page.getByText("something:new")).toBeVisible();
});

test("without apps the page says so", async () => {
  await openPage([]).result;

  await expect.element(page.getByText("You have not connected any applications.")).toBeVisible();
});

test("disconnecting asks first, and only a yes revokes the app and refreshes the list", async () => {
  const revoked: string[] = [];
  const { result, setApps } = openPage([claude, documentApp]);
  worker.use(
    http.delete("*/api/oauth/connected_apps/:id", ({ params }) => {
      revoked.push(String(params.id));
      setApps([documentApp]);
      return new HttpResponse(null, { status: 204 });
    }),
  );
  await result;

  await page.getByRole("button", { name: "Disconnect" }).first().click();
  await expect.element(page.getByText(/Disconnect “Claude Code”\?/)).toBeVisible();
  // Backing out revokes nothing.
  await page.getByRole("button", { name: "No" }).click();
  expect(revoked).toStrictEqual([]);

  await page.getByRole("button", { name: "Disconnect" }).first().click();
  await page.getByRole("dialog").getByRole("button", { name: "Disconnect" }).click();

  await expect.element(page.getByText("The application was disconnected")).toBeVisible();
  expect(revoked).toStrictEqual(["7"]);
  await expect.element(page.getByText("Editor Helper")).toBeVisible();
  await expect.poll(() => page.getByText("Claude Code").elements().length).toBe(0);
});

test("a failed revoke says so and keeps the app listed", async () => {
  const { result } = openPage([claude]);
  worker.use(
    http.delete("*/api/oauth/connected_apps/:id", () =>
      HttpResponse.json({ message: "boom" }, { status: 500 }),
    ),
  );
  await result;

  await page.getByRole("button", { name: "Disconnect" }).first().click();
  await page.getByRole("dialog").getByRole("button", { name: "Disconnect" }).click();

  await expect
    .element(page.getByText("Could not disconnect the application. Please try again."))
    .toBeVisible();
  await expect.element(page.getByText("Claude Code")).toBeVisible();
});
