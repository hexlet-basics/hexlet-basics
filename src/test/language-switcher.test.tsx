import { http, HttpResponse } from "msw";
import { afterAll, beforeEach, expect, test } from "vitest";
import { page } from "vitest/browser";
import { Route as pageRoute } from "@/routes/{-$locale}/pages/$id";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// The header's language menu: every entry is a plain link to the same page
// Under the other locale's prefix, query included. Following one asks the API
// For nothing (the SSR server remembers the locale of the page it serves), so
// The browser can prefetch or open it without side effects.

// The header shows its menus from the desktop width up.
const startingViewport = { width: window.innerWidth, height: window.innerHeight };
beforeEach(() => page.viewport(1200, 800));

afterAll(() => page.viewport(startingViewport.width, startingViewport.height));

function openPage(initialPath: string, locale: "ru") {
  // The header's courses menu.
  worker.use(http.get("*/languages", () => HttpResponse.json([])));
  return renderRoute(pageRoute, { path: "/{-$locale}/pages/$id", initialPath, locale });
}

test("links each language to the same page and query under its prefix", async () => {
  await openPage("/ru/pages/about?from=mail", "ru");

  await page.getByRole("button", { name: "Русский" }).first().click();

  await expect
    .element(page.getByRole("menuitem", { name: "English" }))
    .toHaveAttribute("href", "/pages/about?from=mail");
  await expect
    .element(page.getByRole("menuitem", { name: "Español" }))
    .toHaveAttribute("href", "/es/pages/about?from=mail");
});

// The test page must stay: the link would otherwise navigate it away.
const stay = (event: Event) => event.preventDefault();

test("sends no request to the API when a language is followed", async () => {
  const requested: string[] = [];
  worker.events.on("request:start", ({ request }) => requested.push(new URL(request.url).pathname));
  document.addEventListener("click", stay, true);

  try {
    await openPage("/ru/pages/about", "ru");
    await page.getByRole("button", { name: "Русский" }).first().click();
    await page.getByRole("menuitem", { name: "English" }).click();
    await new Promise((resolve) => {
      setTimeout(resolve, 200);
    });
  } finally {
    document.removeEventListener("click", stay, true);
  }

  expect(requested.filter((path) => path.includes("locale"))).toStrictEqual([]);
});
