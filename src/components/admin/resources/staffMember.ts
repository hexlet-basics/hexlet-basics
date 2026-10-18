import type { StaffMember, StaffMemberInput } from "@/client/types.gen";
import type { CrudFieldSpec } from "@/components/admin/CrudForm";
import { staffMemberFields } from "@/components/admin/generated/staffMember";
import { useGeneratedFields } from "@/components/admin/useGeneratedFields";

// Staff-member form values ARE the generated request body (StaffMemberInput).
// userId/roleId stay raw ids for now (searchable pickers are a later
// refinement); allowedLocales is the site's locale set.

export const emptyStaffMember: StaffMemberInput = {
  allowedLocales: [],
  roleId: 0,
  userId: 0,
};

export function staffMemberToForm(member: StaffMember): StaffMemberInput {
  return {
    allowedLocales: member.allowedLocales,
    roleId: member.roleId,
    userId: member.userId,
  };
}

// Fields come from the contract (generated/staffMember.ts); what is listed here is
// only what the schema cannot say.
export function useStaffMemberFields(): CrudFieldSpec<StaffMemberInput>[] {
  return useGeneratedFields<StaffMemberInput>(staffMemberFields, {
    // The contract only says "array of strings": the site's locale set is not in it.
    set: {
      allowedLocales: {
        options: [
          { label: "en", value: "en" },
          { label: "ru", value: "ru" },
          { label: "es", value: "es" },
        ],
        required: false,
      },
    },
  });
}
