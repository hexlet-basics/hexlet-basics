import { HttpResponse, http } from "msw";
import { describe, expect, it } from "vitest";
import { page } from "vitest/browser";
import type { BookView, User } from "@/client/types.gen";
import { Route as bookRoute } from "@/routes/{-$locale}/book";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

const DESCRIPTION_LENGTH = 160;

// The free book page, driven through its real route with the API faked at the
// HTTP boundary: which button a visitor gets, what a request does, and that
// the page exists for ru only.

const user: User = {
  admin: false,
  assistantMessagesCount: 0,
  canAccessAdmin: false,
  createdAt: "2026-09-01T00:00:00Z",
  createdAtAsTimestamp: null,
  email: "dora@example.com",
  firstName: "Dora",
  id: 1,
  lastName: null,
  name: "Dora",
  type: "user",
};

function serveBook(view: () => BookView) {
  worker.use(http.get("*/api/book", () => HttpResponse.json(view())));
}

function openBook(visitor: User | null, initialPath = "/ru/book") {
  return renderRoute(bookRoute, { initialPath, path: "/{-$locale}/book", user: visitor });
}

describe("book", () => {
  it("a visitor is sent to sign in to request the book, with no lead form", async () => {
    expect.hasAssertions();
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
    expect(match?.meta).toContainEqual({ content: "@hexlethq", name: "twitter:site" });
    expect(match?.links).toContainEqual({
      href: `${globalThis.location.origin}/ru/book`,
      rel: "canonical",
    });
    const description = match?.meta?.find((tag) => tag?.name === "description")?.content;
    expect(description).toHaveLength(DESCRIPTION_LENGTH);
  }, 60_000);

  it("a user requests the book and then gets the download link", async () => {
    expect.hasAssertions();
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
  }, 60_000);

  it("the book is not a page outside ru", async () => {
    expect.hasAssertions();
    serveBook(() => ({ requested: false }));

    const { router } = await openBook(null, "/es/book");

    await expect
      .poll(() => router.state.matches.some((match) => match.status === "notFound"))
      .toBe(true);
    expect(page.getByRole("heading", { name: "Содержание" }).query()).toBeNull();
  }, 60_000);
});
