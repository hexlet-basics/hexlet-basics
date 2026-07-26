import { useTranslation } from "react-i18next";
import AdminLayout from "@/layouts/AdminLayout";
import * as Routes from "@/routes.js";
import type { Language, ReviewUpdate } from "@/types/serializers";
import Form from "./shared/form";
import { Menu } from "./shared/menu";

type Props = {
  reviewDto: ReviewUpdate;
  courses: Language[];
};

export default function Edit({ reviewDto, courses }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout header={t(($) => $.admin.reviews.edit.header)}>
      <Menu data={reviewDto} />
      <Form
        courses={courses}
        method="patch"
        data={reviewDto}
        url={Routes.admin_review_path(reviewDto.id)}
      />
    </AdminLayout>
  );
}
