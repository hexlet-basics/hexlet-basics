import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import AdminLayout from "@/pages/layouts/AdminLayout";
import type { Language, LanguageLandingPageCrud } from "@/types/serializers";
import { Col, Row } from "react-bootstrap";
import Form from "./shared/form";
import { Menu } from "./shared/menu";

type Props = {
  landingPageDto: LanguageLandingPageCrud;
  languages: Language[];
};

export default function New({ landingPageDto, languages }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout header={t("admin.language_landing_pages.new.header")}>
      <Menu />
      <Row>
        <Col className="col-7">
          <Form
            data={landingPageDto}
            languages={languages}
            url={Routes.admin_language_landing_pages_path()}
          />
        </Col>
      </Row>
    </AdminLayout>
  );
}
