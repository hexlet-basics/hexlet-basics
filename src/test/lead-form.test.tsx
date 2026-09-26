import Cookies from "js-cookie";
import { http, HttpResponse } from "msw";
import { afterEach, expect, test } from "vitest";
import { page } from "vitest/browser";
import type { FirstVisit, Lead, LeadInput, User } from "@/client/types.gen";
import { Route as newLeadRoute } from "@/routes/{-$locale}/leads/new";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

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

const created: Lead = {
  id: 7,
  userId: 1,
  fullName: null,
  email: null,
  phone: null,
  whatsapp: null,
  telegram: "@dora",
  surveyAnswersData: "[]",
  coursesData: "[]",
  createdAt: "2026-09-26T00:00:00Z",
};

const firstVisit: FirstVisit = {
  utmSource: "vk",
  utmMedium: "cpc",
  utmCampaign: "autumn",
  utmContent: null,
  utmTerm: null,
  landingPage: "https://code-basics.com/ru?utm_source=vk",
  referrer: "https://vk.com/",
};

afterEach(() => {
  Cookies.remove("first_visit", { path: "/" });
  delete window.ymClientId;
});

function acceptLeads(received: LeadInput[]) {
  worker.use(
    http.post("*/leads", async ({ request }) => {
      received.push((await request.json()) as LeadInput);
      return HttpResponse.json(created, { status: 201 });
    }),
  );
}

function openForm(initialPath = "/ru/leads/new") {
  return renderRoute(newLeadRoute, { path: "/{-$locale}/leads/new", initialPath, user });
}

async function submitTelegram(handle: string) {
  await page.getByRole("textbox", { name: "Телефон / Имя пользователя" }).fill(handle);
  await page.getByRole("button", { name: "Отправить" }).click();
}

test("a lead carries the first visit and the Metrika client id, then returns home", async () => {
  const received: LeadInput[] = [];
  acceptLeads(received);
  Cookies.set("first_visit", JSON.stringify(firstVisit), { path: "/" });
  window.ymClientId = "1700000000123456789";

  const { router } = await openForm();
  await submitTelegram("@dora");

  await expect.poll(() => received).toHaveLength(1);
  expect(received[0]).toEqual({
    contactMethod: "telegram",
    contactValue: "@dora",
    ymClientId: "1700000000123456789",
    firstVisit,
  });
  await expect.element(page.getByText(/Заявка отправлена!/)).toBeVisible();
  await expect.poll(() => router.state.location.pathname).toBe("/ru");
});

test("a lead from a browser with no first visit still goes out, and returns to where it came from", async () => {
  const received: LeadInput[] = [];
  acceptLeads(received);

  const { router } = await openForm("/ru/leads/new?from=%2Fru%2Flanguages%2Fpython");
  await expect.element(page.getByRole("link", { name: "Вернуться" })).toBeVisible();
  await submitTelegram("@dora");

  await expect.poll(() => received).toHaveLength(1);
  expect(received[0]).toMatchObject({ ymClientId: null, firstVisit: null });
  await expect.poll(() => router.state.location.pathname).toBe("/ru/languages/python");
});

test("a refused lead says to check the form and stays on the page", async () => {
  worker.use(
    http.post("*/leads", () => HttpResponse.json({ message: "invalid" }, { status: 422 })),
  );

  const { router } = await openForm();
  await submitTelegram("@dora");

  await expect.element(page.getByText("Проверьте ошибки в форме")).toBeVisible();
  expect(router.state.location.pathname).toBe("/ru/leads/new");
});

test("a visitor is sent to sign in first", async () => {
  const { router } = await renderRoute(newLeadRoute, {
    path: "/{-$locale}/leads/new",
    initialPath: "/ru/leads/new",
  });

  await expect.poll(() => router.state.location.pathname).toBe("/session/new");
});
