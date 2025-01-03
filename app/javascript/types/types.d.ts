import type { PageProps } from "@inertiajs/inertia";
import type i18next from "i18next";
import type { Language, User } from "./serializers";

export type BaseModel = {
  type: "user" | "review" | "language";
};

export type BreadcrumbItem = {
  name: string;
  url: string;
};

interface SharedProps extends PageProps {
  auth: {
    user: User;
  };
  suffix: "ru" | null;
  locale: "ru" | "en"; // TODO: use i18next.i18n["locale"];
  flash: Record<"alert" | "notice", string | null>;
  courses: Language[];
}
