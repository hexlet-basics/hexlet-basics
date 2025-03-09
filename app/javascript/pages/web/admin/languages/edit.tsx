import { DTDateTemplate } from "@/components/dtTemplates";
import AdminLayout from "@/pages/layouts/AdminLayout";
import * as Routes from "@/routes.js";
import type {
  LanguageCrud,
  LanguageLandingPage,
  LanguageVersion,
} from "@/types/serializers";
import { Link } from "@inertiajs/react";
import { Column } from "primereact/column";
import { DataTable } from "primereact/datatable";
import { Col, Row } from "react-bootstrap";
import { useTranslation } from "react-i18next";

import Form from "./shared/form";
import { Menu } from "./shared/menu";

type Props = {
  courseDto: LanguageCrud;
  landingPage: LanguageLandingPage;
  courseVersions: LanguageVersion[];
};

export default function Edit({
  courseDto,
  courseVersions,
  landingPage,
}: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout
      header={t("admin.languages.edit.header", { id: courseDto.meta.slug })}
    >
      <Menu data={courseDto} landingPage={landingPage} />
      <Row>
        <Col className="col-4">
          <Form
            method="patch"
            data={courseDto}
            url={Routes.admin_language_path(courseDto.language.id)}
          />
        </Col>
        <Col>
          <div className="d-flex mb-1">
            <h2 className="me-auto">Versions</h2>
            <Link
              method="post"
              className="btn btn-sm btn-outline-secondary"
              href={Routes.admin_language_versions_path(courseDto.language.id)}
            >
              Load New Version
            </Link>
          </div>
          <DataTable value={courseVersions}>
            <Column field="id" header="id" />
            <Column field="result" header="result" />
            <Column
              field="created_at"
              header="created_at"
              body={DTDateTemplate}
            />
          </DataTable>
        </Col>
      </Row>
    </AdminLayout>
  );
}
