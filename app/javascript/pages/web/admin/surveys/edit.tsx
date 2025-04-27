import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import AdminLayout from "@/pages/layouts/AdminLayout";
import type { Language, SurveyCrud, } from "@/types/serializers";
import Form from "./shared/form";
import { Menu } from "./shared/menu";

type Props = {
  surveyDto: SurveyCrud;
  courses: Language[];
};

export default function Edit({ surveyDto, courses }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout
      header={t("admin.surveys.edit.header", { id: surveyDto.survey.id })}
    >
      <Menu data={surveyDto} />
      <Form
        method="patch"
        data={surveyDto}
        url={Routes.admin_survey_path(surveyDto.survey.id)}
      />
    </AdminLayout>
  );
}

