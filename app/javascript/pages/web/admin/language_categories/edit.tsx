import { useTranslation } from "react-i18next";
import AdminLayout from "@/layouts/AdminLayout";
import * as Routes from "@/routes.js";
import type {
  LanguageCategoryQnaItem,
  LanguageCategoryUpdate,
  LanguageLandingPageForLists,
} from "@/types";
import Form from "./shared/form";
import { Menu } from "./shared/menu";

type Props = {
  categoryDto: LanguageCategoryUpdate;
  qnaItems: LanguageCategoryQnaItem[];
  landingPagesForCategories: LanguageLandingPageForLists[];
};

export default function Edit({
  categoryDto,
  qnaItems,
  landingPagesForCategories,
}: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout
      header={t(($) => $.admin.language_categories.edit.header, {
        id: categoryDto.name,
      })}
    >
      <Menu data={categoryDto} />
      <Form
        method="patch"
        data={categoryDto}
        qnaItems={qnaItems}
        landingPagesForCategories={landingPagesForCategories}
        url={Routes.admin_language_category_path(categoryDto.id)}
      />
    </AdminLayout>
  );
}
