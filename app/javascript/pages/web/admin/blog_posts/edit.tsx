import { List } from "@mantine/core";
import { useTranslation } from "react-i18next";
import AdminLayout from "@/layouts/AdminLayout";
import * as Routes from "@/routes.js";
import type { BlogPostCrud, Language } from "@/types/serializers";
import Form from "./shared/form";
import Menu from "./shared/menu";

type Props = {
  blogPostDto: BlogPostCrud;
  relatedCourses: Language[];
};

export default function Edit({ blogPostDto, relatedCourses }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout
      header={t(($) => $.admin.blog_posts.edit.header, {
        id: blogPostDto.name,
      })}
    >
      <Menu data={blogPostDto} />
      <b>Related Courses</b>
      <List mb="xl">
        {relatedCourses.map((course) => {
          return <List.Item key={course.id}>{course.slug}</List.Item>;
        })}
      </List>
      <Form
        method="patch"
        data={blogPostDto}
        url={Routes.admin_blog_post_path(blogPostDto.id)}
      />
    </AdminLayout>
  );
}
