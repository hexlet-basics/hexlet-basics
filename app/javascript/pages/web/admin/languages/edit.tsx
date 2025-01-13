import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import AdminLayout from "@/pages/layouts/AdminLayout";
import type { OriginalLanguage } from "@/types/serializers";
import Form from "./shared/form";
import { Menu } from "./shared/menu";

type Props = {
  originalCourse: OriginalLanguage;
};

export default function Edit({ originalCourse }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout header={t("admin.languages.edit.header", { id: originalCourse.id })}>
      <Menu data={originalCourse} />
      <Form
        method="patch"
        data={originalCourse}
        url={Routes.admin_language_path(originalCourse.id)}
      />
    </AdminLayout>
  );
}
