import type { PropsWithChildren } from "react";

import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import { XCheck, XForm, XInput } from "@/components/forms";
import AdminLayout from "@/pages/layouts/AdminLayout";
import type { User } from "@/types/serializers";
import { Col, Row } from "react-bootstrap";
import { Submit } from "use-inertia-form";
import { Menu } from "./shared/menu";

type Props = PropsWithChildren & {
  user: User;
};

export default function Edit({ user }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");

  return (
    <AdminLayout header={t("admin.management.users.edit.header")}>
      <Menu data={user} />
      <Row>
        <Col className="col-7">
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
      </Row>
    </AdminLayout>
  );
}
