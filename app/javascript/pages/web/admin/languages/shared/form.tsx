import { useTranslation } from "react-i18next";

import type { LanguageCategory, OriginalLanguage } from "@/types/serializers";
import { Col, Row } from "react-bootstrap";
import { XForm, XInput, XSelect } from "@/components/forms";
import { type HTTPVerb, Submit } from "use-inertia-form";

import * as Routes from "@/routes.js";
import { usePage } from "@inertiajs/react";
import type { SharedProps } from "@/types";

type Props = {
  data: OriginalLanguage;
  url: string;
  method?: HTTPVerb;
};

const locales = [
  { name: "Russian", code: "ru" },
  { name: "English", code: "en" },
];

export default function Form({ data, url, method }: Props) {
  const { t } = useTranslation();
  const { courseCategories } = usePage<SharedProps>().props;
  const { t: tHelpers } = useTranslation("helpers");

  return (
    <Row>
      <Col className="col-7">
        <XForm
          method={method}
          model="language"
          data={{ language: data }}
          to={url}
        >
          <XSelect
            name="category_id"
            labelField="name"
            valueField="id"
            items={courseCategories}
          />
          <XInput name="name" />
          <XInput name="slug" />
          <XInput name="state" />

          <Submit className="btn w-100 btn-lg btn-primary mb-3">
            {tHelpers("submit.save")}
          </Submit>
        </XForm>
      </Col>
    </Row>
  );
}
