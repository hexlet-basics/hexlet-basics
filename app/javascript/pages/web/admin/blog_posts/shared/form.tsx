import { useTranslation } from "react-i18next";

import { XFile, XForm, XInput, XSelect, XStateEvent } from "@/components/forms";
import type { BlogPostCrud } from "@/types/serializers";
import { Col, Row } from "react-bootstrap";
import { type HTTPVerb, Submit } from "use-inertia-form";

type Props = {
  data: BlogPostCrud;
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
        <XForm method={method} model="blog_post" data={data} to={url}>
          <XStateEvent fieldName="state" />
          <XInput name="name" />
          <XFile metaName="cover_thumb_variant" name="cover" />
          <XInput name="slug" />
          <XInput
            name="description"
            as="textarea"
            style={{ height: "100px" }}
          />
          <XInput name="body" as="textarea" style={{ height: "500px" }} />
          <Submit className="btn w-100 btn-lg btn-primary mb-3">
            {tHelpers("submit.save")}
          </Submit>
        </XForm>
      </Col>
    </Row>
  );
}
