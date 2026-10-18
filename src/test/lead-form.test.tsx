import Cookies from "js-cookie";
import { HttpResponse, http } from "msw";
import { describe, expect, it, onTestFinished } from "vitest";
import { page } from "vitest/browser";
import type { FirstVisit, Lead, LeadInput, User } from "@/client/types.gen";
import { Route as newLeadRoute } from "@/routes/{-$locale}/leads/new";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";
import type { RenderedRoute } from "@/test/renderRoute";
import { zLeadInput } from "@/client/zod.gen";

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

const created: Lead = {
  coursesData: "[]",
  createdAt: "2026-09-26T00:00:00Z",
  email: null,
  fullName: null,
  id: 7,
  phone: null,
  surveyAnswersData: "[]",
  telegram: "@dora",
  userId: 1,
  whatsapp: null,
};

const firstVisit: FirstVisit = {
  landingPage: "https://code-basics.com/ru?utm_source=vk",
  referrer: "https://vk.com/",
  utmCampaign: "autumn",
  utmContent: null,
  utmMedium: "cpc",
  utmSource: "vk",
  utmTerm: null,
};

function acceptLeads(received: LeadInput[]): void {
  worker.use(
    http.post("*/leads", async ({ request }) => {
      received.push(zLeadInput.parse(await request.json()));
      return HttpResponse.json(created, { status: 201 });
    }),
  );
}

function openForm(initialPath = "/ru/leads/new"): Promise<RenderedRoute> {
  // What a test seeds for the form, the cookie and the Metrika id, must not leak into the next one.
  onTestFinished(() => {
    Cookies.remove("first_visit", { path: "/" });
    delete globalThis.ymClientId;
  });
  return renderRoute(newLeadRoute, { initialPath, path: "/{-$locale}/leads/new", user });
}

async function submitTelegram(handle: string): Promise<void> {
  await page.getByRole("textbox", { name: "Телефон / Имя пользователя" }).fill(handle);
  await page.getByRole("button", { name: "Отправить" }).click();
}

describe("lead form", () => {
  it("a lead carries the first visit and the Metrika client id, then returns home", async () => {
    expect.hasAssertions();
    const received: LeadInput[] = [];
    acceptLeads(received);
    Cookies.set("first_visit", JSON.stringify(firstVisit), { path: "/" });
    globalThis.ymClientId = "1700000000123456789";

    const { router } = await openForm();
    await submitTelegram("@dora");

    await expect.poll(() => received).toHaveLength(1);
    expect(received[0]).toStrictEqual({
      contactMethod: "telegram",
      contactValue: "@dora",
      firstVisit,
      ymClientId: "1700000000123456789",
    });
    await expect.element(page.getByText(/Заявка отправлена!/u)).toBeVisible();
    await expect.poll(() => router.state.location.pathname).toBe("/ru");
  }, 60_000);

  it("a lead from a browser with no first visit still goes out, and returns to where it came from", async () => {
    expect.hasAssertions();
    const received: LeadInput[] = [];
    acceptLeads(received);

    const { router } = await openForm("/ru/leads/new?from=%2Fru%2Flanguages%2Fpython");
    await expect.element(page.getByRole("link", { name: "Вернуться" })).toBeVisible();
    await submitTelegram("@dora");

    await expect.poll(() => received).toHaveLength(1);
    expect(received[0]).toMatchObject({ firstVisit: null, ymClientId: null });
    await expect.poll(() => router.state.location.pathname).toBe("/ru/languages/python");
  }, 60_000);

  it("a refused lead says to check the form and stays on the page", async () => {
    expect.hasAssertions();
    worker.use(
      http.post("*/leads", () => HttpResponse.json({ message: "invalid" }, { status: 422 })),
    );

    const { router } = await openForm();
    await submitTelegram("@dora");

    await expect.element(page.getByText("Проверьте ошибки в форме")).toBeVisible();
    expect(router.state.location.pathname).toBe("/ru/leads/new");
  }, 60_000);

  it("a visitor is sent to sign in first", async () => {
    expect.hasAssertions();
    const { router } = await renderRoute(newLeadRoute, {
      initialPath: "/ru/leads/new",
      path: "/{-$locale}/leads/new",
    });

    await expect.poll(() => router.state.location.pathname).toBe("/session/new");
  }, 60_000);
});
