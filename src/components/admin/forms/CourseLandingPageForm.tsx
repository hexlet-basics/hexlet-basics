import { Button, Card, Group, Stack } from "@mantine/core";
import { useTranslation } from "react-i18next";
import type { CourseLandingPageInput } from "@/client/types.gen";
import { zCourseLandingPageInput } from "@/client/zod.gen";
import { submitForm, useAppForm } from "@/lib/form";
import { submitWithServerErrors } from "@/lib/server-validation";
import type { ReactElement } from "react";

interface CourseLandingPageFormProps {
  defaultValues: CourseLandingPageInput;
  onSubmit: (values: CourseLandingPageInput) => Promise<unknown>;
  submitLabel: string;
  onCancel: () => void;
}

// The form for both create and edit. The values are the contract's request body
// (CourseLandingPageInput) and the zod schema is the one generated from it.
export function CourseLandingPageForm({
  defaultValues,
  onSubmit,
  submitLabel,
  onCancel,
}: CourseLandingPageFormProps): ReactElement {
  const { t } = useTranslation();
  const form = useAppForm({
    defaultValues,
    validators: {
      onSubmit: zCourseLandingPageInput,
      onSubmitAsync: submitWithServerErrors(onSubmit),
    },
  });

  return (
    <Card withBorder p="xl" maw={720}>
      <Stack component="form" onSubmit={submitForm(form)}>
        <form.AppField name="courseId">
          {(field) => (
            <field.NumberField
              label={t(($) => $.models.attributes.course_landing_page.course_id)}
              required
            />
          )}
        </form.AppField>
        <form.AppField name="name">
          {(field) => (
            <field.TextField
              label={t(($) => $.models.attributes.course_landing_page.name)}
              required
            />
          )}
        </form.AppField>
        <form.AppField name="slug">
          {(field) => (
            <field.TextField
              label={t(($) => $.models.attributes.course_landing_page.slug)}
              required
            />
          )}
        </form.AppField>
        <form.AppField name="main">
          {(field) => (
            <field.CheckboxField label={t(($) => $.models.attributes.course_landing_page.main)} />
          )}
        </form.AppField>
        <form.AppField name="listed">
          {(field) => (
            <field.CheckboxField label={t(($) => $.models.attributes.course_landing_page.listed)} />
          )}
        </form.AppField>
        <form.AppField name="footer">
          {(field) => (
            <field.CheckboxField label={t(($) => $.models.attributes.course_landing_page.footer)} />
          )}
        </form.AppField>
        <form.AppField name="footerName">
          {(field) => (
            <field.TextField
              label={t(($) => $.models.attributes.course_landing_page.footer_name)}
            />
          )}
        </form.AppField>
        <form.AppField name="order">
          {(field) => (
            <field.TextField label={t(($) => $.models.attributes.course_landing_page.order)} />
          )}
        </form.AppField>
        <form.AppField name="header">
          {(field) => (
            <field.TextField label={t(($) => $.models.attributes.course_landing_page.header)} />
          )}
        </form.AppField>
        <form.AppField name="description">
          {(field) => (
            <field.TextareaField
              label={t(($) => $.models.attributes.course_landing_page.description)}
            />
          )}
        </form.AppField>
        <form.AppField name="metaTitle">
          {(field) => (
            <field.TextField label={t(($) => $.models.attributes.course_landing_page.meta_title)} />
          )}
        </form.AppField>
        <form.AppField name="metaDescription">
          {(field) => (
            <field.TextareaField
              label={t(($) => $.models.attributes.course_landing_page.meta_description)}
            />
          )}
        </form.AppField>
        <form.AppField name="usedInHeader">
          {(field) => (
            <field.TextField
              label={t(($) => $.models.attributes.course_landing_page.used_in_header)}
            />
          )}
        </form.AppField>
        <form.AppField name="usedInDescription">
          {(field) => (
            <field.TextareaField
              label={t(($) => $.models.attributes.course_landing_page.used_in_description)}
            />
          )}
        </form.AppField>
        <form.AppField name="outcomesHeader">
          {(field) => (
            <field.TextField
              label={t(($) => $.models.attributes.course_landing_page.outcomes_header)}
            />
          )}
        </form.AppField>
        <form.AppField name="outcomesDescription">
          {(field) => (
            <field.TextareaField
              label={t(($) => $.models.attributes.course_landing_page.outcomes_description)}
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
