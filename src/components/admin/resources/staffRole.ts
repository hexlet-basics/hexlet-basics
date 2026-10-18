import type { RoleInput, StaffRoleDetail } from "@/client/types.gen";
import type { CrudFieldSpec } from "@/components/admin/CrudForm";
import { staffRoleFields } from "@/components/admin/generated/staffRole";
import { useGeneratedFields } from "@/components/admin/useGeneratedFields";

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

// Fields come from the contract (generated/staffRole.ts); what is listed here is
// only what the schema cannot say.
export function useStaffRoleFields(): CrudFieldSpec<RoleInput>[] {
  return useGeneratedFields<RoleInput>(staffRoleFields, {
    set: { description: { type: "textarea" } },
  });
}
