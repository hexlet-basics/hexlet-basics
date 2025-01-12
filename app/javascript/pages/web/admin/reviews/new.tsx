import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import AdminLayout from "@/pages/layouts/AdminLayout";
import type { Review } from "@/types/serializers";
import Form from "./shared/form";
import { Menu } from "./shared/menu";

type Props = {
  review: Review;
};

export default function Edit({ review }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout header={t("admin.reviews.new.header")}>
      <Menu />
      <Form
        data={review}
        url={Routes.admin_reviews_path()}
      />
    </AdminLayout>
  );
}
