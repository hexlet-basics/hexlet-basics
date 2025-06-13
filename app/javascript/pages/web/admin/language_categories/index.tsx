import { DataTable } from 'mantine-datatable';
import type { PropsWithChildren } from "react";

import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import AdminLayout from "@/pages/layouts/AdminLayout";
import type {
  Grid,
  LanguageCategory,
} from "@/types";
import { Menu } from "./shared/menu";
import useConfirmation from "@/hooks/useConfirmation";
import AppAnchor from '@/components/AppAnchor';
import useDataTableProps from '@/hooks/useDataTableProps';
import { Edit, Link } from 'lucide-react';

type Props = PropsWithChildren & {
  categories: LanguageCategory[];
  grid: Grid;
};

export default function Index({ grid, categories }: Props) {
  const { t } = useTranslation();
  const confirmDeleting = useConfirmation();

  const gridProps = useDataTableProps<LanguageCategory>(grid)

  const renderActions = (item: LanguageCategory) => {
    return (
      <>
        <AppAnchor me="xs" href={Routes.edit_admin_language_category_path(item.id)}>
          <Edit size={14} />
        </AppAnchor>
        <a
          target="_blank"
          href={Routes.language_category_path(item.slug!)}
        >
          <Link size={14} />
        </a>
        {/* <Link */}
        {/*   onClick={confirmDeleting} */}
        {/*   className="btn btn-link link-body-emphasis p-0 m-0" */}
        {/*   method="delete" */}
        {/*   href={Routes.admin_language_category_path(data.id)} */}
        {/* > */}
        {/*   {<i className="bi bi-file-x" />} */}
        {/* </Link> */}
      </>
    );
  };

  return (
    <AdminLayout header={t("admin.language_categories.index.header")}>
      <Menu />
      <DataTable
        records={categories}
        columns={[
          { accessor: 'id' },
          { accessor: 'name', sortable: true },
          { accessor: 'slug', },
          {
            accessor: 'actions',
            title: 'actions',
            render: renderActions
          }
        ]}
        {...gridProps}
      />
    </AdminLayout>
  );
}
