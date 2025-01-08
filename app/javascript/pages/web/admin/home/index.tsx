import type { PropsWithChildren } from "react";
import { Container } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type { User } from "@/types/serializers";
import { DataTable } from "primereact/datatable";
import { Column } from "primereact/column";

type Props = PropsWithChildren & {};

export default function Index() {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");

  const products = [{ code: "jopa", name: "lala" }];

  return (
    <ApplicationLayout>
      <Container>
        <DataTable
          value={products}
          tableStyle={{ minWidth: "50rem" }}
        >
          <Column field="code" header="Code" />
          <Column field="name" header="Name" />
        </DataTable>
      </Container>
    </ApplicationLayout>
  );
}
