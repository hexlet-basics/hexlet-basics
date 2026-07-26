import { useTranslation } from "react-i18next";
import AdminLayout from "@/layouts/AdminLayout";
import * as Routes from "@/routes.js";
import type { StaffRoleCrud } from "@/types";
import Form from "./shared/form";
import { Menu } from "./shared/menu";

type Props = {
  roleCrud: StaffRoleCrud;
};

export default function Edit({ roleCrud }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout header={t(($) => $.admin.management.roles.edit.header)}>
      <Menu data={roleCrud} />
      <Form
        data={roleCrud}
        url={Routes.admin_management_role_path(roleCrud.id)}
        method="patch"
      />
    </AdminLayout>
  );
}
