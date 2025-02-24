import type { PageProps } from "@inertiajs/inertia";
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

export interface SharedProps extends PageProps {
  auth: {
    user: User;
  };
  courseCategories: LanguageCategory[];
  suffix: "ru" | null;
  locale: "ru" | "en"; // TODO: use i18next.i18n["locale"];
  events: BackendEvent[] | null;
  flash: Partial<Record<FlashKey, string | null>>;
  courses: Language[];
  mobileBrowser: boolean;
  carrotquest: {
    user_hash: string;
  };
}

export type RootProps = SetupOptions<HTMLElement, SharedProps>["props"];
