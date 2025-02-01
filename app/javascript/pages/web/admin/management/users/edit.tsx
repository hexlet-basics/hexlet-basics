import type { PropsWithChildren } from "react";

import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import { XCheck, XForm, XInput } from "@/components/forms";
import AdminLayout from "@/pages/layouts/AdminLayout";
import type { User } from "@/types/serializers";
import { Col, Row } from "react-bootstrap";
import { Submit } from "use-inertia-form";
import { Menu } from "./shared/menu";
import { DataTable } from "primereact/datatable";
import { Column } from "primereact/column";

type Props = PropsWithChildren & {
  user: User;
  progress: Record<string, number>[]
};

export default function Edit({ user, progress }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");

  return (
    <AdminLayout header={t("admin.management.users.edit.header")}>
      <Menu data={user} />
      <Row>
        <Col className="col-5">
          <XForm
            method="patch"
            model="user"
            data={{ user }}
            to={Routes.admin_management_user_path(user.id)}
          >
            <XInput name="email" autoComplete="email" disabled />
            <XCheck name="admin" />
            <XInput name="first_name" autoComplete="name" />
            <XInput name="last_name" autoComplete="name" />
            <XInput name="last_name" autoComplete="name" />
            <Submit className="btn w-100 btn-lg btn-primary mb-3">
              {tHelpers("submit.save")}
            </Submit>
          </XForm>
        </Col>
        <Col>
          <div className="d-flex mb-1">
            <h2 className="me-auto">Progress</h2>
          </div>
          <DataTable value={progress}>
            <Column field="language" header="Language" />
            <Column field="count" header="Count" />
          </DataTable>
        </Col>
      </Row>
    </AdminLayout>
  );
}
