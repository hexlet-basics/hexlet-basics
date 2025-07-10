import { Image } from '@mantine/core';
import dayjs from 'dayjs';
import { Edit, Link } from 'lucide-react';
import { DataTable } from 'mantine-datatable';
import type { PropsWithChildren } from 'react';
import { useTranslation } from 'react-i18next';
import AppAnchor from '@/components/AppAnchor';
import useDataTableProps from '@/hooks/useDataTableProps';
import AdminLayout from '@/pages/layouts/AdminLayout';
import * as Routes from '@/routes.js';
import type { BlogPost, Grid } from '@/types';
import Menu from './shared/menu';

type Props = PropsWithChildren & {
  blogPosts: BlogPost[];
  grid: Grid;
};

export default function Index({ grid, blogPosts }: Props) {
  const { t } = useTranslation();
  const { gridProps } = useDataTableProps<BlogPost, {}>(grid);

  const renderActions = (item: BlogPost) => (
    <>
      <AppAnchor me="xs" href={Routes.edit_admin_blog_post_path(item.id)}>
        <Edit size={14} />
      </AppAnchor>
      <a target="_blank" href={Routes.blog_post_path(item.slug!)}>
        <Link size={14} />
      </a>
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

  return (
    <AdminLayout header={t('admin.blog_posts.index.header')}>
      <Menu />
      <DataTable
        records={blogPosts}
        columns={[
          { accessor: 'id' },
          { accessor: 'cover', title: 'cover', render: renderCover },
          { accessor: 'name', sortable: true },
          { accessor: 'state', sortable: true },
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
