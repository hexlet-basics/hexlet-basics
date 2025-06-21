import { DataTable } from 'mantine-datatable';
import type { PropsWithChildren } from 'react';

import * as Routes from '@/routes.js';
import { useTranslation } from 'react-i18next';

import AdminLayout from '@/pages/layouts/AdminLayout';
import AppAnchor from '@/components/AppAnchor';
import type { User, Grid } from '@/types/serializers';
import { Menu } from './shared/menu';
import useDataTableProps from '@/hooks/useDataTableProps';
import { Edit } from 'lucide-react';
import { XForm, XInput } from "@/components/forms";
import dayjs from 'dayjs';

type Props = PropsWithChildren & {
  users: User[];
  grid: Grid;
};

export default function Index({ grid, users }: Props) {
  const { t } = useTranslation();
  const { gridProps } = useDataTableProps<User, {}>(grid);

  const renderActions = (item: User) => (
    <AppAnchor href={Routes.edit_admin_management_user_path(item.id)}>
      <Edit size={14} />
    </AppAnchor>
  );

  return (
    <AdminLayout header={t('admin.management.users.index.header')}>
      <Menu />
      <DataTable
        records={users}
        columns={[
          { accessor: 'id' },
          { accessor: 'admin', title: 'Admin?' },
          { accessor: 'assistant_messages_count', title: 'Messages' },
          { accessor: 'name', sortable: true },
          {
            accessor: 'email',
            sortable: true,
            filter: (
              <XForm method="get" to={Routes.admin_management_users_path()}>
                <XInput field="fields[email_cont]" label={t('admin.management.users.index.search_by_email')} />
              </XForm>
            )
          },
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
