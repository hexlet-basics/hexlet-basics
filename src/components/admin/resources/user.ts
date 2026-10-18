import type { UserCrud, UserInput } from "@/client/types.gen";
import type { CrudFieldSpec } from "@/components/admin/CrudForm";
import { userFields } from "@/components/admin/generated/user";
import { useGeneratedFields } from "@/components/admin/useGeneratedFields";

// User form values ARE the generated request body (UserInput). The same
// descriptor backs BOTH /admin/users (full CRUD) and /admin/management/users
// (edit only) — the contract reuses UserCrud/UserInput for the two surfaces.

export const emptyUser: UserInput = {
  admin: false,
  email: "",
  firstName: "",
  lastName: "",
};

export function userToForm(user: UserCrud): UserInput {
  return {
    admin: user.admin ?? false,
    email: user.email ?? "",
    firstName: user.firstName ?? "",
    lastName: user.lastName ?? "",
  };
}

// Fields come from the contract (generated/user.ts); what is listed here is
// only what the schema cannot say.
export function useUserFields(): CrudFieldSpec<UserInput>[] {
  return useGeneratedFields<UserInput>(userFields);
}
