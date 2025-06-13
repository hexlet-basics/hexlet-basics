import type { PropsWithChildren } from "react";
import { Alert, Card, Center, Container, Stack, Title } from '@mantine/core';

import { useTranslation } from "react-i18next";

import SignUpFormBlock from "@/components/SignUpFormBlock";
import XssContent from "@/components/XssContent";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type { User } from "@/types/serializers";
import { useForm } from "@inertiajs/react";

type Props = PropsWithChildren & {
  user: User;
  demo: boolean;
};

export default function New({ user, demo }: Props) {
  const { t } = useTranslation();
  const data = useForm()

  return (
    <ApplicationLayout
      header={t("users.new.sign_up")} center>
      <Container>
        {demo && (
          <Alert color="blue" mb="md">
            <XssContent>{t("users.new.demo_html")}</XssContent>
          </Alert>
        )}
        <Center>
          <Card withBorder p="xl" w={{ base: '100%', sm: '80%', md: '70%', lg: '50%' }}>
            <SignUpFormBlock autoFocus user={user} />
          </Card>
        </Center>
      </Container>
    </ApplicationLayout>
  );
}
