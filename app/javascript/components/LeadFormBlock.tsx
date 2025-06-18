import { Submit } from "use-inertia-form";
import { XForm, XSelect, XInput, XHidden } from "./forms";
import { useTranslation } from "react-i18next";
import { enumToOptions, fromWindow } from "@/lib/utils";
import { LeadCrud } from "@/types";
import * as Routes from "@/routes.js";
import XssContent from "./XssContent";
import { toMerged } from "es-toolkit";
import { Box, Button } from "@mantine/core";

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

  const data = toMerged(
    lead,
    { lead: {
      ym_client_id: fromWindow('ymClientId'),
      contact_method: 'telegram',
    } }
  )

  return (
    <XForm model="lead" data={data} to={Routes.leads_path()}>
      <XHidden field="ym_client_id" />
      <XSelect
        required
        field="contact_method"
        labelField="name"
        valueField="id"
        items={contactMethodOptions}
      />
      <XInput required autoFocus={autoFocus} field="contact_value" />
      <Box fz={14} my="lg">
        <Box>
          <XssContent>
            {tViews('blocks.lead_form_block.description1')}
          </XssContent>
        </Box>
        <Box>
          <XssContent>
            {tViews('blocks.lead_form_block.description2')}
          </XssContent>
        </Box>
      </Box>
      <Button type="submit" fullWidth>
        {tHelpers("send")}
      </Button>
    </XForm>
  )
}
