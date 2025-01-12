import { useTranslation } from "react-i18next";

import type { Review } from "@/types/serializers";
import { Col, Row } from "react-bootstrap";
import { XForm, XInput, XSelect } from "@/components/forms";
import { type HTTPVerb, Submit } from "use-inertia-form";

import * as Routes from "@/routes.js";

type Props = {
  data: Review;
  url: string;
  method?: HTTPVerb;
};

const locales = [
  { name: "Russian", code: "ru" },
  { name: "English", code: "en" },
];

export default function Form({ data, url, method }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");

  return (
    <Row>
      <Col className="col-7">
        <XForm method={method} model="review" data={{ review: data }} to={url}>
          <XSelect
            name="locale"
            labelField="name"
            valueField="code"
            items={locales}
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
