import { useTranslation } from "react-i18next";

import { XCheck, XDynamicInputs, XForm, XHidden, XInput, XSelect, XTextarea } from "@/components/forms";
import { type HTTPVerb } from "use-inertia-form";

import type { LanguageCategoryCrud } from "@/types";
import { Button } from "@mantine/core";

type Props = {
  data: LanguageCategoryCrud;
  url: string;
  method?: HTTPVerb;
};

export default function Form({ data, url, method }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");

  return (
    <XForm method={method} model="language_category" data={data} to={url}>
      <XInput autoFocus field="name" />
      <XInput field="header" />
      <XInput field="slug" />
      <XTextarea field="description" rows={5} />

      <XDynamicInputs
        model="items"
        label="Items"
        emptyData={{ language_landing_page_id: null }}
      >
        <XHidden field="id" />
        <XSelect
          field="language_landing_page_id"
          valueField="id"
          labelField="header"
          items={data.meta.landingPagesForCategories}
        />
        <XCheck field="_destroy" />
      </XDynamicInputs>

      <XDynamicInputs
        model="qna_items"
        label="QNA"
        emptyData={{ question: "", answer: "" }}
      >
        <XHidden field="id" />
        <XInput field="question" />
        <XTextarea field="answer" rows={5} />
        <XCheck field="_destroy" />
      </XDynamicInputs>

      <Button type="submit" mt="xl">
        {tHelpers("submit.save")}
      </Button>
    </XForm>
  );
}
