import { http, HttpResponse } from "msw";
import { expect, test } from "vitest";
import { page } from "vitest/browser";
import type { BookView, User } from "@/client/types.gen";
import { Route as bookRoute } from "@/routes/{-$locale}/book";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// The free book page, driven through its real route with the API faked at the
// HTTP boundary: which button a visitor gets, what a request does, and that
// the page exists for ru only.

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

function serveBook(view: () => BookView) {
  worker.use(http.get("*/api/book", () => HttpResponse.json(view())));
}

function openBook(visitor: User | null, initialPath = "/ru/book") {
  return renderRoute(bookRoute, { path: "/{-$locale}/book", initialPath, user: visitor });
}

test("a visitor is sent to sign in to request the book, with no lead form", async () => {
  serveBook(() => ({ requested: false }));

  const { router } = await openBook(null);

  await expect
    .element(page.getByRole("link", { name: "Запросить книгу" }))
    .toHaveAttribute("href", "/ru/session/new");
  await expect.element(page.getByRole("heading", { name: "Содержание" })).toBeVisible();
  await expect.element(page.getByText("Глава 1", { exact: true })).toBeVisible();
  expect(page.getByRole("button", { name: "Отправить" }).query()).toBeNull();

  const match = router.state.matches.at(-1);
  expect(match?.meta).toContainEqual({
    title: "CodeBasics | Книга: Профессия программист. С нуля до трудоустройства",
  });
  expect(match?.meta).toContainEqual({ name: "twitter:site", content: "@hexlethq" });
  expect(match?.links).toContainEqual({
    rel: "canonical",
    href: `${window.location.origin}/ru/book`,
  });
  const description = match?.meta?.find((tag) => tag?.name === "description")?.content;
  expect(description).toHaveLength(160);
});

test("a user requests the book and then gets the download link", async () => {
  let requested = false;
  serveBook(() => ({ requested }));
  worker.use(
    http.post("*/api/book/create_request", () => {
      requested = true;
      return new HttpResponse(null, { status: 204 });
    }),
  );

  await openBook(user);
  await page.getByRole("button", { name: "Запросить книгу" }).click();

  await expect
    .element(page.getByText("Ура, теперь книга доступна для скачивания!", { exact: false }))
    .toBeVisible();
  const download = page.getByRole("link", { name: "Скачать книгу" });
  await expect.element(download).toHaveAttribute("href", "/api/book/download");
  await expect.element(download).toHaveAttribute("target", "_blank");
  // The consultation form is offered to a signed-in user.
  await expect.element(page.getByRole("button", { name: "Отправить" })).toBeVisible();
});

test("the book is not a page outside ru", async () => {
  serveBook(() => ({ requested: false }));

  const { router } = await openBook(null, "/es/book");

  await expect.poll(() => router.state.matches.some((m) => m.status === "notFound")).toBe(true);
  expect(page.getByRole("heading", { name: "Содержание" }).query()).toBeNull();
});
