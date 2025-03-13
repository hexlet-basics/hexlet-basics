import AdminLayout from "@/pages/layouts/AdminLayout";
import * as Routes from "@/routes.js";
import type { LanguageCategoryCrud } from "@/types/serializers";
import { useTranslation } from "react-i18next";

import Form from "./shared/form";
import { Menu } from "./shared/menu";

type Props = {
  landingPageDto: LanguageCategoryCrud;
  // courseVersions: LanguageVersion[];
};

export default function Edit({ landingPageDto }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout
      header={t("admin.language_categories.edit.header", {
        id: landingPageDto.language_category.name_ru,
      })}
    >
      <Menu data={landingPageDto} />
      <Form
        method="patch"
        data={landingPageDto}
        url={Routes.admin_language_category_path(
          landingPageDto.language_category.id,
        )}
      />
    </AdminLayout>
  );
}
