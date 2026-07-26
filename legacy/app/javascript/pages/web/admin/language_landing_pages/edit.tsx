import { useTranslation } from "react-i18next";
import AdminLayout from "@/layouts/AdminLayout";
import * as Routes from "@/routes.js";
import type {
  Language,
  LanguageLandingPage,
  LanguageLandingPageQnaItem,
  LanguageLandingPageUpdate,
} from "@/types";

import Form from "./shared/form";
import { Menu } from "./shared/menu";
import QnaItemsSection from "./shared/qna_items_section";

type Props = {
  landingPageDto: LanguageLandingPageUpdate;
  landingPages: LanguageLandingPage[];
  languages: Language[];
  qnaItems: LanguageLandingPageQnaItem[];
  // courseVersions: LanguageVersion[];
};

export default function Edit({
  landingPageDto,
  landingPages,
  languages,
  qnaItems,
}: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout
      header={t(($) => $.admin.language_landing_pages.edit.header, {
        id: landingPageDto.header,
      })}
    >
      <Menu data={landingPageDto} />
      <Form
        languages={languages}
        landingPages={landingPages}
        method="patch"
        data={landingPageDto}
        url={Routes.admin_language_landing_page_path(landingPageDto.id)}
      />
      <QnaItemsSection landingPageId={landingPageDto.id} items={qnaItems} />
    </AdminLayout>
  );
}
