import { Button, Card, Group, Stack } from "@mantine/core";
import { useTranslation } from "react-i18next";
import type { CourseInput } from "@/client/types.gen";
import { zCourseInput, zCourseLearnAs } from "@/client/zod.gen";
import { submitForm, useAppForm } from "@/lib/form";
import { submitWithServerErrors } from "@/lib/server-validation";
import type { ReactElement } from "react";

interface CourseFormProps {
  defaultValues: CourseInput;
  onSubmit: (values: CourseInput) => Promise<unknown>;
  submitLabel: string;
  onCancel: () => void;
}

// The form for both create and edit. The values are the contract's request body
// (CourseInput) and the zod schema is the one generated from it.
export function CourseForm({
  defaultValues,
  onSubmit,
  submitLabel,
  onCancel,
}: CourseFormProps): ReactElement {
  const { t } = useTranslation();
  const form = useAppForm({
    defaultValues,
    validators: { onSubmit: zCourseInput, onSubmitAsync: submitWithServerErrors(onSubmit) },
  });

  return (
    <Card withBorder p="xl" maw={720}>
      <Stack component="form" onSubmit={submitForm(form)}>
        <form.AppField name="slug">
          {(field) => (
            <field.TextField label={t(($) => $.models.attributes.course.slug)} required />
          )}
        </form.AppField>
        <form.AppField name="learnAs">
          {(field) => (
            <field.SelectField
              label={t(($) => $.models.attributes.course.learn_as)}
              data={zCourseLearnAs.options.map((value) => ({
                label: t(($) => $.models.attributes.course["learn_as/values"][value]),
                value,
              }))}
            />
          )}
        </form.AppField>
        <form.AppField name="hexletProgramLandingPage">
          {(field) => (
            <field.TextField
              label={t(($) => $.models.attributes.course.hexlet_program_landing_page)}
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
