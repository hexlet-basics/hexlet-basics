import type { PageProps } from "@inertiajs/inertia";
import type { i18n } from "i18next";
import type { SetupOptions } from "node_modules/@inertiajs/react/types/createInertiaApp";
import type { ReactNode } from "react";
import type { BackendEvent } from "./events";
import { LanguageCategory, LanguageLandingPageForLists, User } from "./serializers";
import type { SortOrder } from "primereact/api";

export * from "./serializers";
export { SortOrder }

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
  landingPagesForFooter: LanguageLandingPageForLists[];
  mobileBrowser: boolean;
  carrotQuestUserHash: string | null;
  metaTagsHTMLString: string;
}

export type RootProps = SetupOptions<HTMLElement, SharedProps>["props"];

// Temporary type definition, until @inertiajs/react provides one
export type ResolvedComponent = {
  default: ReactNode & {
    layout?: (page: ReactNode) => ReactNode;
  };
  layout?: (page: ReactNode) => ReactNode;
};
