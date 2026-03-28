import { Link } from "@inertiajs/react";
import { Alert, Button, Group } from "@mantine/core";
import { Trans, useTranslation } from "react-i18next";
import XssContent from "@/components/XssContent";
import * as Routes from "@/routes.js";

export default function ContactMethodRequestingBlock() {
  const { t } = useTranslation();

  return (
    <Alert color="gray" variant="light">
      <Group justify="space-between" align="center">
        <Trans
          t={t}
          i18nKey={($) =>
            $.layouts.shared.contact_method_requesting.description_html
          }
          values={{
            url: Routes.edit_account_profile_path(),
          }}
          components={{
            a: <Link href={Routes.edit_account_profile_path()} />,
          }}
        />
        <Link href={Routes.new_lead_path()} style={{ textDecoration: "none" }}>
          <XssContent>
            {t(($) => $.layouts.shared.contact_method_requesting.go)}
          </XssContent>
        </Link>
      </Group>
    </Alert>
  );
}
