import { ActionIcon, Image, Select, ThemeIcon } from '@mantine/core';
import dayjs from 'dayjs';
import { Edit, GraduationCap, Link, Search } from 'lucide-react';
import { DataTable } from 'mantine-datatable';
import type { PropsWithChildren } from 'react';
import { useTranslation } from 'react-i18next';
import AppAnchor from '@/components/Elements/AppAnchor';
import { enums } from '@/generated/enums';
import useDataTableProps from '@/hooks/useDataTableProps';
import AdminLayout from '@/pages/layouts/AdminLayout';
import * as Routes from '@/routes.js';
import type { BlogPost, Grid } from '@/types';
import Menu from './shared/menu';

type Props = PropsWithChildren & {
  blogPosts: BlogPost[];
  grid: Grid & {
    fields: {
      state_eq: string;
    };
  };
};

export default function Index({ grid, blogPosts }: Props) {
  const { t } = useTranslation();
  const { gridProps, filters } = useDataTableProps<
    BlogPost,
    typeof grid.fields
  >(grid);

  const renderActions = (item: BlogPost) => (
    <>
      <AppAnchor me="xs" href={Routes.edit_admin_blog_post_path(item.id)}>
        <ThemeIcon variant="default" size="xs">
          <Edit />
        </ThemeIcon>
      </AppAnchor>
      <AppAnchor me="xs" external href={Routes.blog_post_path(item.slug!)}>
        <ThemeIcon variant="default" size="xs">
          <Link />
        </ThemeIcon>
      </AppAnchor>
      <AppAnchor
        method="post"
        href={Routes.related_courses_admin_blog_post_path(item.id)}
      >
        <ThemeIcon variant="default" size="xs">
          <GraduationCap />
        </ThemeIcon>
      </AppAnchor>
    </>
  );

  const renderCover = (item: BlogPost) => {
    if (!item.cover_thumb_variant) return null;
    return (
      <Image
        src={item.cover_thumb_variant}
        alt={item.name!}
        style={{ maxWidth: '100%', height: 'auto' }}
      />
    );
  };

  const stateFilterSelect = (
    <Select
      data={enums.blogPostState}
      value={filters.values.state_eq}
      onChange={filters.getOnChange('state_eq')}
      leftSection={<Search size={16} />}
      comboboxProps={{ withinPortal: false }}
      clearable
      searchable
    />
  );

  return (
    <AdminLayout header={t('admin.blog_posts.index.header')}>
      <Menu />
      <DataTable
        records={blogPosts}
        columns={[
          { accessor: 'id' },
          { accessor: 'cover', title: 'cover', render: renderCover },
          { accessor: 'name', sortable: true },
          { accessor: 'state', sortable: true, filter: stateFilterSelect },
          {
            accessor: 'related_language_items_count',
            title: 'Related C',
            sortable: true,
          },
          {
            accessor: 'created_at',
            render: (r) => dayjs(r.created_at).format('LL'),
            sortable: true,
          },
          { accessor: 'actions', title: 'actions', render: renderActions },
        ]}
        {...gridProps}
      />
    </AdminLayout>
  );
}
