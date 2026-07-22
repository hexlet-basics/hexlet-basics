import type { Errors } from "@inertiajs/core";
import { Link } from "@inertiajs/react";
import {
  Box,
  Button,
  Card,
  Container,
  Stack,
  Text,
  TextInput,
} from "@mantine/core";
import { notifications } from "@mantine/notifications";
import type { PropsWithChildren } from "react";
import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { useAppForm } from "@/hooks/useAppForm";
import ApplicationLayout from "@/layouts/ApplicationLayout";
import {
  currentPasswordInputProps,
  loginEmailInputProps,
} from "@/lib/authFieldProps";
import {
  cancelPasskeyCeremony,
  loginWithPasskey,
  passkeyAutofillSupported,
  passkeyCancelled,
  passkeySupported,
} from "@/lib/passkey";
import * as Routes from "@/routes.js";

type Props = PropsWithChildren;

export default function New(_props: Props) {
  const { t } = useTranslation();

  const handlePasskeyLogin = async () => {
    try {
      await loginWithPasskey();
    } catch (error) {
      if (passkeyCancelled(error)) {
        return;
      }

      notifications.show({
        color: "red",
        message: t(($) => $.sessions.new.passkey_error),
      });
    }
  };

  useEffect(() => {
    (async () => {
      if (!(await passkeyAutofillSupported())) {
        return;
      }

      try {
        await loginWithPasskey(true);
      } catch (error) {
        // Autofill runs unprompted in the background, so a toast here
        // would blame the user for an action they never took.
        if (!passkeyCancelled(error)) {
          console.warn(error);
        }
      }
    })();

    return cancelPasskeyCeremony;
  }, []);

  const { onSubmit, processing, form, reset } = useAppForm(
    { email: "", password: "" },
    {
      url: Routes.session_path(),
      method: "post",
      onError: (_errors: Errors) => {
        reset("password");
      },
    },
  );

  return (
    <ApplicationLayout center header={t(($) => $.sessions.new.title)}>
      <Container>
        <Stack align="center">
          <Card
            withBorder
            p="xl"
            w={{ base: "100%", sm: "80%", md: "70%", lg: "50%" }}
          >
            <Stack component="form" onSubmit={onSubmit}>
              <TextInput
                label={t(($) => $.models.attributes.user.email)}
                {...form.getInputProps("email")}
                required
                autoFocus
                {...loginEmailInputProps}
                autoComplete="username webauthn"
              />
              <TextInput
                label={t(($) => $.models.attributes.user.password)}
                {...form.getInputProps("password")}
                required
                {...currentPasswordInputProps}
              />
              <Box my="lg" ta="right">
                {t(($) => $.sessions.new.forgot_password)}{" "}
                <Text
                  component={Link}
                  fw="bold"
                  href={Routes.new_remind_password_path()}
                  span
                >
                  {t(($) => $.sessions.new.reset_password)}
                </Text>
              </Box>
              <Button type="submit" fullWidth loading={processing}>
                {t(($) => $.helpers.submit.user_sign_in_form.create)}
              </Button>
              <Button
                component={Link}
                href={Routes.new_phone_auth_path()}
                variant="default"
                fullWidth
              >
                {t(($) => $.sessions.new.sign_in_with_phone)}
              </Button>
              <Button
                component={Link}
                href={Routes.new_magic_link_path()}
                variant="default"
                fullWidth
              >
                {t(($) => $.sessions.new.sign_in_with_magic_link)}
              </Button>
              {passkeySupported() ? (
                <Button
                  variant="default"
                  fullWidth
                  onClick={handlePasskeyLogin}
                >
                  {t(($) => $.sessions.new.sign_in_with_passkey)}
                </Button>
              ) : null}
            </Stack>
          </Card>

          <Text mt="xs">
            {t(($) => $.sessions.new.dont_have_account)}{" "}
            <Text component={Link} href={Routes.new_user_path()} fw="bold" span>
              {t(($) => $.sessions.new.register)}
            </Text>
          </Text>
        </Stack>
      </Container>
    </ApplicationLayout>
  );
}
