import { expect } from "@playwright/test";
import type { Browser } from "./browser";

/** The lesson player: the lesson, its navigation and the assistant. */
export class Lesson {
  readonly #browser: Browser;

  constructor(browser: Browser) {
    this.#browser = browser;
  }

  async open(path: string): Promise<void> {
    await this.#browser.visit(path);
  }

  async expectLessonTab(): Promise<void> {
    await expect(
      this.#browser.page.getByRole("tab", { name: "Lesson", selected: true }),
    ).toBeVisible();
  }

  async expectCodeEditor(): Promise<void> {
    await expect(this.#browser.page.getByRole("textbox", { name: "Code editor" })).toBeAttached();
  }

  async openNavigation(): Promise<void> {
    await this.#browser.page.getByRole("tab", { name: "Navigation" }).click();
  }

  async expectLinkTo(lesson: string): Promise<void> {
    await expect(this.#browser.page.getByRole("link", { name: lesson })).toBeVisible();
  }

  async openAssistant(): Promise<void> {
    await this.#browser.page.getByRole("tab", { name: "Assistant" }).click();
  }

  async expectSignInInvitation(): Promise<void> {
    await expect(this.#browser.page.getByText("please sign up or log in")).toBeVisible();
  }

  async expectNoSignInInvitation(): Promise<void> {
    await expect(this.#browser.page.getByText("please sign up or log in")).toBeHidden();
  }

  async expectAssistantGreeting(): Promise<void> {
    await expect(
      this.#browser.page.getByText("I'm here to help you with any questions"),
    ).toBeVisible();
  }

  /** The question box is live: Send stays disabled until there is something to send. */
  async expectCannotSend(): Promise<void> {
    await expect(this.#browser.page.getByRole("button", { name: "Send" })).toBeDisabled();
  }

  async expectCanSend(): Promise<void> {
    await expect(this.#browser.page.getByRole("button", { name: "Send" })).toBeEnabled();
  }

  /** Types a question without sending it: no language model is configured in e2e. */
  async typeQuestion(question: string): Promise<void> {
    await this.#browser.page.getByRole("textbox", { name: "Your question" }).fill(question);
  }
}
