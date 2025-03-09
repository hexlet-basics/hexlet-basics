import type { PageProps } from "@inertiajs/inertia";
import type { i18n } from "i18next";
import type { SetupOptions } from "node_modules/@inertiajs/react/types/createInertiaApp";
import type { SortOrder } from "primereact/api";
import type { BackendEvent } from "./events";
import type {
  Language,
  LanguageCategory,
  LanguageLandingPageForLists,
  LanguageVersion,
  User,
} from "./serializers";

export type { SortOrder, User, LanguageVersion, Language };

export type BreadcrumbItem = {
  name: string;
  url: string;
};

export type FlashKey = "success" | "notice" | "error";

export type FlashVariant =
  | "alert"
  | "notice"
  | "primary"
  | "secondary"
  | "success"
  | "danger"
  | "warning"
  | "info"
  | "light"
  | "dark";

export type FlashVariants = Record<FlashKey, FlashVariant>;

export type Locale = i18n["language"];

export interface SharedProps extends PageProps {
  auth: {
    user: User;
  };
  railsDirectUploadsUrl: string;
  courseCategories: LanguageCategory[];
  suffix: "ru" | null;
  locale: Locale;
  events: BackendEvent[] | null;
  flash: Partial<Record<FlashKey, string | null>>;
  landingPagesForLists: LanguageLandingPageForLists[];
  mobileBrowser: boolean;
  carrotQuestUserHash: string | null;
  metaTagsHTMLString: string;
}

export type RootProps = SetupOptions<HTMLElement, SharedProps>["props"];
