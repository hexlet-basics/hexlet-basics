import type { PropsWithChildren } from "react";
import { Container } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import Application from "@/pages/layouts/Application";
import type { Language, LanguageCategory, User } from "@/types/serializers";

type Props = PropsWithChildren & {
  languageCategories: LanguageCategory[];
  languages: Language[];
  user: User;
};

export default function New({ languageCategories, languages, user }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");
  const { t: tAr } = useTranslation("activerecord");

  return (
    <Application languageCategories={languageCategories} languages={languages}>
      <Container>heu</Container>
    </Application>
  );
}
