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
    <AdminLayout header={t("admin.reviews.edit.header", { id: review.id })}>
      <Menu data={review} />
      <Form
        method="patch"
        data={review}
        url={Routes.edit_admin_review_path(review.id)}
      />
    </AdminLayout>
  );
}
