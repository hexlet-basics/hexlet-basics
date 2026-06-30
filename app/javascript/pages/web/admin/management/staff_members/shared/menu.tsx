import { useTranslation } from "react-i18next";
import {
  CrudHorizontalMenu,
  type CrudHorizontalMenuItem,
} from "@/components/CrudHorizontalMenu";
import * as Routes from "@/routes.js";
import type { StaffMemberCrud } from "@/types";

type Props = {
  data?: StaffMemberCrud;
};

export function Menu({ data }: Props) {
  const { t } = useTranslation();

  const items: CrudHorizontalMenuItem[] = [
    {
      href: Routes.admin_management_staff_members_path(),
      label: t(($) => $.helpers.crud.list),
    },
    {
      href: Routes.new_admin_management_staff_member_path(),
      label: t(($) => $.helpers.crud.add),
    },
  ];

  if (data) {
    items.push({
      href: Routes.edit_admin_management_staff_member_path(data.id),
      label: t(($) => $.helpers.crud.editing),
    });
  }

  return <CrudHorizontalMenu items={items} />;
}
