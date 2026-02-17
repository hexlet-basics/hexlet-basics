import { Group } from "@mantine/core";
import { IconEdit, IconFileSearch } from "@tabler/icons-react";
import dayjs from "dayjs";
import { DataTable } from "mantine-datatable";
import type { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import AppAnchor from "@/components/Elements/AppAnchor";
import useDataTableProps from "@/hooks/useDataTableProps";
import AdminLayout from "@/layouts/AdminLayout";
import * as Routes from "@/routes.js";
import type { Grid, Language } from "@/types";
import { Menu } from "./shared/menu";

type Props = PropsWithChildren & {
  courses: Language[];
  grid: Grid;
};

export default function Index({ grid, courses }: Props) {
  const { t } = useTranslation();
  const { gridProps } = useDataTableProps<Language, {}>(grid);

  const renderActions = (item: Language) => (
    <Group gap={5}>
      <AppAnchor href={Routes.edit_admin_language_path(item.id)}>
        <IconEdit size={14} />
      </AppAnchor>
      <AppAnchor
        method="post"
        href={Routes.review_admin_language_path(item.id)}
      >
        <IconFileSearch size={14} />
      </AppAnchor>
    </Group>
  );

  return (
    <AdminLayout header={t(($) => $.admin.languages.index.header)}>
      <Menu />
      <DataTable
        records={courses}
        columns={[
          { accessor: "id" },
          { accessor: "current_version_id", title: "Version Id" },
          { accessor: "slug" },
          { accessor: "progress", sortable: true },
          { accessor: "learn_as", sortable: true },
          { accessor: "order" },
          {
            accessor: "created_at",
            sortable: true,
            render: (r) => dayjs(r.created_at).format("LL"),
          },
          { accessor: "actions", title: "actions", render: renderActions },
        ]}
        {...gridProps}
      />
    </AdminLayout>
  );
}
