import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import AdminLayout from "@/pages/layouts/AdminLayout";
import type { BlogPostCrud } from "@/types/serializers";
import Form from "./shared/form";
import { Menu } from "./shared/menu";

type Props = {
  blogPostDto: BlogPostCrud;
};

export default function Edit({ blogPostDto }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout
      header={t("admin.blog_posts.edit.header", {
        id: blogPostDto.blog_post.name,
      })}
    >
      <Menu data={blogPostDto} />
      <Form
        method="patch"
        data={blogPostDto}
        url={Routes.admin_blog_post_path(blogPostDto.blog_post.id)}
      />
    </AdminLayout>
  );
}
