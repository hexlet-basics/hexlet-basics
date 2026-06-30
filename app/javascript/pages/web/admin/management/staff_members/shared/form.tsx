import type { Method } from "@inertiajs/core";
import { Button, MultiSelect, Select, TextInput } from "@mantine/core";
import { useTranslation } from "react-i18next";
import { useAjaxSelectProps, useAppForm } from "@/hooks/useAppForm";
import * as Routes from "@/routes.js";
import type { StaffMemberCrud, StaffRole } from "@/types";

type Props = {
  data: StaffMemberCrud;
  roles: StaffRole[];
  locales: string[];
  url: string;
  method?: Method;
  isNew?: boolean;
};

export default function Form({
  data,
  roles,
  locales,
  url,
  method,
  isNew,
}: Props) {
  const { t } = useTranslation();

  const { onSubmit, processing, form } = useAppForm(data, {
    url,
    method: method ?? "post",
  });

  const rolesSelectData = roles.map((role) => ({
    value: role.id.toString(),
    label: role.name,
  }));
  const localesSelectData = locales.map((locale) => ({
    value: locale,
    label: locale,
  }));

  const userSelectProps = useAjaxSelectProps(
    form,
    "user_id",
    Routes.search_admin_api_users_path(),
  );

  return (
    <form onSubmit={onSubmit}>
      {isNew ? (
        <Select
          label={t(($) => $.admin.management.staff_members.form.user)}
          {...userSelectProps}
        />
      ) : (
        <TextInput
          label={t(($) => $.admin.management.staff_members.form.user)}
          value={[data.user.name, data.user.email].filter(Boolean).join(" — ")}
          disabled
          mb="sm"
        />
      )}
      <Select
        label={t(($) => $.admin.management.staff_members.form.role)}
        {...form.getSelectProps("role_id", rolesSelectData)}
      />
      <MultiSelect
        label={t(($) => $.admin.management.staff_members.form.allowed_locales)}
        {...form.getMultiSelectProps("allowed_locales", localesSelectData)}
      />
      <Button type="submit" loading={processing} mt="md">
        {t(($) => $.helpers.submit.save)}
      </Button>
    </form>
  );
}
