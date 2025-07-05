import * as Routes from "@/routes.js";
import type { User } from "@/types";
import type { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import XssContent from "./XssContent";
import { XForm, XInput } from "./forms";
import { Box, Button, Text } from "@mantine/core";
import AppAnchor from "./AppAnchor";

type Props = PropsWithChildren & {
  user: User;
  autoFocus?: boolean
};

export default function SignUpFormBlock({ user, autoFocus = false }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");

  return (
    <XForm model="user_sign_up_form" data={user} to={Routes.users_path()}>
      <XInput field="first_name" autoComplete="name" autoFocus={autoFocus} />
      <XInput field="email" required autoComplete="email" />
      <XInput field="password" required type="password" autoComplete="current-password" />
      <Box my="lg" ta="right">
        {t("users.new.have_account")}{" "}
        <AppAnchor fw="bold" href={Routes.new_session_path()}>
          {t("users.new.sign_in")}
        </AppAnchor>
      </Box>
      <Button type="submit" fullWidth>
        {tHelpers("submit.user_sign_up_form.create")}
      </Button>

      <XssContent fz="sm" mt="xs">
        {t("users.new.confirmation_html", {
          url: Routes.page_path("tos"),
        })}
      </XssContent>

    </XForm>
  );
}
