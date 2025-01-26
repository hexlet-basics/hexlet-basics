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

export default function New({ review, courses }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout header={t("admin.reviews.new.header")}>
      <Menu />
      <Form courses={courses} data={review} url={Routes.admin_reviews_path()} />
    </AdminLayout>
  );
}
