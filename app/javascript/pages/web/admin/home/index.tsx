import * as Routes from "@/routes.js";
import { DataTable } from 'mantine-datatable';
import type { PropsWithChildren } from 'react';

import { useTranslation } from 'react-i18next';

import AdminLayout from '@/pages/layouts/AdminLayout';
import type { User, Grid } from '@/types/serializers';
import useDataTableProps from '@/hooks/useDataTableProps';
import AppAnchor from '@/components/AppAnchor';
import { Pencil } from "lucide-react";

type Props = PropsWithChildren & {
  admins: User[];
  grid: Grid;
};

export default function Index({ admins, grid }: Props) {
  const { t } = useTranslation();
  const { gridProps } = useDataTableProps<User, {}>(grid);

  return (
    <AdminLayout header={t('admin.home.index.dashboard')}>
      <DataTable
        records={admins}
        columns={[
          { accessor: 'id' },
          { accessor: 'name', sortable: true },
          { accessor: 'email', sortable: true },
          {
            accessor: 'Actions',
            render: (r) => <AppAnchor href={Routes.edit_admin_management_user_path(r)}><Pencil size={14} /></AppAnchor>,
          },
        ]}
        {...gridProps}
      />
    </AdminLayout>
  );
}
