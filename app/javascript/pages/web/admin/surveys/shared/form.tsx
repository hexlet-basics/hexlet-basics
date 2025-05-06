import { useTranslation } from "react-i18next";

import {
    XCheck,
  XDynamicInputs,
  XForm,
  XInput,
  XSelect
} from "@/components/forms";
import { Col, Row } from "react-bootstrap";
import { type HTTPVerb, Submit } from "use-inertia-form";

import { SurveyCrud, SurveyItem, SurveyItemCrud } from "@/types";

type Props = {
  data: SurveyCrud;
  surveysItems: SurveyItemCrud[]
  url: string;
  method?: HTTPVerb;
};

export default function Form({ data, url, method, surveysItems }: Props) {
  const { t } = useTranslation();
  console.log(surveysItems)
  const { t: tHelpers } = useTranslation("helpers");

  return (
    <Row>
      <Col className="col-7">
        <XForm method={method} model="survey" data={data} to={url}>
          {/* <XStateEvent fieldName="state" /> */}
          <XInput name="question" />
          <XInput name="slug" />
          <XSelect name="parent_survey_item_id" items={surveysItems} valueField="id" labelField="value_for_select" />
          <XInput name="description" as="textarea" style={{ height: "200px" }} />

          <XDynamicInputs
            model="items"
            label="Items"
            emptyData={{ value: '', order: 100 }}
          >
            <XInput name="value" />
            <XInput name="slug" />
            <XSelect name="state" valueField="value" labelField="key" items={data.meta.item_states} />
            <XInput name="order" />
            <XCheck name="run_always" />
            <XInput name="run_after_finishing_lessons_count" />
            {/* <XCheck name="_destroy" /> */}
          </XDynamicInputs>

          <Submit className="btn w-100 btn-lg btn-primary mb-3">
            {tHelpers("submit.save")}
          </Submit>
        </XForm>
      </Col>
    </Row>
  );
}

