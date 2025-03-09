import { DTDateTemplate } from "@/components/dtTemplates";
import AdminLayout from "@/pages/layouts/AdminLayout";
import * as Routes from "@/routes.js";
import type { Language, LanguageVersion } from "@/types/serializers";
import { Link } from "@inertiajs/react";
import { Column } from "primereact/column";
import { DataTable } from "primereact/datatable";
import { Col, Row } from "react-bootstrap";
import { useTranslation } from "react-i18next";

import Form from "./shared/form";
import { Menu } from "./shared/menu";

type Props = {
  course: Language;
  courseVersions: LanguageVersion[];
};

export default function Edit({ course, courseVersions }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout header={t("admin.languages.edit.header", { id: course.id })}>
      <Menu data={course} />
      <Row>
        <Col className="col-4">
          <Form
            method="patch"
            data={course}
            url={Routes.admin_language_path(course.id)}
          />
        </Col>
        <Col>
          <div className="d-flex mb-1">
            <h2 className="me-auto">Versions</h2>
            <Link
              method="post"
              className="btn btn-sm btn-outline-secondary"
              href={Routes.admin_language_versions_path(course.id)}
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
