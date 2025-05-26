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

import { Survey, SurveyScenarioCrud, SurveyItemCrud, SurveyScenarioItemCrud } from "@/types";
import { eventNames } from "@/event_names";

type Props = {
  data: SurveyScenarioCrud;
  surveys: Survey[],
  surveysItems: SurveyItemCrud[]
  url: string;
  method?: HTTPVerb;
};

export default function Form({
  data,
  url,
  method,
  surveys,
  surveysItems,
}: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");

  const eventNameOptions: Record<string, string>[] = eventNames.map((event) => ({
    label: event,
    value: event,
  }));

  return (
    <Row>
      <Col className="col-7">
        <XForm method={method} model="survey_scenario" data={data} to={url}>
          <XInput name="name" />
          <XSelect name="survey_item_id" has="survey_item" items={surveysItems} valueField="id" labelField="value_for_select" />
          {/* <XSelect name="event_name" items={eventNameOptions} labelField="label" valueField="value" /> */}
          {/* <XInput name="threshold_count" /> */}

          <XDynamicInputs
            model="triggers"
            label="Triggers"
            emptyData={{ event_name: null, event_threshold_count: null }}
          >
            <XSelect name="event_name" items={eventNameOptions} valueField="value" labelField="label" />
            <XInput name="event_threshold_count" type="number" />
            <XCheck name="_destroy" />
          </XDynamicInputs>

          <XDynamicInputs
            model="items"
            label="Items"
            emptyData={{ survey_id: null }}
          >
            <XSelect name="survey_id" has="survey" valueField="id" labelField="question" items={surveys} />
            <XInput name="order" type="number" />
            {/* <XSelect name="event_name" valueField="id" labelField="name" items={surveys} /> */}
            <XCheck name="_destroy" />
          </XDynamicInputs>

          <Submit className="btn w-100 btn-lg btn-primary mb-3">
            {tHelpers("submit.save")}
          </Submit>
        </XForm>
      </Col>
    </Row>
  );
}

