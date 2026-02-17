import { Link } from "@inertiajs/react";
import { Alert, Button, Group } from "@mantine/core";
import { useTranslation } from "react-i18next";
import XssContent from "@/components/XssContent";
import * as Routes from "@/routes.js";

export default function ContactMethodRequestingBlock() {
  const { t } = useTranslation();

  return (
    <Alert color="gray" variant="light">
      <Group justify="space-between" align="center">
        <XssContent>
          {t(
            ($) => $.layouts.shared.contact_method_requesting.description_html,
            {
              url: Routes.edit_account_profile_url(),
            },
          )}
        </XssContent>
        <Link href={Routes.new_lead_path()} style={{ textDecoration: "none" }}>
          <XssContent>
            {t(($) => $.layouts.shared.contact_method_requesting.go)}
          </XssContent>
        </Link>
      </Group>
    </Alert>
  );
}
