import type { StaffMember, StaffMemberInput } from "@/client/types.gen";
import type { CrudFieldSpec } from "@/components/admin/CrudForm";
import { staffMemberFields } from "@/components/admin/generated/staffMember";
import { useGeneratedFields } from "@/components/admin/useGeneratedFields";

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

// Fields come from the contract (generated/staffMember.ts); what is listed here is
// only what the schema cannot say.
export function useStaffMemberFields(): CrudFieldSpec<StaffMemberInput>[] {
  return useGeneratedFields<StaffMemberInput>(staffMemberFields, {
    // The contract only says "array of strings": the site's locale set is not in it.
    set: {
      allowedLocales: {
        required: false,
        options: [
          { value: "en", label: "en" },
          { value: "ru", label: "ru" },
          { value: "es", label: "es" },
        ],
      },
    },
  });
}
