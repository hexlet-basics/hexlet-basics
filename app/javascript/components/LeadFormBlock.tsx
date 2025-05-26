import { Submit } from "use-inertia-form";
import { XForm, XSelect, XInput, XHidden } from "./forms";
import { useTranslation } from "react-i18next";
import { enumToOptions } from "@/lib/utils";
import { LeadCrud } from "@/types";
import * as Routes from "@/routes.js";
import XssContent from "./XssContent";

type Props = {
  lead: LeadCrud
  autoFocus?: boolean
}

export default function LeadFormBlock({ lead, autoFocus = false }: Props) {
  const { t: tAr } = useTranslation("activerecord");
  const { t: tHelpers } = useTranslation("helpers");
  const { t: tViews } = useTranslation();

  const contactMethodEnum = tAr('attributes.user.contact_method/values', { returnObjects: true })
  const contactMethodOptions = enumToOptions(contactMethodEnum);

  const data = {
    lead: {
      ...lead?.lead,
      ym_client_id: window.ymClientId
    }
  }

  return (
    <XForm className="d-flex flex-column h-100" model="lead" data={data} to={Routes.leads_path()}>
      <XHidden name="ym_client_id" />
      <XSelect
        name="contact_method"
        labelField="name"
        valueField="id"
        items={contactMethodOptions}
      />
      <XInput autoFocus={autoFocus} name="contact_value" />
      <div className="small text-end mb-5">
        <div>
          <XssContent>
            {tViews('blocks.lead_form_block.description1')}
          </XssContent>
        </div>
        <div>
          <XssContent>
            {tViews('blocks.lead_form_block.description2')}
          </XssContent>
        </div>
      </div>
      <Submit className="btn d-block btn-lg btn-primary mt-auto shadow-sm">
        {tHelpers("send")}
      </Submit>
    </XForm>
  )
}
