import { Submit } from "use-inertia-form";
import { XForm, XSelect, XInput } from "./forms";
import { useTranslation } from "react-i18next";
import { enumToOptions } from "@/lib/utils";
import { LeadCrud } from "@/types";
import * as Routes from "@/routes.js";
import XssContent from "./XssContent";

type Props = {
  lead: LeadCrud
}

export default function LeadFormBlock({ lead }: Props) {
  const { t: tAr } = useTranslation("activerecord");
  const { t: tHelpers } = useTranslation("helpers");
  const { t: tViews } = useTranslation();

  const contactMethodEnum = tAr('attributes.user.contact_method/values', { returnObjects: true })
  const contactMethodOptions = enumToOptions(contactMethodEnum);

  return (
    <XForm className="d-flex flex-column h-100" model="lead" data={lead} to={Routes.leads_path()}>
      <XSelect
        name="contact_method"
        labelField="name"
        valueField="id"
        items={contactMethodOptions}
      />
      <XInput name="contact_value" />
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
