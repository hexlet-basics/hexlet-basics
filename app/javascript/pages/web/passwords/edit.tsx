import {
  Box,
  Button,
  Card,
  Container,
  Stack,
  TextInput,
  Title,
} from "@mantine/core";
import type { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import { useAppForm } from "@/hooks/useAppForm";
import ApplicationLayout from "@/layouts/ApplicationLayout";
import * as Routes from "@/routes.js";
import type { UserPassword } from "@/types/serializers";

type Props = PropsWithChildren & {
  userPassword: UserPassword;
  resetPasswordToken: string;
};

export default function New({ userPassword, resetPasswordToken }: Props) {
  const { t } = useTranslation();

  const payload = userPassword;

  const { onSubmit, processing, form } = useAppForm(payload, {
    url: Routes.password_path({ reset_password_token: resetPasswordToken }),
    method: "patch",
  });

  return (
    <ApplicationLayout>
      <Container>
        <Stack align="center" gap="md">
          <Title order={1} ta="center" mb="md">
            {t(($) => $.passwords.edit.title)}
          </Title>
          <Card
            withBorder
            p="xl"
            w={{ base: "100%", sm: "80%", md: "70%", lg: "50%" }}
          >
            <form onSubmit={onSubmit}>
              <TextInput
                {...form.getInputProps("password")}
                type="password"
                autoComplete="new-password"
                required
              />
              <Box mt="lg" ta="right">
                <Button type="submit" loading={processing}>
                  {t(($) => $.helpers.submit.replace)}
                </Button>
              </Box>
            </form>
          </Card>
        </Stack>
      </Container>
    </ApplicationLayout>
  );
}
