import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import AdminLayout from "@/pages/layouts/AdminLayout";
import type { Language } from "@/types/serializers";
import { Col, Row } from "react-bootstrap";
import Form from "./shared/form";
import { Menu } from "./shared/menu";

type Props = {
  course: Language;
};

export default function New({ course }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout header={t("admin.languages.new.header")}>
      <Menu />
      <Row>
        <Col className="col-7">
          <Form data={course} url={Routes.admin_languages_path()} />
        </Col>
      </Row>
    </AdminLayout>
  );
}
