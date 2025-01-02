import type { PropsWithChildren } from "react";
import { Container } from "react-bootstrap";
import RuAbout from "./parts/about.ru.tsx";

import { useTranslation } from "react-i18next";

import Application from "@/pages/layouts/Application";
import type { Language, LanguageCategory, User } from "@/types/serializers";
import i18next from "i18next";

type Props = PropsWithChildren & {
  languageCategories: LanguageCategory[];
  languages: Language[];
  page: "about" | "tos";
  user: User;
};

const mapping = {
  ru: {
    about: RuAbout,
    tos: RuAbout,
    privacy: RuAbout,
    cookie_policy: RuAbout,
    authors: RuAbout,
  },
  en: {
    about: RuAbout,
    tos: RuAbout,
    privacy: RuAbout,
    cookie_policy: RuAbout,
    authors: RuAbout,
  },
};

export default function New({ languageCategories, languages, page }: Props) {
  const Component = mapping[i18next.language][page];
  return (
    <Application languageCategories={languageCategories} languages={languages}>
      <Container>
        <Component />
      </Container>
    </Application>
  );
}
