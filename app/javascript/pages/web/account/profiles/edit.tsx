import { Box, Button, Card, Center, Container, TextInput } from "@mantine/core";
import type { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import AppAnchor from "@/components/Elements/AppAnchor";
import { useAppForm } from "@/hooks/useAppForm";
import ApplicationLayout from "@/layouts/ApplicationLayout";
import * as Routes from "@/routes.js";
import type { UserProfileForm } from "@/types/serializers";

type Props = PropsWithChildren & {
  form: UserProfileForm;
};

export default function Edit({ form }: Props) {
  const { t } = useTranslation();

  const payload = form;

  const {
    onSubmit,
    processing,
    form: appForm,
  } = useAppForm(payload, {
    url: Routes.account_profile_path(),
    method: "patch",
  });

  return (
    <ApplicationLayout center header={t(($) => $.account.profiles.edit.title)}>
      <Container mt="xl">
        <Center>
          <Card
            shadow="sm"
            withBorder
            p="xl"
            w={{ base: "100%", xs: "70%", sm: "50%" }}
          >
            <form onSubmit={onSubmit}>
              <TextInput
                {...appForm.getInputProps("first_name")}
                autoComplete="name"
              />
              <TextInput
                {...appForm.getInputProps("last_name")}
                autoComplete="name"
              />
              <Button type="submit" fullWidth mt="xl" loading={processing}>
                {t(($) => $.helpers.submit.save)}
              </Button>
            </form>

            <Box>
              <AppAnchor
                withConfirmation
                mt="xl"
                href={Routes.account_profile_path()}
                method="delete"
                c="red"
              >
                {t(($) => $.account.profiles.edit.delete)}
              </AppAnchor>
            </Box>
          </Card>
        </Center>
      </Container>
    </ApplicationLayout>
  );
}
