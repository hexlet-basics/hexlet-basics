import { useTranslation } from "react-i18next";

import { XForm, XInput, XSelect, XStateEvent } from "@/components/forms";
import type { OriginalLanguage, Review } from "@/types/serializers";
import { Col, Row } from "react-bootstrap";
import { type HTTPVerb, Submit } from "use-inertia-form";

import * as Routes from "@/routes.js";
import { locales } from "@/lib/utils";

type Props = {
  data: Review;
  url: string;
  courses: OriginalLanguage[];
  method?: HTTPVerb;
};

export default function Form({ courses, data, url, method }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");

  return (
    <Row>
      <Col className="col-7">
        <XForm method={method} model="review" data={{ review: data }} to={url}>
          <XStateEvent
            fieldName="state"
          />
          <XSelect
            name="locale"
            labelField="name"
            valueField="code"
            items={locales}
          />
          <XSelect
            name="language_id"
            has="language"
            labelField="name"
            valueField="id"
            items={courses}
          />
          <XSelect
            name="user_id"
            has="user"
            labelField="email"
            valueField="id"
            source={Routes.search_admin_api_users_path()}
          />
          <XInput name="first_name" autoComplete="name" />
          <XInput name="last_name" autoComplete="name" />
          <XInput name="body" as="textarea" style={{ height: "200px" }} />
          <Submit className="btn w-100 btn-lg btn-primary mb-3">
            {tHelpers("submit.save")}
          </Submit>
        </XForm>
      </Col>
    </Row>
  );
}
