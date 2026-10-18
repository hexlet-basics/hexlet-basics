import { HttpResponse, http } from "msw";
import { describe, expect, it, onTestFinished } from "vitest";
import { page } from "vitest/browser";
import { Route as pageRoute } from "@/routes/{-$locale}/pages/$id";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";
import { DESKTOP_VIEWPORT } from "@/test/viewport";
import { sleep } from "@/test/wait";

// Long enough for a request a click would have sent to show up.
const SETTLE_MS = 200;

// The header's language menu: every entry is a plain link to the same page
// under the other locale's prefix, query included. Following one asks the API
// for nothing (the SSR server remembers the locale of the page it serves), so
// the browser can prefetch or open it without side effects.

// The header shows its menus from the desktop width up.
const startingViewport = { height: window.innerHeight, width: window.innerWidth };
function openPage(initialPath: string, locale: "ru") {
  // The header's courses menu.
  worker.use(http.get("*/languages", () => HttpResponse.json([])));
  return renderRoute(pageRoute, { initialPath, locale, path: "/{-$locale}/pages/$id" });
}

// The test page must stay: the link would otherwise navigate it away.
const stay = (event: Event) => event.preventDefault();

// The switcher is read at a desktop viewport, put back when the test is done.
async function atDesktop() {
  onTestFinished(() => page.viewport(startingViewport.width, startingViewport.height));
  await page.viewport(DESKTOP_VIEWPORT.width, DESKTOP_VIEWPORT.height);
}

describe("language switcher", () => {
  it("links each language to the same page and query under its prefix", async () => {
    expect.hasAssertions();
    await atDesktop();
    await openPage("/ru/pages/about?from=mail", "ru");

    await page.getByRole("button", { name: "Русский" }).first().click();

    await expect
      .element(page.getByRole("menuitem", { name: "English" }))
      .toHaveAttribute("href", "/pages/about?from=mail");
    await expect
      .element(page.getByRole("menuitem", { name: "Español" }))
      .toHaveAttribute("href", "/es/pages/about?from=mail");
  }, 60_000);

  it("sends no request to the API when a language is followed", async () => {
    expect.hasAssertions();
    await atDesktop();
    const requested: string[] = [];
    worker.events.on("request:start", ({ request }) =>
      requested.push(new URL(request.url).pathname),
    );
    document.addEventListener("click", stay, true);

    try {
      await openPage("/ru/pages/about", "ru");
      await page.getByRole("button", { name: "Русский" }).first().click();
      await page.getByRole("menuitem", { name: "English" }).click();
      await sleep(SETTLE_MS);
    } finally {
      document.removeEventListener("click", stay, true);
    }

    expect(requested.filter((path) => path.includes("locale"))).toStrictEqual([]);
  }, 60_000);
});
