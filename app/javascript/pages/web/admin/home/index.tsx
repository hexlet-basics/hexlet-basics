import { DataTable } from 'mantine-datatable';
import type { PropsWithChildren } from 'react';

import { useTranslation } from 'react-i18next';

import AdminLayout from '@/pages/layouts/AdminLayout';
import type { User, Grid } from '@/types/serializers';
import useDataTableProps from '@/hooks/useDataTableProps';

type Props = PropsWithChildren & {
  admins: User[];
  grid: Grid;
};

export default function Index({ admins, grid }: Props) {
  const { t } = useTranslation();
  const gridProps = useDataTableProps<User>(grid);

  return (
    <AdminLayout header={t('admin.home.index.dashboard')}>
      <DataTable
        records={admins}
        columns={[
          { accessor: 'id' },
          { accessor: 'name', sortable: true },
          { accessor: 'email', sortable: true },
        ]}
        {...gridProps}
      />
    </AdminLayout>
  );
}
