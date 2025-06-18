import { DataTable } from 'mantine-datatable';
import type { PropsWithChildren } from 'react';

import * as Routes from '@/routes.js';
import { useTranslation } from 'react-i18next';

import AdminLayout from '@/pages/layouts/AdminLayout';
import AppAnchor from '@/components/AppAnchor';
import type { LanguageLandingPage, Grid } from '@/types/serializers';
import { Menu } from './shared/menu';
import useDataTableProps from '@/hooks/useDataTableProps';
import { Edit, Link } from 'lucide-react';
import dayjs from 'dayjs';

type Props = PropsWithChildren & {
  landingPages: LanguageLandingPage[];
  grid: Grid;
};

export default function Index({ grid, landingPages }: Props) {
  const { t } = useTranslation();
  const gridProps = useDataTableProps<LanguageLandingPage>(grid);

  const renderActions = (item: LanguageLandingPage) => (
    <>
      <AppAnchor me="xs" href={Routes.edit_admin_language_landing_page_path(item.id)}>
        <Edit size={14} />
      </AppAnchor>
      <a href={Routes.language_path(item.slug!)} target="_blank">
        <Link size={14} />
      </a>
    </>
  );

  const renderLanguage = (item: LanguageLandingPage) => (
    <AppAnchor href={Routes.edit_admin_language_path(item.language_id)}>
      {item.language_slug}
    </AppAnchor>
  );

  return (
    <AdminLayout header={t('admin.language_landing_pages.index.header')}>
      <Menu />
      <DataTable
        records={landingPages}
        columns={[
          { accessor: 'id' },
          { accessor: 'main' },
          { accessor: 'listed' },
          { accessor: 'state' },
          { accessor: 'order' },
          { accessor: 'header' },
          { accessor: 'language', title: 'language', render: renderLanguage },
          { accessor: 'slug' },
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
