import { Link } from "@inertiajs/react";
import { IconEdit, IconTrash } from "@tabler/icons-react";
import dayjs from "dayjs";
import { DataTable } from "mantine-datatable";
import type { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import useDataTableProps from "@/hooks/useDataTableProps";
import AdminLayout from "@/layouts/AdminLayout";
import * as Routes from "@/routes.js";
import type { Banner, Grid } from "@/types/serializers";
import { Menu } from "./shared/menu";

type Props = PropsWithChildren & {
  banners: Banner[];
  grid: Grid;
};

const formatDate = (value: string | null) =>
  value ? dayjs(value).format("LLL") : "—";

export default function Index({ grid, banners }: Props) {
  const { t } = useTranslation();
  const { gridProps } = useDataTableProps<Banner, {}>(grid);

  const renderActions = (item: Banner) => (
    <Link href={Routes.edit_admin_banner_path(item.id)}>
      <IconEdit size={14} />
    </Link>
  );

  const renderDestroy = (item: Banner) => (
    <Link href={Routes.admin_banner_path(item.id)} method="delete" as="button">
      <IconTrash size={14} />
    </Link>
  );

  return (
    <AdminLayout header={t(($) => $.admin.banners.index.header)}>
      <Menu />
      <DataTable
        records={banners}
        columns={[
          { accessor: "id", sortable: true },
          { accessor: "locale", sortable: true },
          { accessor: "state", sortable: true },
          { accessor: "background" },
          {
            accessor: "starts_at",
            sortable: true,
            render: (r) => formatDate(r.starts_at),
          },
          {
            accessor: "finishes_at",
            sortable: true,
            render: (r) => formatDate(r.finishes_at),
          },
          { accessor: "actions", title: "actions", render: renderActions },
          { accessor: "destroy", title: "", render: renderDestroy },
        ]}
        {...gridProps}
      />
    </AdminLayout>
  );
}
