import { Button, Card, Group, Stack } from "@mantine/core";
import { useTranslation } from "react-i18next";
import type { ReviewInput } from "@/client/types.gen";
import { zReviewInput } from "@/client/zod.gen";
import { submitForm, useAppForm } from "@/lib/form";
import { submitWithServerErrors } from "@/lib/server-validation";
import type { ReactElement } from "react";

interface ReviewFormProps {
  defaultValues: ReviewInput;
  onSubmit: (values: ReviewInput) => Promise<unknown>;
  submitLabel: string;
  onCancel: () => void;
}

// The form for both create and edit. The values are the contract's request body
// (ReviewInput) and the zod schema is the one generated from it.
export function ReviewForm({
  defaultValues,
  onSubmit,
  submitLabel,
  onCancel,
}: ReviewFormProps): ReactElement {
  const { t } = useTranslation();
  const form = useAppForm({
    defaultValues,
    validators: { onSubmit: zReviewInput, onSubmitAsync: submitWithServerErrors(onSubmit) },
  });

  return (
    <Card withBorder p="xl" maw={720}>
      <Stack component="form" onSubmit={submitForm(form)}>
        <form.AppField name="pinned">
          {(field) => <field.CheckboxField label={t(($) => $.models.attributes.review.pinned)} />}
        </form.AppField>
        <form.AppField name="courseId">
          {(field) => <field.NumberField label={t(($) => $.models.attributes.review.course_id)} />}
        </form.AppField>
        <form.AppField name="userId">
          {(field) => <field.NumberField label={t(($) => $.models.attributes.review.user_id)} />}
        </form.AppField>
        <form.AppField name="firstName">
          {(field) => <field.TextField label={t(($) => $.models.attributes.review.first_name)} />}
        </form.AppField>
        <form.AppField name="lastName">
          {(field) => <field.TextField label={t(($) => $.models.attributes.review.last_name)} />}
        </form.AppField>
        <form.AppField name="body">
          {(field) => (
            <field.TextareaField label={t(($) => $.models.attributes.review.body)} required />
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
