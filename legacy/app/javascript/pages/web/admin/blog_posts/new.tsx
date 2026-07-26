import { useTranslation } from "react-i18next";
import AdminLayout from "@/layouts/AdminLayout";
import * as Routes from "@/routes.js";
import type BlogPostCreate from "@/types/serializers/BlogPostCreate";
import Form from "./shared/form";
import Menu from "./shared/menu";

type Props = {
  blogPostDto: BlogPostCreate;
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
