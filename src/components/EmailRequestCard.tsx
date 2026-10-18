import { Alert, Card, Stack } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import type { UseMutationResult } from "@tanstack/react-query";
import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { zEmailInput } from "@/client/zod.gen";
import { emailInputProps } from "@/lib/authFieldProps";
import { submitForm, useAppForm } from "@/lib/form";

// The card of a page that asks for an email and answers the same whatever it is
// (password reminder, magic link): the form until the request goes through, then
// one confirmation. The mutation's own failure toast stays with the page.
export function EmailRequestCard({
  mutation,
  success,
  submit,
}: {
  mutation: Pick<
    UseMutationResult<unknown, unknown, { body: { email: string } }>,
    "isPending" | "isSuccess" | "mutateAsync"
  >;
  success: string;
  submit: string;
}) {
  const { t } = useTranslation();
  const form = useAppForm({
    defaultValues: { email: "" },
    validators: { onSubmit: zEmailInput },
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync({ body: value });
    },
  });

  return (
    <Card withBorder p="xl" w={{ base: "100%", sm: "80%", md: "70%", lg: "50%" }}>
      {mutation.isSuccess ? (
        <Alert color="green">{success}</Alert>
      ) : (
        <Stack component="form" onSubmit={submitForm(form)}>
          <form.AppField name="email">
            {(field) => (
              <field.TextField
                label={t(($) => $.models.attributes.user.email)}
                required
                {...emailInputProps}
              />
            )}
          </form.AppField>

          <form.AppForm>
            <form.SubmitButton fullWidth>{submit}</form.SubmitButton>
          </form.AppForm>
        </Stack>
      )}
    </Card>
  );
}

// The refused-link notice is a flash, so it is shown once on arrival like the
// other flashes; the id keeps StrictMode's double effect from stacking two.
export function useInvalidLinkNotice(invalid: boolean | undefined, message: string) {
  useEffect(() => {
    if (invalid) notifications.show({ id: "invalid-link", color: "red", message });
  }, [invalid, message]);
}
