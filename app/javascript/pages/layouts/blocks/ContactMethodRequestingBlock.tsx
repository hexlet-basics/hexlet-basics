import XssContent from "@/components/XssContent";
import * as Routes from "@/routes.js";
import { Alert } from "react-bootstrap"
import { useTranslation } from "react-i18next"

export default function ContactMethodRequestingBlock() {
  const { t: tLayouts } = useTranslation("layouts")

  return <Alert variant="light">
    <XssContent>
      {tLayouts('shared.contact_method_requesting.description_html', { url: Routes.edit_account_profile_url() })}
    </XssContent>
  </Alert>
}
