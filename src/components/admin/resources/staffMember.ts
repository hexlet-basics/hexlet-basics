import { useTranslation } from "react-i18next";
import type { StaffMember, StaffMemberInput } from "@/client/types.gen";
import type { CrudFieldSpec } from "@/components/admin/CrudForm";

// Staff-member form values ARE the generated request body (StaffMemberInput).
// userId/roleId stay raw ids for now (searchable pickers are a later
// refinement); allowedLocales is the site's locale set.

export const emptyStaffMember: StaffMemberInput = {
  userId: 0,
  roleId: 0,
  allowedLocales: [],
};

export function staffMemberToForm(member: StaffMember): StaffMemberInput {
  return {
    userId: member.userId,
    roleId: member.roleId,
    allowedLocales: member.allowedLocales,
  };
}

export function useStaffMemberFields(): CrudFieldSpec<StaffMemberInput>[] {
  const { t } = useTranslation();
  return [
    {
      name: "userId",
      label: t(($) => $.models.attributes.staff_member.user_id),
      type: "number",
      required: true,
    },
    {
      name: "roleId",
      label: t(($) => $.models.attributes.staff_member.role_id),
      type: "number",
      required: true,
    },
    {
      name: "allowedLocales",
      label: t(($) => $.models.attributes.staff_member.allowed_locales),
      type: "multiselect",
      options: [
        { value: "en", label: "en" },
        { value: "ru", label: "ru" },
        { value: "es", label: "es" },
      ],
    },
  ];
}
