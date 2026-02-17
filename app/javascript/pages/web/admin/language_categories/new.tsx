import { useTranslation } from "react-i18next";
import AdminLayout from "@/layouts/AdminLayout";
import * as Routes from "@/routes.js";
import type { LanguageCategoryCrud } from "@/types";
import Form from "./shared/form";
import { Menu } from "./shared/menu";

type Props = {
  categoryDto: LanguageCategoryCrud;
};

export default function New({ categoryDto }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout header={t(($) => $.admin.language_categories.new.header)}>
      <Menu />
      <Form data={categoryDto} url={Routes.admin_language_categories_path()} />
    </AdminLayout>
  );
}
