import { Select } from '@mantine/core';
import dayjs from 'dayjs';
import { Edit, Link, Search } from 'lucide-react';
import { DataTable } from 'mantine-datatable';
import type { PropsWithChildren } from 'react';
import { useTranslation } from 'react-i18next';
import AppAnchor from '@/components/AppAnchor';
import useDataTableProps from '@/hooks/useDataTableProps';
import AdminLayout from '@/pages/layouts/AdminLayout';
import * as Routes from '@/routes.js';
import type { Grid, LanguageLandingPage } from '@/types/serializers';
import { Menu } from './shared/menu';

type Props = PropsWithChildren & {
  landingPages: LanguageLandingPage[];
  grid: Grid & {
    fields: {
      state_eq: string;
    };
  };
};

export default function Index({ grid, landingPages }: Props) {
  const { t } = useTranslation();
  const { gridProps, filters } = useDataTableProps<
    LanguageLandingPage,
    typeof grid.fields
  >(grid);

  const renderActions = (item: LanguageLandingPage) => (
    <>
      <AppAnchor
        me="xs"
        href={Routes.edit_admin_language_landing_page_path(item.id)}
      >
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

  const filterState = (
    <Select
      data={['published', 'archived']}
      value={filters.values.state_eq}
      onChange={filters.getOnChange('state_eq')}
      leftSection={<Search size={16} />}
      comboboxProps={{ withinPortal: false }}
      clearable
      searchable
    />
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
          { accessor: 'state', filter: filterState },
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
