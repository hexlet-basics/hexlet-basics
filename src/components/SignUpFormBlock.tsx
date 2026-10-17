import { Anchor, Box, Button, Stack, Text } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Link, useNavigate, useRouter } from "@tanstack/react-router";
import { Trans, useTranslation } from "react-i18next";
import { z } from "zod";
import { createUserMutation } from "@/client/@tanstack/react-query.gen";
import { zSignUpInput } from "@/client/zod.gen";
import {
  firstNameInputProps,
  newPasswordInputProps,
  registrationEmailInputProps,
} from "@/lib/authFieldProps";
import { TextLink } from "@/components/RouterLink";
import { useAppForm } from "@/lib/form";
import { onIdentityChanged } from "@/lib/identity";
import { isValidationError, submitWithServerErrors } from "@/lib/server-validation";

// The registration form, ported from legacy SignUpFormBlock and shared by the
// registration page and the home page. Submits through the generated
// `createUser` mutation, which creates the account and starts the session
// server-side. Validation reuses the generated `zSignUpInput` schema,
// tightening `firstName` to a plain (optional) string so the empty-string
// default validates instead of the contract's nullable form.
const signUpFormSchema = zSignUpInput.extend({ firstName: z.string() });

// `redirect` is where to go once the account exists: the lesson player sends a
// guest to sign up with the lesson they were on, so signing up does not lose
// their place. The caller has already made it safe (a path on this site).
export default function SignUpFormBlock({ redirect }: { redirect?: string }) {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const router = useRouter();

  const mutation = useMutation({
    ...createUserMutation(),
    onSuccess: async (user) => {
      await onIdentityChanged(queryClient, router, user);
      notifications.show({
        color: "green",
        message: t(($) => $.flash.users.create.success),
      });
      if (redirect) await navigate({ href: redirect });
      else await navigate({ to: "/{-$locale}" });
    },
    // A 422 is shown at the fields (submitWithServerErrors), not as a toast.
    onError: (error) => {
      if (isValidationError(error)) return;
      notifications.show({
        color: "red",
        message: t(($) => $.flash.users.create.error),
      });
    },
  });

  const form = useAppForm({
    defaultValues: { firstName: "", email: "", password: "" },
    validators: {
      onSubmit: signUpFormSchema,
      onSubmitAsync: submitWithServerErrors((value: z.infer<typeof signUpFormSchema>) =>
        mutation.mutateAsync({ body: value }),
      ),
    },
  });

  return (
    <Stack
      component="form"
      onSubmit={(event) => {
        event.preventDefault();
        void form.handleSubmit();
      }}
    >
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
            a: <TextLink to="/{-$locale}/pages/$id" params={{ id: "tos" }} fz="inherit" span />,
          }}
        />
      </Text>
    </Stack>
  );
}
