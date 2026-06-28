import { router } from "@inertiajs/react";
import {
  Box,
  Button,
  Card,
  Center,
  Container,
  Divider,
  Group,
  Stack,
  Text,
  TextInput,
} from "@mantine/core";
import type { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import { useAppForm } from "@/hooks/useAppForm";
import useConfirmation from "@/hooks/useConfirmation";
import ApplicationLayout from "@/layouts/ApplicationLayout";
import { passkeySupported, registerPasskey } from "@/lib/passkey";
import * as Routes from "@/routes.js";
import type { UserCredential, UserProfileForm } from "@/types/serializers";

type Props = PropsWithChildren & {
  form: UserProfileForm;
  passkeys: UserCredential[];
};

export default function Edit({ form, passkeys }: Props) {
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

  const deletePasskey = (id: number) => {
    router.visit(Routes.account_passkey_path(id), { method: "delete" });
  };

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

            <Divider my="xl" />

            <Stack gap="xs">
              <Text fw="bold">
                {t(($) => $.account.profiles.edit.passkeys.title)}
              </Text>
              {passkeys.length === 0 ? (
                <Text c="dimmed" size="sm">
                  {t(($) => $.account.profiles.edit.passkeys.empty)}
                </Text>
              ) : (
                passkeys.map((passkey) => (
                  <Group key={passkey.id} justify="space-between">
                    <Text size="sm">
                      {passkey.nickname ||
                        t(($) => $.account.profiles.edit.passkeys.unnamed)}
                    </Text>
                    <Box
                      component="a"
                      href={Routes.account_passkey_path(passkey.id)}
                      c="red"
                      onClick={(event) => {
                        event.preventDefault();
                        deletePasskey(passkey.id);
                      }}
                    >
                      {t(($) => $.account.profiles.edit.passkeys.delete)}
                    </Box>
                  </Group>
                ))
              )}
              {passkeySupported() ? (
                <Button
                  variant="default"
                  mt="sm"
                  onClick={() => {
                    registerPasskey();
                  }}
                >
                  {t(($) => $.account.profiles.edit.passkeys.add)}
                </Button>
              ) : null}
            </Stack>

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
