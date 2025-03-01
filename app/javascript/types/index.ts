import type { PageProps } from "@inertiajs/inertia";
import type { i18n } from "i18next";
import type { SetupOptions } from "node_modules/@inertiajs/react/types/createInertiaApp";
import type { SortOrder } from "primereact/api";
import type { BackendEvent } from "./events";
import type {
  Language,
  LanguageCategory,
  LanguageVersion,
  OriginalLanguage,
  User,
} from "./serializers";

export type { SortOrder, User, LanguageVersion, OriginalLanguage };

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
  courseCategories: LanguageCategory[];
  suffix: "ru" | null;
  locale: Locale;
  events: BackendEvent[] | null;
  flash: Partial<Record<FlashKey, string | null>>;
  courses: Language[];
  mobileBrowser: boolean;
  carrotQuestUserHash: string | null;
  metaTagsHTMLString: string;
}

export type RootProps = SetupOptions<HTMLElement, SharedProps>["props"];
