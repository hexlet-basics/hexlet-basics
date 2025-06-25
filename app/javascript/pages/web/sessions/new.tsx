import type { PropsWithChildren } from "react";
import { Card, Container, Stack, Text, Anchor, Box, Button } from '@mantine/core';

import { useTranslation } from "react-i18next";

import * as Routes from "@/routes.js";

import { XForm, XInput } from "@/components/forms";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type { SignInForm } from "@/types/serializers";
import { Link } from "@inertiajs/react";
import AppAnchor from "@/components/AppAnchor";

type Props = PropsWithChildren & {
  signInForm: SignInForm;
};

export default function New({ signInForm }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");

  return (
    <ApplicationLayout
      center
      header={t("sessions.new.title")}
    >
      <Container>
        <Stack align="center">
          <Card withBorder p="xl" w={{ base: '100%', sm: '80%', md: '70%', lg: '50%' }}>
            <XForm
              to={Routes.session_path()}
              data={signInForm}
              model="user_sign_in_form"
              className="d-flex flex-column"
            >
              <XInput field="email" required autoFocus autoComplete="email" />
              <XInput
                required
                field="password"
                type="password"
                autoComplete="current-password"
              />
              <Box my="lg" ta="right">
                {t("sessions.new.forgot_password")}{" "}
                <Anchor
                  fw="bold"
                  component={Link}
                  href={Routes.new_remind_password_path()}
                >
                  {t("sessions.new.reset_password")}
                </Anchor>
              </Box>
              <Button type="submit" fullWidth>
                {tHelpers("submit.user_sign_in_form.create")}
              </Button>
            </XForm>
          </Card>

          <Text mt="xs">
            {t('sessions.new.dont_have_account')}{' '}
            <AppAnchor href={Routes.new_user_path()} fw="bold">
              {t('sessions.new.register')}
            </AppAnchor>
          </Text>

        </Stack>
      </Container>
    </ApplicationLayout>
  );
}
