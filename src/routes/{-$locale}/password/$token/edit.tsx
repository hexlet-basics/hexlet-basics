import { Alert, Button, Card, Container, Stack, Title } from "@mantine/core";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createFileRoute, redirect, useNavigate } from "@tanstack/react-router";
import { isAxiosError } from "axios";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { getCurrentUserQueryKey, updatePasswordMutation } from "@/client/@tanstack/react-query.gen";
import { checkPasswordResetToken } from "@/client/sdk.gen";
import { zResetPasswordInput } from "@/client/zod.gen";
import { newPasswordInputProps } from "@/lib/authFieldProps";
import { useAppForm } from "@/lib/form";

// Where a refused link lands, from the loader and from a failed submission.
const refusedLink = (locale: string | undefined) => ({
  to: "/{-$locale}/remind_password/new" as const,
  params: { locale },
  search: { invalid: true },
});

// The emailed Password Reset link (legacy passwords#edit). The loader checks the
// link before the form is shown, so nobody types a new password into a dead
// one; a refused link goes back to the request form, which says it is no longer
// valid. The check has no answer to cache, so it calls the SDK directly rather
// than going through ensureQueryData. Submitting signs the visitor in and takes
// them home.
export const Route = createFileRoute("/{-$locale}/password/$token/edit")({
  loader: async ({ params }) => {
    try {
      await checkPasswordResetToken({ path: { token: params.token }, throwOnError: true });
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) {
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
  const [serverError, setServerError] = useState<string | null>(null);

  const mutation = useMutation({
    ...updatePasswordMutation(),
    onSuccess: (user) => {
      // Seeded from the response, as after a Magic Link: the root guard reads
      // the current user through ensureQueryData and would keep the guest.
      queryClient.setQueryData(getCurrentUserQueryKey(), { user });
      navigate({ to: "/{-$locale}", params: { locale } });
    },
    onError: (error) => {
      // The link can die between the check and the submission: it expired,
      // or it was already used from another tab.
      if (isAxiosError(error) && error.response?.status === 404) {
        navigate(refusedLink(locale));
        return;
      }
      setServerError(t(($) => $.common.errors.network));
    },
  });

  const form = useAppForm({
    defaultValues: { password: "" },
    validators: { onSubmit: zResetPasswordInput },
    onSubmit: async ({ value }) => {
      setServerError(null);
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
          <Stack
            component="form"
            onSubmit={(event) => {
              event.preventDefault();
              form.handleSubmit();
            }}
          >
            {serverError && <Alert color="red">{serverError}</Alert>}

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

            <Button type="submit" fullWidth loading={mutation.isPending}>
              {t(($) => $.passwords.edit.submit)}
            </Button>
          </Stack>
        </Card>
      </Stack>
    </Container>
  );
}
