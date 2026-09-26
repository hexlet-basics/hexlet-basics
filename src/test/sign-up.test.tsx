import { http, HttpResponse } from "msw";
import { expect, test } from "vitest";
import { page } from "vitest/browser";
import { getCourseLessonQueryKey } from "@/client/@tanstack/react-query.gen";
import { Route as signUpRoute } from "@/routes/{-$locale}/users/new";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// The sign-up page, driven through its real route with the API faked at the
// HTTP boundary: where a new account lands, and what it is shown there.

async function renderSignUp(redirect: string) {
  worker.use(http.post("*/users", () => HttpResponse.json({}, { status: 201 })));
  const rendered = await renderRoute(signUpRoute, {
    path: "/{-$locale}/users/new",
    initialPath: `/users/new?redirect=${encodeURIComponent(redirect)}`,
  });
  return rendered;
}

async function signUp() {
  await page.getByLabelText("Email", { exact: false }).fill("ada@example.com");
  await page.getByLabelText("Password", { exact: false }).fill("correct-horse-battery");
  await page.getByRole("button", { name: "Create an Account" }).click();
}

test("returns a new account to the lesson it signed up from, read afresh", async () => {
  const { router, queryClient } = await renderSignUp("/languages/javascript/lessons/variables");
  // What the guest's visit left in the cache: the lesson as a guest saw it.
  const lessonKey = getCourseLessonQueryKey({
    path: { courseSlug: "javascript", slug: "variables" },
  });
  queryClient.setQueryData(lessonKey, { guest: true });

  await signUp();

  await expect
    .poll(() => router.state.location.pathname)
    .toBe("/languages/javascript/lessons/variables");
  // The guest's copy is gone, so the player reads the progress the sign-up
  // merged into the account instead of the guest's locks.
  expect(queryClient.getQueryData(lessonKey)).toBeUndefined();
});

test("drops a redirect that would leave the site", async () => {
  for (const hostile of ["/\\evil.com", "//evil.com", "https://evil.com/languages"]) {
    const { router, screen } = await renderSignUp(hostile);
    await signUp();
    await expect.poll(() => router.state.location.href).toBe("/");
    await screen.unmount();
  }
});
