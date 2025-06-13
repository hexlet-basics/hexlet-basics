import { DataTable } from 'mantine-datatable';
import type { PropsWithChildren } from 'react';

import * as Routes from '@/routes.js';
import { useTranslation } from 'react-i18next';

import AdminLayout from '@/pages/layouts/AdminLayout';
import AppAnchor from '@/components/AppAnchor';
import type { Survey, Grid } from '@/types';
import { Menu } from './shared/menu';
import useDataTableProps from '@/hooks/useDataTableProps';
import { Edit } from 'lucide-react';

type Props = PropsWithChildren & {
  surveys: Survey[];
  grid: Grid;
};

export default function Index({ grid, surveys }: Props) {
  const { t } = useTranslation();
  const gridProps = useDataTableProps<Survey>(grid);

  const renderActions = (item: Survey) => (
    <AppAnchor href={Routes.edit_admin_survey_path(item.id)}>
      <Edit size={14} />
    </AppAnchor>
  );

  return (
    <AdminLayout header={t('admin.surveys.index.header')}>
      <Menu />
      <DataTable
        records={surveys}
        columns={[
          { accessor: 'id' },
          { accessor: 'slug' },
          { accessor: 'state', sortable: true },
          { accessor: 'question' },
          { accessor: 'created_at', sortable: true },
          { accessor: 'actions', title: 'actions', render: renderActions },
        ]}
        {...gridProps}
      />
    </AdminLayout>
  );
}

