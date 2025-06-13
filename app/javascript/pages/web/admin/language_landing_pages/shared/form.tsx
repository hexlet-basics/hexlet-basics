import { useTranslation } from "react-i18next";
import { Button } from "@mantine/core";

import {
  XCheck,
  XDynamicInputs,
  XFile,
  XForm,
  XHidden,
  XInput,
  XSelect,
  XTextarea,
} from "@/components/forms";
import type { HTTPVerb } from "use-inertia-form";

import type { Language, LanguageLandingPage } from "@/types";
import type { LanguageLandingPageCrud } from "@/types/serializers";

type Props = {
  data: LanguageLandingPageCrud;
  url: string;
  method?: HTTPVerb;
  languages: Language[];
  landingPages: LanguageLandingPage[];
};

// const locales = [
//   { name: "Russian", code: "ru" },
//   { name: "English", code: "en" },
// ];

export default function Form({ data, landingPages, url, method, languages }: Props) {
  const { t: tHelpers } = useTranslation("helpers");
  return (
      <XForm method={method} model="language_landing_page" data={data} to={url}>
      <XCheck field="main" />
      <XCheck field="listed" />
      <XCheck field="footer" />
      <XSelect
        field="state"
        items={data.meta.state_events}
        labelField="value"
        valueField="key"
      />
      <XSelect
        field="language_id"
        labelField="slug"
        valueField="id"
        items={languages}
      />
      <XSelect
        field="landing_page_to_redirect_id"
        labelField="header"
        valueField="id"
        items={landingPages}
      />
      <XInput field="slug" />
      <XInput field="order" />
      <XInput field="meta_title" />
      <XTextarea field="meta_description" rows={3} />
      <XInput field="name" />
      <XInput field="header" />
      <XTextarea field="description" rows={5} />

      <XInput field="used_in_header" />
      <XTextarea field="used_in_description" rows={5} />

      <XFile metaName="outcomes_image_thumb_url" field="outcomes_image" />
      <XInput field="outcomes_header" />
      <XTextarea field="outcomes_description" rows={5} />

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

      <Button type="submit">
        {tHelpers("submit.save")}
      </Button>
    </XForm>
  );
}
