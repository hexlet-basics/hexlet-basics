import { useTranslation } from "react-i18next";
import AdminLayout from "@/layouts/AdminLayout";
import * as Routes from "@/routes.js";
import type { StaffMemberCrud, StaffRole } from "@/types";
import Form from "./shared/form";
import { Menu } from "./shared/menu";

type Props = {
  staffMemberCrud: StaffMemberCrud;
  roles: StaffRole[];
  locales: string[];
};

export default function New({ staffMemberCrud, roles, locales }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout header={t(($) => $.admin.management.staff_members.new.header)}>
      <Menu />
      <Form
        data={staffMemberCrud}
        roles={roles}
        locales={locales}
        url={Routes.admin_management_staff_members_path()}
        isNew
      />
    </AdminLayout>
  );
}
