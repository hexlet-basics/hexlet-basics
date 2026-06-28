import { Link } from "@inertiajs/react";
import { Box, Button, Card, Center, Container, TextInput } from "@mantine/core";
import type { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import { useAppForm } from "@/hooks/useAppForm";
import ApplicationLayout from "@/layouts/ApplicationLayout";
import { emailInputProps } from "@/lib/authFieldProps";
import * as Routes from "@/routes.js";

type Props = PropsWithChildren;

export default function New(_props: Props) {
  const { t } = useTranslation();

  const { onSubmit, processing, form } = useAppForm(
    { email: "" },
    {
      url: Routes.remind_password_path(),
      method: "post",
    },
  );

  return (
    <ApplicationLayout center header={t(($) => $.remind_passwords.new.title)}>
      <Container>
        <Center>
          <Card
            withBorder
            shadow="sm"
            p="xl"
            w={{ base: "100%", sm: "80%", md: "70%", lg: "50%" }}
          >
            <form onSubmit={onSubmit}>
              <TextInput
                label={t(($) => $.models.attributes.user.email)}
                {...form.getInputProps("email")}
                required
                autoFocus
                enterKeyHint="send"
                {...emailInputProps}
              />
              <Box my="lg" ta="right">
                {t(($) => $.users.new.have_account)}{" "}
                <Box
                  component={Link}
                  fw="bold"
                  href={Routes.new_session_path()}
                >
                  {t(($) => $.users.new.sign_in)}
                </Box>
              </Box>
              <Button fullWidth type="submit" loading={processing}>
                {t(($) => $.helpers.submit.remind_password_form.create)}
              </Button>
            </form>
          </Card>
        </Center>
      </Container>
    </ApplicationLayout>
  );
}
