import type { Errors } from "@inertiajs/core";
import { Link } from "@inertiajs/react";
import { Button, Card, Container, Stack, Text, TextInput } from "@mantine/core";
import type { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import { useAppForm } from "@/hooks/useAppForm";
import ApplicationLayout from "@/layouts/ApplicationLayout";
import * as Routes from "@/routes.js";

type Props = PropsWithChildren;

const phoneInputProps = {
  autoComplete: "tel" as const,
  enterKeyHint: "go" as const,
  inputMode: "tel" as const,
  type: "tel" as const,
};

export default function New(_props: Props) {
  const { t } = useTranslation();

  const { onSubmit, processing, form } = useAppForm(
    { phone: "" },
    {
      url: Routes.phone_auth_path(),
      method: "post",
      onError: (_errors: Errors) => {},
    },
  );

  return (
    <ApplicationLayout center header={t(($) => $.phone_auth.new.title)}>
      <Container>
        <Stack align="center">
          <Card
            withBorder
            p="xl"
            w={{ base: "100%", sm: "80%", md: "70%", lg: "50%" }}
          >
            <Stack component="form" onSubmit={onSubmit}>
              <Text c="dimmed" size="sm">
                {t(($) => $.phone_auth.new.hint)}
              </Text>
              <TextInput
                label={t(($) => $.phone_auth.new.phone)}
                placeholder="+7 999 123-45-67"
                {...form.getInputProps("phone")}
                required
                autoFocus
                {...phoneInputProps}
              />
              <Button type="submit" fullWidth loading={processing}>
                {t(($) => $.phone_auth.new.submit)}
              </Button>
            </Stack>
          </Card>

          <Text mt="xs">
            <Text
              component={Link}
              href={Routes.new_session_path()}
              fw="bold"
              span
            >
              {t(($) => $.phone_auth.new.other_methods)}
            </Text>
          </Text>
        </Stack>
      </Container>
    </ApplicationLayout>
  );
}
