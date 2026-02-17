import { Box, Button, Select, TextInput } from "@mantine/core";
import { toMerged } from "es-toolkit";
import { Trans, useTranslation } from "react-i18next";
import { useAppForm } from "@/hooks/useAppForm";
import { arrayToSelectData, enumToOptions, fromWindow } from "@/lib/utils";
import * as Routes from "@/routes.js";
import type { LeadCrud } from "@/types";
import AppAnchor from "./Elements/AppAnchor";

type Props = {
  leadDto: LeadCrud;
  autoFocus?: boolean;
};

export default function LeadFormBlock({ leadDto, autoFocus = false }: Props) {
  const { t } = useTranslation();

  const contactMethodEnum = t(
    ($) => $.models.attributes.user["contact_method/values"],
    {
      returnObjects: true,
    },
  );
  const contactMethodOptions = enumToOptions(contactMethodEnum);

  const payload = toMerged(leadDto, {
    ym_client_id: fromWindow("ymClientId"),
    contact_method: "telegram",
  });
  const contactMethodSelectData = arrayToSelectData(
    contactMethodOptions,
    "id",
    "name",
  );

  const { onSubmit, processing, form } = useAppForm(payload, {
    url: Routes.leads_path(),
    method: "post",
  });

  return (
    <form onSubmit={onSubmit}>
      {/* Поле ym_client_id передаём скрыто */}
      <input type="hidden" {...form.getInputProps("ym_client_id")} />
      <Select
        {...form.getSelectProps("contact_method", contactMethodSelectData)}
        required
      />
      <TextInput
        {...form.getInputProps("contact_value")}
        required
        autoFocus={autoFocus}
      />
      <Box fz="sm" my="lg">
        {t(($) => $.blocks.lead_form_block.description1)}
        <Trans
          t={t}
          i18nKey={($) => $.blocks.lead_form_block.description2}
          components={{
            a: <AppAnchor external href="https://t.me/WelcomeCodebasicsBot" />,
          }}
        />
      </Box>
      <Button type="submit" fullWidth loading={processing}>
        {t(($) => $.helpers.send)}
      </Button>
    </form>
  );
}
