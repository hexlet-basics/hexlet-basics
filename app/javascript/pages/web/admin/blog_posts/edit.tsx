import { useTranslation } from 'react-i18next';
import AdminLayout from '@/pages/layouts/AdminLayout';
import * as Routes from '@/routes.js';
import type { BlogPostCrud } from '@/types/serializers';
import type { Language } from '@/types/serializers';
import Form from './shared/form';
import Menu from './shared/menu';
import {List} from "@mantine/core";

type Props = {
  blogPostDto: BlogPostCrud;
  relatedCourses: Language[];
};

export default function Edit({ blogPostDto, relatedCourses }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout
      header={t('admin.blog_posts.edit.header', {
        id: blogPostDto.data.name,
      })}
    >
      <Menu data={blogPostDto} />
      <b>Related Courses</b>
      <List mb="xl">
        {relatedCourses.map((course) => {
          return (
            <List.Item>
              {course.slug}
            </List.Item>
          );
        })}
      </List>
      <Form
        method="patch"
        data={blogPostDto}
        url={Routes.admin_blog_post_path(blogPostDto.data.id)}
      />
    </AdminLayout>
  );
}
