import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import AdminLayout from "@/pages/layouts/AdminLayout";
import type { Language, SurveyCrud, SurveyItemCrud, } from "@/types";
import Form from "./shared/form";
import { Menu } from "./shared/menu";

type Props = {
  surveyDto: SurveyCrud;
  surveyItems: SurveyItemCrud[];
  courses: Language[];
};

export default function New({ surveyDto, surveyItems }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout header={t("admin.surveys.new.header")}>
      <Menu />
      <Form
        surveyItems={surveyItems}
        data={surveyDto}
        url={Routes.admin_surveys_path()}
      />
    </AdminLayout>
  );
}

