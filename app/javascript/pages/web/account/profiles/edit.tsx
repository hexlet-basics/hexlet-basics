import { router } from "@inertiajs/react";
import { Box, Button, Card, Center, Container, TextInput } from "@mantine/core";
import type { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import { useAppForm } from "@/hooks/useAppForm";
import useConfirmation from "@/hooks/useConfirmation";
import ApplicationLayout from "@/layouts/ApplicationLayout";
import * as Routes from "@/routes.js";
import type { UserProfileForm } from "@/types/serializers";

type Props = PropsWithChildren & {
  form: UserProfileForm;
};

export default function Edit({ form }: Props) {
  const { t } = useTranslation();
  const confirmDeleting = useConfirmation({
    callback: () => {
      router.visit(Routes.account_profile_path(), { method: "delete" });
    },
  });

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
              <Box
                component="a"
                mt="xl"
                href={Routes.account_profile_path()}
                c="red"
                onClick={confirmDeleting}
              >
                {t(($) => $.account.profiles.edit.delete)}
              </Box>
            </Box>
          </Card>
        </Center>
      </Container>
    </ApplicationLayout>
  );
}
