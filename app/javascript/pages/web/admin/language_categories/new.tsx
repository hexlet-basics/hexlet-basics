import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import AdminLayout from "@/pages/layouts/AdminLayout";
import type { LanguageCategoryCrud } from "@/types/serializers";
import Form from "./shared/form";
import { Menu } from "./shared/menu";

type Props = {
  categoryDto: LanguageCategoryCrud;
};

export default function New({ categoryDto }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout header={t("admin.language_categories.new.header")}>
      <Menu />
      <Form
        data={categoryDto}
        url={Routes.admin_language_categories_path()}
      />
    </AdminLayout>
  );
}
