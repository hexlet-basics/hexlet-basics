import { ActionIcon } from "@mantine/core";
import { IconEdit, IconLink } from "@tabler/icons-react";
import dayjs from "dayjs";
import { DataTable } from "mantine-datatable";
import type { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import AppAnchor from "@/components/Elements/AppAnchor";
// import useConfirmation from '@/hooks/useConfirmation';
import useDataTableProps from "@/hooks/useDataTableProps";
import AdminLayout from "@/layouts/AdminLayout";
import * as Routes from "@/routes.js";
import type { Grid, LanguageCategory } from "@/types";
import { Menu } from "./shared/menu";

type Props = PropsWithChildren & {
  categories: LanguageCategory[];
  grid: Grid;
};

export default function Index({ grid, categories }: Props) {
  const { t } = useTranslation();
  // const confirmDeleting = useConfirmation();

  const { gridProps } = useDataTableProps<LanguageCategory, {}>(grid);

  const renderActions = (item: LanguageCategory) => {
    return (
      <>
        <AppAnchor
          me="xs"
          href={Routes.edit_admin_language_category_path(item.id)}
        >
          <ActionIcon variant="default" size="xs">
            <IconEdit />
          </ActionIcon>
        </AppAnchor>
        <AppAnchor external href={Routes.language_category_path(item.slug!)}>
          <ActionIcon variant="default" size="xs">
            <IconLink />
          </ActionIcon>
        </AppAnchor>
        {/* <IconLink */}
        {/*   onClick={confirmDeleting} */}
        {/*   className="btn btn-link link-body-emphasis p-0 m-0" */}
        {/*   method="delete" */}
        {/*   href={Routes.admin_language_category_path(data.id)} */}
        {/* > */}
        {/*   {<IconFileX />} */}
        {/* </IconLink> */}
      </>
    );
  };

  return (
    <AdminLayout header={t(($) => $.admin.language_categories.index.header)}>
      <Menu />
      <DataTable
        records={categories}
        columns={[
          { accessor: "id" },
          { accessor: "name", sortable: true },
          { accessor: "slug" },
          {
            accessor: "created_at",
            sortable: true,
            render: (r) => dayjs(r.created_at).format("LL"),
          },
          {
            accessor: "actions",
            title: "actions",
            render: renderActions,
          },
        ]}
        {...gridProps}
      />
    </AdminLayout>
  );
}
