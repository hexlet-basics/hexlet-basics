import type { PageProps } from "@inertiajs/inertia";
import type i18next from "i18next";
import { SortOrder } from "primereact/api";
import type {
  Language,
  LanguageCategory,
  User,
  LanguageVersion,
} from "./serializers";

export interface BaseModel extends object {
  type: "user" | "review" | "language";
}

export { SortOrder, type User, type LanguageVersion, type OriginalLanguage };

export type BreadcrumbItem = {
  name: string;
  url: string;
};

export interface SharedProps extends PageProps {
  auth: {
    user: User;
  };
  courseCategories: LanguageCategory[];
  suffix: "ru" | null;
  locale: "ru" | "en"; // TODO: use i18next.i18n["locale"];
  events: Record<string, unknown>[] | null;
  flash: Record<
    | "alert"
    | "notice"
    | "primary"
    | "secondary"
    | "success"
    | "danger"
    | "warning"
    | "info"
    | "light"
    | "dark"
    | string
    | null
  >;
  courses: Language[];
  language;
  mobileBrowser: boolean;
}

export type RootProps = SetupOptions<HTMLElement, SharedProps>["props"];
