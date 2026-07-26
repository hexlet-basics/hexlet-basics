import { useTranslation } from "react-i18next";
import AdminLayout from "@/layouts/AdminLayout";
import * as Routes from "@/routes.js";
import type { StaffRoleCrud } from "@/types";
import Form from "./shared/form";
import { Menu } from "./shared/menu";

type Props = {
  roleCrud: StaffRoleCrud;
};

export default function New({ roleCrud }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout header={t(($) => $.admin.management.roles.new.header)}>
      <Menu />
      <Form data={roleCrud} url={Routes.admin_management_roles_path()} />
    </AdminLayout>
  );
}
