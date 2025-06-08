import { useTranslation } from "react-i18next";

import { XCheck, XDynamicInputs, XForm, XHidden, XInput, XSelect } from "@/components/forms";
import { type HTTPVerb, Submit } from "use-inertia-form";

import type { SharedProps } from "@/types";
import type { LanguageCategoryCrud } from "@/types/serializers";
import { usePage } from "@inertiajs/react";

type Props = {
  data: LanguageCategoryCrud;
  url: string;
  method?: HTTPVerb;
};

// const locales = [
//   { name: "Russian", code: "ru" },
//   { name: "English", code: "en" },
// ];

export default function Form({ data, url, method }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");
  // const { t: tEnums } = useTranslation("enumerize");
  // const languageProgressEnum = tEnums("language.progress", {
  //   returnObjects: true,
  // });
  // const languageProgressEnumOptions = enumToOptions(languageProgressEnum);
  //
  // const languageLearnAsEnum = tEnums("language.learn_as", {
  //   returnObjects: true,
  // });
  // const languageLearnAsEnumOptions = enumToOptions(languageLearnAsEnum);

  return (
    <XForm method={method} model="language_category" data={data} to={url}>
      <XInput name="name" />
      <XInput name="header" />
      <XInput name="slug" />
      <XInput as="textarea" style={{ height: "100px" }} name="description" />

      <XDynamicInputs
        model="items"
        label="Items"
        emptyData={{ language_landing_page_id: null }}
      >
        <XHidden type="hidden" name="id" />
        <XSelect name="language_landing_page_id" valueField="id" labelField="header" items={data.meta.landingPagesForCategories} />
        <XCheck name="_destroy" />
      </XDynamicInputs>

      <XDynamicInputs
        model="qna_items"
        label="QNA"
        emptyData={{ question: "", answer: "" }}
      >
        <XHidden type="hidden" name="id" />
        <XInput name="question" />
        <XInput name="answer" />
        <XCheck name="_destroy" />
      </XDynamicInputs>


      <Submit className="btn w-100 btn-lg btn-primary mb-3">
        {tHelpers("submit.save")}
      </Submit>
    </XForm>
  );
}
