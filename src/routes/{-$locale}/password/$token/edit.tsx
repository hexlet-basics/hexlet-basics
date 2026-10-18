import { Card, Container, Stack, Title } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createFileRoute, redirect, useNavigate, useRouter } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { updatePasswordMutation } from "@/client/@tanstack/react-query.gen";
import { checkPasswordResetToken } from "@/client/sdk.gen";
import { zResetPasswordInput } from "@/client/zod.gen";
import { newPasswordInputProps } from "@/lib/authFieldProps";
import { useAppForm, submitForm } from "@/lib/form";
import { onIdentityChanged } from "@/lib/identity";
import { isApiError } from "@/lib/api-error";

// Where a refused link lands, from the loader and from a failed submission.
const refusedLink = (locale: string | undefined) => ({
  to: "/{-$locale}/remind_password/new" as const,
  params: { locale },
  search: { invalid: true },
});

// The emailed Password Reset link (legacy passwords#edit). The loader checks the
// link before the form is shown, so nobody types a new password into a dead
// one; a refused link goes back to the request form, which says it is no longer
// valid. The check answers 204 with no data to cache, so it calls the SDK
// directly rather than going through queryClient.query (ADR-0008). Submitting signs the visitor in and takes
// them home.
export const Route = createFileRoute("/{-$locale}/password/$token/edit")({
  loader: async ({ params }) => {
    try {
      await checkPasswordResetToken({ path: { token: params.token }, throwOnError: true });
    } catch (error) {
      if (isApiError(error, 404)) {
        throw redirect(refusedLink(params.locale));
      }
      throw error;
    }
  },
  component: Edit,
});

function Edit() {
  const { t } = useTranslation();
  const { locale, token } = Route.useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const router = useRouter();

  const mutation = useMutation({
    ...updatePasswordMutation(),
    onSuccess: async (user) => {
      await onIdentityChanged(queryClient, router, user);
      notifications.show({
        color: "green",
        message: t(($) => $.flash.passwords.update.success),
      });
      void navigate({ to: "/{-$locale}", params: { locale } });
    },
    onError: (error) => {
      // The link can die between the check and the submission: it expired,
      // or it was already used from another tab.
      if (isApiError(error, 404)) {
        void navigate(refusedLink(locale));
        return;
      }
      notifications.show({ color: "red", message: t(($) => $.common.errors.network) });
    },
  });

  const form = useAppForm({
    defaultValues: { password: "" },
    validators: { onSubmit: zResetPasswordInput },
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync({ path: { token }, body: value });
    },
  });

  return (
    <Container my="xl">
      <Stack align="center">
        <Title order={1} ta="center">
          {t(($) => $.passwords.edit.title)}
        </Title>

        <Card withBorder p="xl" w={{ base: "100%", sm: "80%", md: "70%", lg: "50%" }}>
          <Stack component="form" onSubmit={submitForm(form)}>
            <form.AppField name="password">
              {(field) => (
                <field.TextField
                  label={t(($) => $.passwords.edit.new_password)}
                  description={t(($) => $.passwords.edit.password_hint)}
                  required
                  {...newPasswordInputProps}
                />
              )}
            </form.AppField>

            <form.AppForm>
              <form.SubmitButton fullWidth>{t(($) => $.passwords.edit.submit)}</form.SubmitButton>
            </form.AppForm>
          </Stack>
        </Card>
      </Stack>
    </Container>
  );
}
