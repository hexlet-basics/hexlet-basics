import { useTranslation } from "react-i18next";
import { Button } from '@mantine/core';

import { XCheck, XDynamicInputs, XForm, XInput, XSelect, XTextarea } from "@/components/forms";
import type { HTTPVerb } from "use-inertia-form";

import type { SurveyCrud, SurveyItemCrud } from "@/types";

type Props = {
  data: SurveyCrud;
  surveyItems: SurveyItemCrud[]
  url: string;
  method?: HTTPVerb;
};

export default function Form({ data, url, method, surveyItems }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");

  return (
    <XForm method={method} model="survey" data={data} to={url}>
      <XInput field="question" />
      <XInput field="slug" />
      <XTextarea field="description" rows={8} />

      <XDynamicInputs
        model="items"
        label="Items"
        emptyData={{ value: '', order: 100 }}
      >
        <XInput field="value" />
        <XInput field="tag_list" />
        <XSelect field="state" valueField="value" labelField="key" items={data.meta.item_states} />
        <XInput field="order" />
      </XDynamicInputs>

      <Button type="submit">
        {tHelpers("submit.save")}
      </Button>
    </XForm>
  );
}

