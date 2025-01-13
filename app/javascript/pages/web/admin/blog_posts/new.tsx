import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import AdminLayout from "@/pages/layouts/AdminLayout";
import type { BlogPost } from "@/types/serializers";
import Form from "./shared/form";
import { Menu } from "./shared/menu";

type Props = {
  blog_post: BlogPost;
};

export default function New({ blog_post }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout header={t("admin.blog_posts.new.header")}>
      <Menu />
      <Form
        data={blog_post}
        url={Routes.admin_blog_posts_path()}
      />
    </AdminLayout>
  );
}
