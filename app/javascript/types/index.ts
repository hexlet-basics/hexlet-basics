import type { PageProps } from "@inertiajs/core";
import type { ResolvedComponent } from "@inertiajs/react";
import type { i18n } from "i18next";
import type { AssistantMessage } from "./assistantMessage";
import type {
  LanguageCategory,
  LanguageLandingPageForLists,
  LanguageLesson,
  LanguageLessonMember,
} from "./serializers";

export * from "./serializers";

export type BreadcrumbItem = {
  name: string;
  url: string;
};

export type FlashKey = "success" | "notice" | "error" | "warning";

export type FlashVariant =
  | "blue"
  | "cyan"
  | "dark"
  | "grape"
  | "gray"
  | "green"
  | "indigo"
  | "lime"
  | "orange"
  | "pink"
  | "red"
  | "teal"
  | "violet"
  | "yellow";

export type FlashVariants = Record<FlashKey, FlashVariant>;

export type Locale = "ru" | "en" | "es";

export type InertiaPageModule = {
  default: ResolvedComponent;
};

export type LessonSharedProps = PageProps & {
  canCreateAssistantMessage: boolean;
  course: Locale;
  courseCategory?: LanguageCategory;
  landingPage: LanguageLandingPageForLists;
  prevLesson?: LanguageLesson;
  nextLesson?: LanguageLesson;
  lessonMember?: LanguageLessonMember;
  lesson: LanguageLesson;
  lessons: LanguageLesson[];
  previousMessages: AssistantMessage[];
};
