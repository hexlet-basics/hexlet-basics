import { useTranslation } from 'react-i18next';
import {
  CrudHorizontalMenu,
  type CrudHorizontalMenuItem,
} from '@/components/CrudHorizontalMenu';
import * as Routes from '@/routes.js';
import type { SurveyCrud } from '@/types';

type Props = {
  data?: SurveyCrud;
};

export function Menu({ data }: Props) {
  const { t: tHelpers } = useTranslation('helpers');

  const items: CrudHorizontalMenuItem[] = [
    { href: Routes.admin_surveys_path(), label: tHelpers('crud.list') },
    { href: Routes.new_admin_survey_path(), label: tHelpers('crud.add') },
  ];

  if (data) {
    items.push({
      href: Routes.edit_admin_survey_path(data.survey.id),
      label: tHelpers('crud.editing'),
    });
  }

  return <CrudHorizontalMenu items={items} />;
}
