import { useTranslation } from "react-i18next";
import AdminLayout from "@/layouts/AdminLayout";
import * as Routes from "@/routes.js";
import type {
  LanguageCategoryCreate,
  LanguageLandingPageForLists,
} from "@/types";
import Form from "./shared/form";
import { Menu } from "./shared/menu";

type Props = {
  categoryDto: LanguageCategoryCreate;
  landingPagesForCategories: LanguageLandingPageForLists[];
};

export default function New({ categoryDto, landingPagesForCategories }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout header={t(($) => $.admin.language_categories.new.header)}>
      <Menu />
      <Form
        data={categoryDto}
        landingPagesForCategories={landingPagesForCategories}
        url={Routes.admin_language_categories_path()}
      />
    </AdminLayout>
  );
}
