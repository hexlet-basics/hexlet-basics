import type { PageProps } from "@inertiajs/inertia";
import type i18next from "i18next";
import { SortOrder } from "primereact/api";
import type { Language, LanguageCategory, User } from "./serializers";

export interface BaseModel extends object {
  type: "user" | "review" | "language";
}

export { SortOrder, type User };

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
}
