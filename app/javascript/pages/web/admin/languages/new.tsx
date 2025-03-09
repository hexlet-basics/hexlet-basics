import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import AdminLayout from "@/pages/layouts/AdminLayout";
import type LanguageCrud from "@/types/serializers/LanguageCrud";
import { Col, Row } from "react-bootstrap";
import Form from "./shared/form";
import { Menu } from "./shared/menu";

type Props = {
  courseDto: LanguageCrud;
};

export default function New({ courseDto }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout header={t("admin.languages.new.header")}>
      <Menu />
      <Row>
        <Col className="col-7">
          <Form data={courseDto} url={Routes.admin_languages_path()} />
        </Col>
      </Row>
    </AdminLayout>
  );
}
