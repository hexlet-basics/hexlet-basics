import { Link } from "@inertiajs/react";
import { Box, Button, Text, TextInput } from "@mantine/core";
import type { PropsWithChildren } from "react";
import { Trans, useTranslation } from "react-i18next";
import { useAppForm } from "@/hooks/useAppForm";
import {
  firstNameInputProps,
  newPasswordInputProps,
  registrationEmailInputProps,
} from "@/lib/authFieldProps";
import * as Routes from "@/routes.js";
import type { UserSignUpForm } from "@/types";

type Props = PropsWithChildren & {
  userDto: UserSignUpForm;
  autoFocus?: boolean;
};

export default function SignUpFormBlock({ userDto, autoFocus = false }: Props) {
  const { t } = useTranslation();

  const payload = userDto;

  const { onSubmit, processing, form, reset } = useAppForm(payload, {
    url: Routes.users_path(),
    method: "post",
    onError: () => {
      reset("password");
    },
  });

  return (
    <form onSubmit={onSubmit}>
      <TextInput
        {...form.getInputProps("first_name")}
        autoFocus={autoFocus}
        {...firstNameInputProps}
      />
      <TextInput
        {...form.getInputProps("email")}
        required
        {...registrationEmailInputProps}
      />
      <TextInput
        {...form.getInputProps("password")}
        required
        {...newPasswordInputProps}
      />
      <Box my="lg" ta="right">
        {t(($) => $.pages.users.new.have_account)}{" "}
        <Text component={Link} fw="bold" href={Routes.new_session_path()} span>
          {t(($) => $.pages.users.new.sign_in)}
        </Text>
      </Box>
      <Button type="submit" fullWidth loading={processing}>
        {t(($) => $.helpers.submit.user_sign_up_form.create)}
      </Button>
      <Text fz="sm" mt="xs">
        <Trans
          t={t}
          i18nKey={($) => $.pages.users.new.confirmation_html}
          components={{
            a: (
              <Text
                component={Link}
                fz="inherit"
                href={Routes.page_path("tos")}
                span
              />
            ),
          }}
        />
      </Text>
    </form>
  );
}
