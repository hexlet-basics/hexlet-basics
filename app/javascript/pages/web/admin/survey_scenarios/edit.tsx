import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import AdminLayout from "@/pages/layouts/AdminLayout";
import type { SurveyCrud, SurveyScenarioCrud, SurveyItemCrud, Survey } from "@/types";
import Form from "./shared/form";
import { Menu } from "./shared/menu";

type Props = {
  surveyScenarioDto: SurveyScenarioCrud;
  surveys: Survey[];
  surveysItems: SurveyItemCrud[]
};

export default function Edit({ surveys, surveyScenarioDto, surveysItems }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout
      header={t("admin.survey_scenarios.edit.header")}
    >
      <Menu data={surveyScenarioDto} />
      <Form
        surveys={surveys}
        surveysItems={surveysItems}
        method="patch"
        data={surveyScenarioDto}
        url={Routes.admin_survey_scenario_path(surveyScenarioDto.survey_scenario.id)}
      />
    </AdminLayout>
  );
}

