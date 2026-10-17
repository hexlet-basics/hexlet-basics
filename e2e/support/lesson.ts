import { expect } from "@playwright/test";
import type { Site } from "./site";

/** The lesson player: the lesson, its navigation and the assistant. */
export class Lesson {
  readonly #site: Site;

  constructor(site: Site) {
    this.#site = site;
  }

  async open(path: string): Promise<void> {
    await this.#site.visit(path);
  }

  async expectLessonTab(): Promise<void> {
    await expect(
      this.#site.page.getByRole("tab", { name: "Lesson", selected: true }),
    ).toBeVisible();
  }

  async expectCodeEditor(): Promise<void> {
    await expect(this.#site.page.getByRole("textbox", { name: "Code editor" })).toBeAttached();
  }

  async openNavigation(): Promise<void> {
    await this.#site.page.getByRole("tab", { name: "Navigation" }).click();
  }

  async expectLinkTo(lesson: string): Promise<void> {
    await expect(this.#site.page.getByRole("link", { name: lesson })).toBeVisible();
  }

  async openAssistant(): Promise<void> {
    await this.#site.page.getByRole("tab", { name: "Assistant" }).click();
  }

  async expectSignInInvitation(): Promise<void> {
    await expect(this.#site.page.getByText("please sign up or log in")).toBeVisible();
  }

  async expectNoSignInInvitation(): Promise<void> {
    await expect(this.#site.page.getByText("please sign up or log in")).toBeHidden();
  }

  async expectAssistantGreeting(): Promise<void> {
    await expect(
      this.#site.page.getByText("I'm here to help you with any questions"),
    ).toBeVisible();
  }

  /** The question box is live: Send stays disabled until there is something to send. */
  async expectCannotSend(): Promise<void> {
    await expect(this.#site.page.getByRole("button", { name: "Send" })).toBeDisabled();
  }

  async expectCanSend(): Promise<void> {
    await expect(this.#site.page.getByRole("button", { name: "Send" })).toBeEnabled();
  }

  /** Types a question without sending it: no language model is configured in e2e. */
  async typeQuestion(question: string): Promise<void> {
    await this.#site.page.getByRole("textbox", { name: "Your question" }).fill(question);
  }
}
