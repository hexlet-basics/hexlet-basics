import { Box, Button, Text, TextInput } from "@mantine/core";
import type { PropsWithChildren } from "react";
import { Trans, useTranslation } from "react-i18next";
import { useAppForm } from "@/hooks/useAppForm";
import * as Routes from "@/routes.js";
import type { UserSignUpForm } from "@/types";
import AppAnchor from "./Elements/AppAnchor";

type Props = PropsWithChildren & {
  userDto: UserSignUpForm;
  autoFocus?: boolean;
};

export default function SignUpFormBlock({ userDto, autoFocus = false }: Props) {
  const { t } = useTranslation();

  const payload = userDto;

  const { onSubmit, processing, form } = useAppForm(payload, {
    url: Routes.users_path(),
    method: "post",
  });

  return (
    <form onSubmit={onSubmit}>
      <TextInput
        {...form.getInputProps("first_name")}
        autoComplete="name"
        autoFocus={autoFocus}
      />
      <TextInput
        {...form.getInputProps("email")}
        required
        autoComplete="email"
      />
      <TextInput
        {...form.getInputProps("password")}
        required
        type="password"
        autoComplete="current-password"
      />
      <Box my="lg" ta="right">
        {t(($) => $.pages.users.new.have_account)}{" "}
        <AppAnchor fw="bold" href={Routes.new_session_path()}>
          {t(($) => $.pages.users.new.sign_in)}
        </AppAnchor>
      </Box>
      <Button type="submit" fullWidth loading={processing}>
        {t(($) => $.helpers.submit.user_sign_up_form.create)}
      </Button>
      <Text fz="sm" mt="xs">
        <Trans
          t={t}
          i18nKey={($) => $.pages.users.new.confirmation_html}
          components={{
            a: <AppAnchor fz="inherit" href={Routes.page_path("tos")} />,
          }}
        />
      </Text>
    </form>
  );
}
