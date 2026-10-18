import { Button, Card, Group, Stack } from "@mantine/core";
import { useTranslation } from "react-i18next";
import type { CourseCategoryInput } from "@/client/types.gen";
import { zCourseCategoryInput } from "@/client/zod.gen";
import { submitForm, useAppForm } from "@/lib/form";
import { submitWithServerErrors } from "@/lib/server-validation";
import type { ReactElement } from "react";

interface CourseCategoryFormProps {
  defaultValues: CourseCategoryInput;
  onSubmit: (values: CourseCategoryInput) => Promise<unknown>;
  submitLabel: string;
  onCancel: () => void;
}

// The form for both create and edit. The values are the contract's request body
// (CourseCategoryInput) and the zod schema is the one generated from it.
export function CourseCategoryForm({
  defaultValues,
  onSubmit,
  submitLabel,
  onCancel,
}: CourseCategoryFormProps): ReactElement {
  const { t } = useTranslation();
  const form = useAppForm({
    defaultValues,
    validators: { onSubmit: zCourseCategoryInput, onSubmitAsync: submitWithServerErrors(onSubmit) },
  });

  return (
    <Card withBorder p="xl" maw={720}>
      <Stack component="form" onSubmit={submitForm(form)}>
        <form.AppField name="name">
          {(field) => (
            <field.TextField label={t(($) => $.models.attributes.course_category.name)} required />
          )}
        </form.AppField>
        <form.AppField name="header">
          {(field) => (
            <field.TextField
              label={t(($) => $.models.attributes.course_category.header)}
              required
            />
          )}
        </form.AppField>
        <form.AppField name="slug">
          {(field) => (
            <field.TextField label={t(($) => $.models.attributes.course_category.slug)} required />
          )}
        </form.AppField>
        <form.AppField name="description">
          {(field) => (
            <field.TextareaField
              label={t(($) => $.models.attributes.course_category.description)}
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
