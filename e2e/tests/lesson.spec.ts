import { expect, test } from "@playwright/test";
import { learner, signIn, visit } from "../support/app";

// The player titles the lesson with the course's landing page in the request
// locale; the fixtures give the JavaScript course no English page but "Frontend".
const lesson = "/languages/javascript/lessons/hello-world";

test("a guest reads a lesson, finds the later ones locked and is invited to sign in", async ({
  page,
}) => {
  await visit(page, lesson);

  await expect(
    page.getByRole("heading", { level: 1, name: "Frontend: Hello, World!" }),
  ).toBeVisible();
  await expect(page.getByRole("tab", { name: "Lesson", selected: true })).toBeVisible();
  await expect(page.getByRole("textbox", { name: "Code editor" })).toBeAttached();

  await page.getByRole("tab", { name: "Navigation" }).click();
  await expect(page.getByRole("link", { name: "Locked Variables" })).toBeVisible();

  await page.getByRole("tab", { name: "Assistant" }).click();
  await expect(page.getByText("please sign up or log in")).toBeVisible();
});

test("a signed-in learner opens the same lesson and can use the assistant", async ({ page }) => {
  await signIn(page, learner.email, learner.name);
  await visit(page, lesson);

  await expect(
    page.getByRole("heading", { level: 1, name: "Frontend: Hello, World!" }),
  ).toBeVisible();

  await page.getByRole("tab", { name: "Assistant" }).click();
  await expect(page.getByText("please sign up or log in")).toBeHidden();
  await expect(page.getByText("I'm here to help you with any questions")).toBeVisible();

  // The question box is live: Send stays disabled until there is something to
  // send. Nothing is sent, because no language model is configured in e2e.
  const send = page.getByRole("button", { name: "Send" });
  await expect(send).toBeDisabled();
  await page.getByRole("textbox", { name: "Your question" }).fill("How do I print a line?");
  await expect(send).toBeEnabled();
});
