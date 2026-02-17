import { Alert, Card, Center, Container, Stack, Title } from "@mantine/core";
import type { PropsWithChildren } from "react";
import { Trans, useTranslation } from "react-i18next";
import SignUpFormBlock from "@/components/SignUpFormBlock";
import ApplicationLayout from "@/layouts/ApplicationLayout";
import type { User, UserSignUpForm } from "@/types/serializers";

type Props = PropsWithChildren & {
  user: UserSignUpForm;
  demo: boolean;
};

export default function New({ user, demo }: Props) {
  const { t } = useTranslation();

  return (
    <ApplicationLayout header={t(($) => $.pages.users.new.sign_up)} center>
      <Container>
        {demo && (
          <Alert color="blue" mb="md">
            <Trans
              t={t}
              i18nKey={($) => $.pages.users.new.demo_html}
              components={{
                b: <strong />,
              }}
            />
          </Alert>
        )}
        <Center>
          <Card
            withBorder
            p="xl"
            w={{ base: "100%", sm: "80%", md: "70%", lg: "50%" }}
          >
            <SignUpFormBlock autoFocus userDto={user} />
          </Card>
        </Center>
      </Container>
    </ApplicationLayout>
  );
}
