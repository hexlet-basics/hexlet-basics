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
import { TextLink } from "@/components/RouterLink";
import { useAppForm } from "@/lib/form";
import { safeRedirectPath } from "@/lib/safe-redirect";

// Registration page, ported from legacy users/new + SignUpFormBlock. Submits
// through the generated `createUser` mutation, which creates the account and
// sets the JWT cookie server-side. Validation reuses the generated
// `zSignUpInput` schema, tightening `firstName` to a plain (optional) string so
// the empty-string default validates instead of the contract's nullable form.
const signUpFormSchema = zSignUpInput.extend({ firstName: z.string() });

// Where to go once the account exists. The lesson player sends a guest here with
// the lesson they were on, so signing up does not lose their place. Only a path
// on this site is honoured — anything else is dropped rather than followed, so
// the page cannot be used to bounce a new account off to another origin.
const signUpSearchSchema = z.object({
  redirect: z.string().transform(safeRedirectPath).optional().catch(undefined),
});

export const Route = createFileRoute("/{-$locale}/users/new")({
  validateSearch: signUpSearchSchema,
  component: New,
});

function New() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { redirect } = Route.useSearch();
  const [serverError, setServerError] = useState<string | null>(null);

  const mutation = useMutation({
    ...createUserMutation(),
    onSuccess: async () => {
      // Everything read as a guest describes the guest. The course and lesson
      // reads carry progress, which the sign-up has just merged into the new
      // account — so the lesson a guest is sent back to must be read again,
      // not served from the cache with the guest's locks. They are dropped
      // rather than invalidated: an inactive query is only marked stale by an
      // invalidation, and a loader's queryClient.query serves stale data.
      queryClient.removeQueries({ queryKey: [{ _id: "getCourse" }] });
      queryClient.removeQueries({ queryKey: [{ _id: "getCourseLesson" }] });
      await queryClient.invalidateQueries({ queryKey: getCurrentUserQueryKey() });
      if (redirect) await navigate({ href: redirect });
      else await navigate({ to: "/{-$locale}" });
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
              void form.handleSubmit();
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
                components={{
                  a: (
                    <TextLink to="/{-$locale}/pages/$id" params={{ id: "tos" }} fz="inherit" span />
                  ),
                }}
              />
            </Text>
          </Stack>
        </Card>
      </Stack>
    </Container>
  );
}
