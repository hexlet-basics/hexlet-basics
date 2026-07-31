import { useTranslation } from "react-i18next";
import type { RoleInput, StaffRoleDetail } from "@/client/types.gen";
import type { CrudFieldSpec } from "@/components/admin/CrudForm";

// Role form values ARE the generated request body (RoleInput). The permission
// matrix is a separate endpoint (rolePermissions) and lands in the specials
// wave — this form covers only the role identity.

export const emptyStaffRole: RoleInput = {
  name: "",
  description: "",
};

export function staffRoleToForm(role: StaffRoleDetail): RoleInput {
  return {
    name: role.name,
    description: role.description ?? "",
  };
}

export function useStaffRoleFields(): CrudFieldSpec<RoleInput>[] {
  const { t } = useTranslation();
  return [
    {
      name: "name",
      label: t(($) => $.models.attributes.base.name),
      required: true,
    },
    {
      name: "description",
      label: t(($) => $.models.attributes.base.description),
      type: "textarea",
    },
  ];
}
