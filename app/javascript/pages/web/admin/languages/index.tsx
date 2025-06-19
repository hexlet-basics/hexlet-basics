import { DataTable } from 'mantine-datatable';
import type { PropsWithChildren } from 'react';

import * as Routes from '@/routes.js';
import { useTranslation } from 'react-i18next';

import AdminLayout from '@/pages/layouts/AdminLayout';
import AppAnchor from '@/components/AppAnchor';
import type { Language, Grid } from '@/types/serializers';
import { Menu } from './shared/menu';
import useDataTableProps from '@/hooks/useDataTableProps';
import { Edit } from 'lucide-react';
import dayjs from 'dayjs';

type Props = PropsWithChildren & {
  courses: Language[];
  grid: Grid;
};

export default function Index({ grid, courses }: Props) {
  const { t } = useTranslation();
  const { gridProps } = useDataTableProps<Language, {}>(grid);

  const renderActions = (item: Language) => (
    <AppAnchor href={Routes.edit_admin_language_path(item.id)}>
      <Edit size={14} />
    </AppAnchor>
  );

  return (
    <AdminLayout header={t('admin.languages.index.header')}>
      <Menu />
      <DataTable
        records={courses}
        columns={[
          { accessor: 'id' },
          { accessor: 'current_version_id', title: 'Version Id' },
          { accessor: 'slug' },
          { accessor: 'progress', sortable: true },
          { accessor: 'learn_as', sortable: true },
          { accessor: 'order' },
          {
            accessor: 'created_at',
            sortable: true,
            render: (r) => dayjs(r.created_at).format('LL'),
          },
          { accessor: 'actions', title: 'actions', render: renderActions },
        ]}
        {...gridProps}
      />
    </AdminLayout>
  );
}
