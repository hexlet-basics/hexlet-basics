import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import AdminLayout from "@/pages/layouts/AdminLayout";
import type { BlogPost } from "@/types/serializers";
import Form from "./shared/form";
import { Menu } from "./shared/menu";

type Props = {
  blog_post: BlogPost;
};

export default function Edit({ blog_post }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout
      header={t("admin.blog_posts.edit.header", { id: blog_post.id })}
    >
      <Menu data={blog_post} />
      <Form
        method="patch"
        data={blog_post}
        url={Routes.admin_blog_post_path(blog_post.id)}
      />
    </AdminLayout>
  );
}
