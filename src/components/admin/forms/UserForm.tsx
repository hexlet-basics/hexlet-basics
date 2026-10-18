import { Button, Card, Group, Stack } from "@mantine/core";
import { useTranslation } from "react-i18next";
import type { UserInput } from "@/client/types.gen";
import { zUserInput } from "@/client/zod.gen";
import { submitForm, useAppForm } from "@/lib/form";
import { submitWithServerErrors } from "@/lib/server-validation";
import type { ReactElement } from "react";

interface UserFormProps {
  defaultValues: UserInput;
  onSubmit: (values: UserInput) => Promise<unknown>;
  submitLabel: string;
  onCancel: () => void;
}

// The form for both create and edit. The values are the contract's request body
// (UserInput) and the zod schema is the one generated from it.
export function UserForm({
  defaultValues,
  onSubmit,
  submitLabel,
  onCancel,
}: UserFormProps): ReactElement {
  const { t } = useTranslation();
  const form = useAppForm({
    defaultValues,
    validators: { onSubmit: zUserInput, onSubmitAsync: submitWithServerErrors(onSubmit) },
  });

  return (
    <Card withBorder p="xl" maw={720}>
      <Stack component="form" onSubmit={submitForm(form)}>
        <form.AppField name="email">
          {(field) => <field.TextField label={t(($) => $.models.attributes.user.email)} required />}
        </form.AppField>
        <form.AppField name="firstName">
          {(field) => <field.TextField label={t(($) => $.models.attributes.user.first_name)} />}
        </form.AppField>
        <form.AppField name="lastName">
          {(field) => <field.TextField label={t(($) => $.models.attributes.user.last_name)} />}
        </form.AppField>
        <form.AppField name="admin">
          {(field) => <field.CheckboxField label={t(($) => $.models.attributes.user.admin)} />}
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
