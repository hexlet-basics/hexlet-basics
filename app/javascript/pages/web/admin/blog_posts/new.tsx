import { useTranslation } from 'react-i18next';
import AdminLayout from '@/pages/layouts/AdminLayout';
import * as Routes from '@/routes.js';
import type BlogPostCrud from '@/types/serializers/BlogPostCrud';
import Form from './shared/form';
import Menu from './shared/menu';

type Props = {
  blogPostDto: BlogPostCrud;
};

export default function New({ blogPostDto }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout header={t(($) => $.admin.blog_posts.new.header)}>
      <Menu />
      <Form data={blogPostDto} url={Routes.admin_blog_posts_path()} />
    </AdminLayout>
  );
}
