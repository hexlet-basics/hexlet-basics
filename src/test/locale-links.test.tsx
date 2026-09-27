import { expect, test } from "vitest";
import { page } from "vitest/browser";
import { Route as magicLinkRoute } from "@/routes/{-$locale}/magic_links/new";
import { Route as remindPasswordRoute } from "@/routes/{-$locale}/remind_password/new";
import { Route as signInRoute } from "@/routes/{-$locale}/session/new";
import { Route as signUpRoute } from "@/routes/{-$locale}/users/new";
import { renderRoute } from "@/test/renderRoute";

// A visitor on a prefixed page stays on its language when following a link,
// the way legacy's `scope "(:suffix)"` routes kept `/ru` in every URL they built.

test("keeps the visitor's locale prefix in the sign-in page's links", async () => {
  await renderRoute(signInRoute, {
    path: "/{-$locale}/session/new",
    initialPath: "/ru/session/new",
    locale: "ru",
    also: [
      { route: signUpRoute, path: "/{-$locale}/users/new" },
      { route: magicLinkRoute, path: "/{-$locale}/magic_links/new" },
      { route: remindPasswordRoute, path: "/{-$locale}/remind_password/new" },
    ],
  });

  // The page renders once its route resolves; read its links only after that.
  await expect.element(page.getByRole("heading", { level: 1 })).toBeVisible();
  const hrefs = page
    .getByRole("link")
    .elements()
    .map((link) => link.getAttribute("href"));

  expect(hrefs).toEqual(
    expect.arrayContaining(["/ru/users/new", "/ru/magic_links/new", "/ru/remind_password/new"]),
  );
});
