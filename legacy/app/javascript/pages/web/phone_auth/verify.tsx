import type { Errors } from "@inertiajs/core";
import { Link } from "@inertiajs/react";
import { Button, Card, Container, PinInput, Stack, Text } from "@mantine/core";
import { useTranslation } from "react-i18next";
import { useAppForm } from "@/hooks/useAppForm";
import ApplicationLayout from "@/layouts/ApplicationLayout";
import * as Routes from "@/routes.js";

type Props = {
  phone: string;
};

export default function Verify({ phone }: Props) {
  const { t } = useTranslation();

  const { onSubmit, processing, form, reset } = useAppForm(
    { phone, code: "" },
    {
      url: Routes.confirm_phone_auth_path(),
      method: "post",
      onError: (_errors: Errors) => {
        reset("code");
      },
    },
  );

  const codeProps = form.getPhoneInputProps("code");

  return (
    <ApplicationLayout center header={t(($) => $.phone_auth.verify.title)}>
      <Container>
        <Stack align="center">
          <Card
            withBorder
            p="xl"
            w={{ base: "100%", sm: "80%", md: "70%", lg: "50%" }}
          >
            <Stack component="form" onSubmit={onSubmit} align="center">
              <Text c="dimmed" size="sm" ta="center">
                {t(($) => $.phone_auth.verify.hint, { phone })}
              </Text>
              <PinInput
                length={4}
                type="number"
                oneTimeCode
                autoFocus
                value={codeProps.value}
                onChange={codeProps.onChange}
                error={Boolean(codeProps.error)}
              />
              {codeProps.error ? (
                <Text c="red" size="sm">
                  {codeProps.error}
                </Text>
              ) : null}
              <Button type="submit" fullWidth loading={processing}>
                {t(($) => $.phone_auth.verify.submit)}
              </Button>
            </Stack>
          </Card>

          <Text mt="xs">
            <Text
              component={Link}
              href={Routes.new_phone_auth_path()}
              fw="bold"
              span
            >
              {t(($) => $.phone_auth.verify.change_phone)}
            </Text>
          </Text>
        </Stack>
      </Container>
    </ApplicationLayout>
  );
}
