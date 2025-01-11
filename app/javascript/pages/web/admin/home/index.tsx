import { type PropsWithChildren } from "react";
import { Col, Form, Row } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import useFilterSubmit from "@/hooks/useFilterSubmit";
import { fieldsToFilters, url } from "@/lib/utils";
import AdminLayout from "@/pages/layouts/AdminLayout";
import type { Grid, User } from "@/types/serializers";
import { Link } from "@inertiajs/react";
import { Column } from "primereact/column";
import { DataTable } from "primereact/datatable";
import useDataTable from "@/hooks/useDataTable";

type Props = PropsWithChildren & {
  admins: User[];
  grid: Grid;
};

export default function Index({ admins, grid }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");

  // const handleFilterSubmit = useFilterSubmit({ q: grid });
  const handleDataTable = useDataTable(grid);

  // const header = (
  //   <div className="mb-4 bg-body-tertiary p-3">
  //     <Form action={url()} onSubmit={handleFilterSubmit}>
  //       <Row>
  //         <Col>
  //           <Form.Control
  //             name="fields[email_cont]"
  //             // defaultValue={grid.fields.email_cont}
  //             placeholder="email"
  //           />
  //         </Col>
  //         <Col>
  //           <Link className="btn btn-outline-primary" href={url()}>
  //             {tHelpers("reset")}
  //           </Link>
  //         </Col>
  //         {/* <Col> */}
  //         {/*   <Form.Control placeholder="Last name" /> */}
  //         {/* </Col> */}
  //       </Row>
  //     </Form>
  //   </div>
  // );

  return (
    <AdminLayout header={t("admin.home.index.dashboard")}>
      <DataTable
        // header={header}
        value={admins}
        filterDisplay="row"
        lazy
        totalRecords={grid.tr}
        rows={grid.per}
        sortField={grid.sf}
        sortOrder={grid.so}
        filters={fieldsToFilters(grid.fields)}
        onSort={handleDataTable}
        onFilter={handleDataTable}
        globalFilterFields={["email"]}
      >
        <Column field="id" header="id" />
        <Column field="name" header="name" sortable />
        <Column sortable filter filterPlaceholder="Search" field="email" header="email" />
      </DataTable>
    </AdminLayout>
  );
}
