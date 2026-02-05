import dayjs from 'dayjs';
import { Edit } from 'lucide-react';
import { DataTable } from 'mantine-datatable';
import type { PropsWithChildren } from 'react';
import { useTranslation } from 'react-i18next';
import AppAnchor from '@/components/Elements/AppAnchor';
import useDataTableProps from '@/hooks/useDataTableProps';
import AdminLayout from '@/pages/layouts/AdminLayout';
import * as Routes from '@/routes.js';
import type { Grid, Review } from '@/types/serializers';
import { Menu } from './shared/menu';

type Props = PropsWithChildren & {
  reviews: Review[];
  grid: Grid;
};

export default function Index({ grid, reviews }: Props) {
  const { t } = useTranslation();
  const { gridProps } = useDataTableProps<Review, {}>(grid);

  const renderActions = (item: Review) => (
    <AppAnchor href={Routes.edit_admin_review_path(item.id)}>
      <Edit size={14} />
    </AppAnchor>
  );

  return (
    <AdminLayout header={t(($) => $.admin.reviews.index.header)}>
      <Menu />
      <DataTable
        records={reviews}
        columns={[
          { accessor: 'id', sortable: true },
          { accessor: 'locale' },
          { accessor: 'user_id', title: 'user' },
          { accessor: 'first_name', title: 'First Name', sortable: true },
          { accessor: 'last_name', title: 'Last Name', sortable: true },
          { accessor: 'state', sortable: true },
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
