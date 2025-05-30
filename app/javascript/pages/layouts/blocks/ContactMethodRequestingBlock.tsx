import XssContent from "@/components/XssContent";
import * as Routes from "@/routes.js";
import { Link } from "@inertiajs/react";
import { Alert, Button } from "react-bootstrap"
import { useTranslation } from "react-i18next"

export default function ContactMethodRequestingBlock() {
  const { t: tLayouts } = useTranslation("layouts")

  return <Alert variant="light" className="small">
    <div className="d-flex">
      <XssContent className="my-auto me-auto">
        {tLayouts('shared.contact_method_requesting.description_html', { url: Routes.edit_account_profile_url() })}
      </XssContent>
      <Link href={Routes.new_lead_path()} className="my-auto text-decoration-none">
        <XssContent>{tLayouts('shared.contact_method_requesting.go')}</XssContent>
      </Link>
    </div>
  </Alert>
}
