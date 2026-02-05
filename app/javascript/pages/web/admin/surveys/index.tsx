import dayjs from 'dayjs';
import { Edit } from 'lucide-react';
import { DataTable } from 'mantine-datatable';
import type { PropsWithChildren } from 'react';
import { useTranslation } from 'react-i18next';
import AppAnchor from '@/components/Elements/AppAnchor';
import useDataTableProps from '@/hooks/useDataTableProps';
import AdminLayout from '@/pages/layouts/AdminLayout';
import * as Routes from '@/routes.js';
import type { Grid, Survey } from '@/types';
import { Menu } from './shared/menu';

type Props = PropsWithChildren & {
  surveys: Survey[];
  grid: Grid;
};

export default function Index({ grid, surveys }: Props) {
  const { t } = useTranslation();
  const { gridProps } = useDataTableProps<Survey, {}>(grid);

  const renderActions = (item: Survey) => (
    <AppAnchor href={Routes.edit_admin_survey_path(item.id)}>
      <Edit size={14} />
    </AppAnchor>
  );

  return (
    <AdminLayout header={t(($) => $.admin.surveys.index.header)}>
      <Menu />
      <DataTable
        records={surveys}
        columns={[
          { accessor: 'id' },
          { accessor: 'slug' },
          { accessor: 'state', sortable: true },
          { accessor: 'question' },
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
