import type { PropsWithChildren } from "react";
import { Card, Container, Button, Anchor, Center, Box } from '@mantine/core';

import { useTranslation } from "react-i18next";

import * as Routes from "@/routes.js";

import { XForm, XInput } from "@/components/forms";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type { PasswordReminderForm } from "@/types";
import { Link } from "@inertiajs/react";
import AppAnchor from "@/components/AppAnchor";

type Props = PropsWithChildren & {
  passwordReminder: PasswordReminderForm;
};

export default function New({ passwordReminder }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");

  return (
    <ApplicationLayout center header={t("remind_passwords.new.title")}>
      <Container>
        <Center>
          <Card withBorder shadow="sm" p="xl" w={{ base: '100%', sm: '80%', md: '70%', lg: '50%' }}>
            <XForm
              model="remind_password_form"
              data={{ remind_password_form: passwordReminder }}
              to={Routes.remind_password_path()}
            >
              <XInput field="email" required autoFocus autoComplete="email" />
              <Box my="lg" ta="right">
                {t("users.new.have_account")}{" "}
                <AppAnchor
                  fw="bold"
                  href={Routes.new_session_path()}
                  style={{ textDecoration: 'none' }}
                >
                  {t("users.new.sign_in")}
                </AppAnchor>
              </Box>
              <Button fullWidth type="submit">
                {tHelpers("submit.remind_password_form.create")}
              </Button>
            </XForm>
          </Card>
        </Center>
      </Container>
    </ApplicationLayout>
  );
}
