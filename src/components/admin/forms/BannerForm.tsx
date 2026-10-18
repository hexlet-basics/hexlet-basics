import { Button, Card, Group, Stack } from "@mantine/core";
import { useTranslation } from "react-i18next";
import type { BannerInput } from "@/client/types.gen";
import { zBannerInput } from "@/client/zod.gen";
import { submitForm, useAppForm } from "@/lib/form";
import { submitWithServerErrors } from "@/lib/server-validation";
import type { ReactElement } from "react";

interface BannerFormProps {
  defaultValues: BannerInput;
  onSubmit: (values: BannerInput) => Promise<unknown>;
  submitLabel: string;
  onCancel: () => void;
}

// The banner form for both create and edit. The values are the contract's
// request body (BannerInput) and the zod schema is the one generated from it,
// so the form validates against the contract with no schema of its own.
export function BannerForm({
  defaultValues,
  onSubmit,
  submitLabel,
  onCancel,
}: BannerFormProps): ReactElement {
  const { t } = useTranslation();
  const form = useAppForm({
    defaultValues,
    validators: { onSubmit: zBannerInput, onSubmitAsync: submitWithServerErrors(onSubmit) },
  });

  return (
    <Card withBorder p="xl" maw={720}>
      <Stack component="form" onSubmit={submitForm(form)}>
        <form.AppField name="background">
          {(field) => (
            <field.SelectField
              label={t(($) => $.models.attributes.banner.background)}
              data={zBannerInput.shape.background.options.map((value) => ({
                label: t(($) => $.models.attributes.banner["background/values"][value]),
                value,
              }))}
              required
            />
          )}
        </form.AppField>
        <form.AppField name="locale">
          {(field) => (
            <field.SelectField
              label={t(($) => $.models.attributes.banner.locale)}
              data={zBannerInput.shape.locale.options.map((value) => ({
                label: t(($) => $.models.attributes.banner["locale/values"][value]),
                value,
              }))}
              required
            />
          )}
        </form.AppField>
        <form.AppField name="body">
          {(field) => (
            <field.TextareaField label={t(($) => $.models.attributes.banner.body)} required />
          )}
        </form.AppField>
        <form.AppField name="url">
          {(field) => <field.TextField label={t(($) => $.models.attributes.banner.url)} />}
        </form.AppField>
        <form.AppField name="startsAt">
          {(field) => (
            <field.DateTimeField label={t(($) => $.models.attributes.banner.starts_at)} />
          )}
        </form.AppField>
        <form.AppField name="finishesAt">
          {(field) => (
            <field.DateTimeField label={t(($) => $.models.attributes.banner.finishes_at)} />
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
