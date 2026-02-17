import { useTranslation } from "react-i18next";
import AdminLayout from "@/layouts/AdminLayout";
import * as Routes from "@/routes.js";
import type { Survey, SurveyItemCrud, SurveyScenarioCrud } from "@/types";
import Form from "./shared/form";
import { Menu } from "./shared/menu";

type Props = {
  surveyScenarioDto: SurveyScenarioCrud;
  surveysItems: SurveyItemCrud[];
  surveys: Survey[];
};

export default function New({
  surveys,
  surveyScenarioDto,
  surveysItems,
}: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout header={t(($) => $.admin.survey_scenarios.new.header)}>
      <Menu />
      <Form
        surveys={surveys}
        surveysItems={surveysItems}
        data={surveyScenarioDto}
        url={Routes.admin_survey_scenarios_path()}
      />
    </AdminLayout>
  );
}
