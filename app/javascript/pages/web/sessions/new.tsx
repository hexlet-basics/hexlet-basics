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
import type { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import { useAppForm } from "@/hooks/useAppForm";
import ApplicationLayout from "@/layouts/ApplicationLayout";
import {
  currentPasswordInputProps,
  loginEmailInputProps,
} from "@/lib/authFieldProps";
import * as Routes from "@/routes.js";
import type { SignInForm } from "@/types/serializers";

type Props = PropsWithChildren & {
  signInForm: SignInForm;
};

export default function New({ signInForm }: Props) {
  const { t } = useTranslation();

  const payload = signInForm;

  const { onSubmit, processing, form } = useAppForm(payload, {
    url: Routes.session_path(),
    method: "post",
  });

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
                {...form.getInputProps("email")}
                required
                autoFocus
                {...loginEmailInputProps}
              />
              <TextInput
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
