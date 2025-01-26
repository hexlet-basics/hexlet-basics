import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import AdminLayout from "@/pages/layouts/AdminLayout";
import type { OriginalLanguage, Review } from "@/types/serializers";
import Form from "./shared/form";
import { Menu } from "./shared/menu";

type Props = {
  review: Review;
  courses: OriginalLanguage[];
};

export default function Edit({ review, courses }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout header={t("admin.reviews.edit.header", { id: review.id })}>
      <Menu data={review} />
      <Form
        courses={courses}
        method="patch"
        data={review}
        url={Routes.admin_review_path(review.id)}
      />
    </AdminLayout>
  );
}
