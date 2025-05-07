import { usePage } from "@inertiajs/react";
import { Nav } from "react-bootstrap";
import { useTranslation } from "react-i18next";

import * as Routes from "@/routes.js";
import { SurveyCrud, SurveyScenarioCrud } from "@/types";

type Props = {
  data?: SurveyScenarioCrud;
};

export function Menu({ data }: Props) {
  const { url } = usePage();
  const { t: tHelpers } = useTranslation("helpers");

  return (
    <Nav variant="tabs" className="mb-4" activeKey={url}>
      <Nav.Item>
        <Nav.Link
          className="link-body-emphasis"
          href={Routes.admin_survey_scenarios_path()}
        >
          {tHelpers("crud.list")}
        </Nav.Link>
      </Nav.Item>
      <Nav.Item>
        <Nav.Link
          className="link-body-emphasis"
          href={Routes.new_admin_survey_scenario_path()}
        >
          {tHelpers("crud.add")}
        </Nav.Link>
      </Nav.Item>
      {data && (
        <Nav.Item>
          <Nav.Link
            className="link-body-emphasis"
            href={Routes.edit_admin_survey_scenario_path(data.survey_scenario.id)}
          >
            {tHelpers("crud.editing")}
          </Nav.Link>
        </Nav.Item>
      )}
    </Nav>
  );
}

