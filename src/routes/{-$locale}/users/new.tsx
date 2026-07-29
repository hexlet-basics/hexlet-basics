import { Alert, Anchor, Box, Button, Card, Container, Stack, Text, Title } from "@mantine/core";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { Trans, useTranslation } from "react-i18next";
import { z } from "zod";
import { createUserMutation, getCurrentUserQueryKey } from "@/client/@tanstack/react-query.gen";
import { zSignUpInput } from "@/client/zod.gen";
import {
  firstNameInputProps,
  newPasswordInputProps,
  registrationEmailInputProps,
} from "@/lib/authFieldProps";
import { useAppForm } from "@/lib/form";

// Registration page, ported from legacy users/new + SignUpFormBlock. Submits
// through the generated `createUser` mutation, which creates the account and
// sets the JWT cookie server-side. Validation reuses the generated
// `zSignUpInput` schema, tightening `firstName` to a plain (optional) string so
// the empty-string default validates instead of the contract's nullable form.
const signUpFormSchema = zSignUpInput.extend({ firstName: z.string() });

export const Route = createFileRoute("/{-$locale}/users/new")({
  component: New,
});

function New() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [serverError, setServerError] = useState<string | null>(null);

  const mutation = useMutation({
    ...createUserMutation(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: getCurrentUserQueryKey() });
      navigate({ to: "/{-$locale}" });
    },
    onError: () => setServerError(t(($) => $.flash.users.create.error)),
  });

  const form = useAppForm({
    defaultValues: { firstName: "", email: "", password: "" },
    validators: { onSubmit: signUpFormSchema },
    onSubmit: async ({ value }) => {
      setServerError(null);
      await mutation.mutateAsync({ body: value });
    },
  });

  return (
    <Container my="xl">
      <Stack align="center">
        <Title order={1} ta="center">
          {t(($) => $.users.new.title)}
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

            <form.AppField name="firstName">
              {(field) => (
                <field.TextField
                  label={t(($) => $.models.attributes.user.first_name)}
                  {...firstNameInputProps}
                />
              )}
            </form.AppField>

            <form.AppField name="email">
              {(field) => (
                <field.TextField
                  label={t(($) => $.models.attributes.user.email)}
                  required
                  {...registrationEmailInputProps}
                />
              )}
            </form.AppField>

            <form.AppField name="password">
              {(field) => (
                <field.TextField
                  label={t(($) => $.models.attributes.user.password)}
                  required
                  {...newPasswordInputProps}
                />
              )}
            </form.AppField>

            <Box my="lg" ta="right">
              {t(($) => $.users.new.have_account)}{" "}
              <Anchor component={Link} to="/session/new" fw="bold">
                {t(($) => $.users.new.sign_in)}
              </Anchor>
            </Box>

            <Button type="submit" fullWidth loading={mutation.isPending}>
              {t(($) => $.helpers.submit.user_sign_up_form.create)}
            </Button>

            <Text fz="sm" mt="xs">
              <Trans
                t={t}
                i18nKey={($) => $.users.new.confirmation_html}
                components={{ a: <Text span fw="bold" /> }}
              />
            </Text>
          </Stack>
        </Card>
      </Stack>
    </Container>
  );
}
