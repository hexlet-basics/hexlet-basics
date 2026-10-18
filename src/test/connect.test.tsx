import { HttpResponse, http } from "msw";
import { describe, expect, it } from "vitest";
import { page } from "vitest/browser";
import type { OAuthConsent, User } from "@/client/types.gen";
import { Route as connectRoute } from "@/routes/{-$locale}/connect";
import { Route as signInRoute } from "@/routes/{-$locale}/session/new";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// The OAuth consent screen and the sign-in that leads back to it, driven through
// their real routes with the API faked at the HTTP boundary. The consent
// operations are contract operations; the protocol endpoints they hand the
// browser on to are not under test here.

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

const consent: OAuthConsent = {
  clientName: "Claude Code",
  redirectHost: "127.0.0.1:8976",
  scopes: ["catalog:read", "progress:read", "offline_access"],
  accountName: "Dora",
};

function openConsent(visitor: User | null = user) {
  worker.use(http.get("*/api/oauth/consents/req-1", () => HttpResponse.json(consent)));
  return renderRoute(connectRoute, {
    path: "/{-$locale}/connect",
    initialPath: "/connect?id=req-1",
    user: visitor,
  });
}

describe("connect", () => {
  it("a guest is sent to sign in and the request id travels with them", async () => {
    expect.hasAssertions();
    const { router } = await openConsent(null);

    await expect.poll(() => router.state.location.pathname).toBe("/session/new");
    expect({ ...(router.state.location.search as Record<string, unknown>) }).toStrictEqual({
      redirect: "/connect?id=req-1",
    });
  }, 60_000);

  it("the consent screen shows who asks, where it returns and what it asks for", async () => {
    expect.hasAssertions();
    await openConsent();

    await expect.element(page.getByText(/Claude Code/u)).toBeVisible();
    await expect.element(page.getByText("Read the course catalog")).toBeVisible();
    await expect.element(page.getByText("Read your progress in courses")).toBeVisible();
    await expect.element(page.getByText(/sent back to 127\.0\.0\.1:8976/u)).toBeVisible();
    await expect.element(page.getByText(/Dora/u).first()).toBeVisible();
  }, 60_000);

  it("a redirect that stays on this computer carries a visible warning", async () => {
    expect.hasAssertions();
    await openConsent();

    await expect.element(page.getByText(/this computer only/u)).toBeVisible();
  }, 60_000);

  it("a redirect to another host carries no such warning", async () => {
    expect.hasAssertions();
    worker.use(
      http.get("*/api/oauth/consents/req-1", () =>
        HttpResponse.json({ ...consent, redirectHost: "app.example.com" }),
      ),
    );
    await renderRoute(connectRoute, {
      path: "/{-$locale}/connect",
      initialPath: "/connect?id=req-1",
      user,
    });

    await expect.element(page.getByText(/sent back to app\.example\.com/u)).toBeVisible();
    await expect.element(page.getByText(/this computer only/u)).not.toBeInTheDocument();
  }, 60_000);

  it("approving sends the decision for that request", async () => {
    expect.hasAssertions();
    const approved: string[] = [];
    worker.use(
      http.post("*/api/oauth/consents/req-1/approve", () => {
        approved.push("req-1");
        // Staying on this origin keeps the test page where it is.
        return HttpResponse.json({ redirectUrl: `${globalThis.location.href}#decided` });
      }),
    );
    await openConsent();

    await page.getByRole("button", { name: "Allow" }).click();

    await expect.poll(() => approved).toEqual(["req-1"]);
  }, 60_000);

  it("denying sends the refusal for that request", async () => {
    expect.hasAssertions();
    const denied: string[] = [];
    worker.use(
      http.post("*/api/oauth/consents/req-1/deny", () => {
        denied.push("req-1");
        return HttpResponse.json({ redirectUrl: `${globalThis.location.href}#decided` });
      }),
    );
    await openConsent();

    await page.getByRole("button", { name: "Deny" }).click();

    await expect.poll(() => denied).toEqual(["req-1"]);
  }, 60_000);

  it("a request that is gone is explained instead of offered", async () => {
    expect.hasAssertions();
    worker.use(
      http.get("*/api/oauth/consents/stale", () =>
        HttpResponse.json({ message: "not found" }, { status: 404 }),
      ),
    );
    await renderRoute(connectRoute, {
      path: "/{-$locale}/connect",
      initialPath: "/connect?id=stale",
      user,
    });

    await expect.element(page.getByText(/already been handled/u)).toBeVisible();
    await expect.element(page.getByRole("button", { name: "Allow" })).not.toBeInTheDocument();
  }, 60_000);

  it("signing in returns the visitor to the page they came from", async () => {
    expect.hasAssertions();
    worker.use(http.post("*/session", () => HttpResponse.json(user)));
    const { router } = await renderRoute(signInRoute, {
      path: "/{-$locale}/session/new",
      initialPath: `/session/new?redirect=${encodeURIComponent("/connect?id=req-1")}`,
    });

    await page.getByLabelText("Email", { exact: false }).fill("dora@example.com");
    await page.getByLabelText("Password", { exact: false }).fill("correct-horse-battery");
    await page.getByRole("button", { name: "Sign In" }).click();

    await expect.poll(() => router.state.location.href).toBe("/connect?id=req-1");
  }, 60_000);

  it("signing in does not follow a redirect that leaves the site", async () => {
    expect.hasAssertions();
    worker.use(http.post("*/session", () => HttpResponse.json(user)));
    const { router } = await renderRoute(signInRoute, {
      path: "/{-$locale}/session/new",
      initialPath: `/session/new?redirect=${encodeURIComponent("https://evil.example/")}`,
    });

    await page.getByLabelText("Email", { exact: false }).fill("dora@example.com");
    await page.getByLabelText("Password", { exact: false }).fill("correct-horse-battery");
    await page.getByRole("button", { name: "Sign In" }).click();

    await expect.poll(() => router.state.location.href).toBe("/");
  }, 60_000);
});
