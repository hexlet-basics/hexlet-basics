import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import AdminLayout from "@/pages/layouts/AdminLayout";
import type { LanguageVersion, OriginalLanguage } from "@/types/serializers";
import Form from "./shared/form";
import { Menu } from "./shared/menu";
import { Col, Row } from "react-bootstrap";
import { DTDateTemplate } from "@/components/dtTemplates";
import { DataTable } from "primereact/datatable";
import { Column } from "primereact/column";
import { Link } from "@inertiajs/react";

type Props = {
  originalCourse: OriginalLanguage;
  courseVersions: LanguageVersion[];
};

export default function Edit({ originalCourse, courseVersions }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout
      header={t("admin.languages.edit.header", { id: originalCourse.id })}
    >
      <Menu data={originalCourse} />
      <Row>
        <Col className="col-4">
          <Form
            method="patch"
            data={originalCourse}
            url={Routes.admin_language_path(originalCourse.id)}
          />
        </Col>
        <Col>
          <div className="d-flex mb-1">
            <h2 className="me-auto">Versions</h2>
            <Link
              method="post"
              className="btn btn-sm btn-outline-secondary"
              href={Routes.admin_language_versions_path(originalCourse.id)}
            >Load New Version</Link>
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
