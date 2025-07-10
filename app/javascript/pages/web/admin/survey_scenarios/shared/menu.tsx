import { useTranslation } from 'react-i18next';
import {
  CrudHorizontalMenu,
  type CrudHorizontalMenuItem,
} from '@/components/CrudHorizontalMenu';
import * as Routes from '@/routes.js';
import type { SurveyScenarioCrud } from '@/types';

type Props = {
  data?: SurveyScenarioCrud;
};

export function Menu({ data }: Props) {
  const { t: tHelpers } = useTranslation('helpers');

  const items: CrudHorizontalMenuItem[] = [
    {
      href: Routes.admin_survey_scenarios_path(),
      label: tHelpers('crud.list'),
    },
    {
      href: Routes.new_admin_survey_scenario_path(),
      label: tHelpers('crud.add'),
    },
  ];

  if (data) {
    items.push({
      href: Routes.edit_admin_survey_scenario_path(data.survey_scenario.id),
      label: tHelpers('crud.editing'),
    });
  }

  return <CrudHorizontalMenu items={items} />;
}
