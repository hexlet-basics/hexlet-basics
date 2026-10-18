import { Button, Card, Group, Stack } from "@mantine/core";
import { useTranslation } from "react-i18next";
import type { StaffMemberInput } from "@/client/types.gen";
import { zStaffMemberInput } from "@/client/zod.gen";
import { submitForm, useAppForm } from "@/lib/form";
import { submitWithServerErrors } from "@/lib/server-validation";
import type { ReactElement } from "react";

// The locales a staff member may be limited to.
const LOCALES = ["en", "es", "ru"];

interface StaffMemberFormProps {
  defaultValues: StaffMemberInput;
  onSubmit: (values: StaffMemberInput) => Promise<unknown>;
  submitLabel: string;
  onCancel: () => void;
}

// The form for both create and edit. The values are the contract's request body
// (StaffMemberInput) and the zod schema is the one generated from it.
export function StaffMemberForm({
  defaultValues,
  onSubmit,
  submitLabel,
  onCancel,
}: StaffMemberFormProps): ReactElement {
  const { t } = useTranslation();
  const form = useAppForm({
    defaultValues,
    validators: { onSubmit: zStaffMemberInput, onSubmitAsync: submitWithServerErrors(onSubmit) },
  });

  return (
    <Card withBorder p="xl" maw={720}>
      <Stack component="form" onSubmit={submitForm(form)}>
        <form.AppField name="userId">
          {(field) => (
            <field.NumberField
              label={t(($) => $.models.attributes.staff_member.user_id)}
              required
            />
          )}
        </form.AppField>
        <form.AppField name="roleId">
          {(field) => (
            <field.NumberField
              label={t(($) => $.models.attributes.staff_member.role_id)}
              required
            />
          )}
        </form.AppField>
        <form.AppField name="allowedLocales">
          {(field) => (
            <field.MultiSelectField
              label={t(($) => $.models.attributes.staff_member.allowed_locales)}
              data={LOCALES}
            />
          )}
        </form.AppField>

        <Group justify="flex-end" mt="md">
          <Button variant="default" onClick={onCancel}>
            {t(($) => $.admin.crud.cancel)}
          </Button>
          <form.AppForm>
            <form.SubmitButton>{submitLabel}</form.SubmitButton>
          </form.AppForm>
        </Group>
      </Stack>
    </Card>
  );
}
