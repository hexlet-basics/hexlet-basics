import type { PropsWithChildren } from "react";
import { Col, Container, Form, Row } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type { Q, User } from "@/types/serializers";
import { DataTable } from "primereact/datatable";
import { Column } from "primereact/column";
import { handleOnSort, url } from "@/lib/utils";
import { Link } from "@inertiajs/react";
import useFilterSubmit from "@/hooks/useFilterSubmit";

type Props = PropsWithChildren & {
  admins: User[];
  q: Q;
};

export default function Index({ admins, q }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");

  const handleFilterSubmit = useFilterSubmit({ q });

  const header = (
    <div className="mb-4 bg-body-tertiary p-3">
      <Form action={url()} onSubmit={handleFilterSubmit}>
        <Row>
          <Col>
            <Form.Control
              name="q[fields][email_cont]"
              defaultValue={q.fields.email_cont}
              placeholder="email"
            />
          </Col>
          <Col>
            <Link className="btn btn-outline-primary" href={url()}>
              {tHelpers("reset")}
            </Link>
          </Col>
          {/* <Col> */}
          {/*   <Form.Control placeholder="Last name" /> */}
          {/* </Col> */}
        </Row>
      </Form>
    </div>
  );

  return (
    <ApplicationLayout>
      <Container>
        <DataTable
          header={header}
          sortField={q.sf}
          sortOrder={q.so}
          value={admins}
          onSort={handleOnSort}
          globalFilterFields={["email"]}
        >
          <Column field="id" header="id" />
          <Column field="name" header="name" sortable />
          <Column sortable field="email" header="email" />
        </DataTable>
      </Container>
    </ApplicationLayout>
  );
}
