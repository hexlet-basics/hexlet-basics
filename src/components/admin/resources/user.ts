import { useTranslation } from "react-i18next";
import type { UserCrud, UserInput } from "@/client/types.gen";
import type { CrudFieldSpec } from "@/components/admin/CrudForm";

// User form values ARE the generated request body (UserInput). The same
// descriptor backs BOTH /admin/users (full CRUD) and /admin/management/users
// (edit only) — the contract reuses UserCrud/UserInput for the two surfaces.

export const emptyUser: UserInput = {
  email: "",
  firstName: "",
  lastName: "",
  admin: false,
};

export function userToForm(user: UserCrud): UserInput {
  return {
    email: user.email ?? "",
    firstName: user.firstName ?? "",
    lastName: user.lastName ?? "",
    admin: user.admin ?? false,
  };
}

export function useUserFields(): CrudFieldSpec<UserInput>[] {
  const { t } = useTranslation();
  return [
    {
      name: "email",
      label: t(($) => $.models.attributes.user.email),
      required: true,
    },
    {
      name: "firstName",
      label: t(($) => $.models.attributes.user.first_name),
    },
    {
      name: "lastName",
      label: t(($) => $.models.attributes.user.last_name),
    },
    {
      name: "admin",
      label: t(($) => $.models.attributes.user.admin),
      type: "checkbox",
    },
  ];
}
