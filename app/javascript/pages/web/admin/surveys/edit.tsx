import { useTranslation } from 'react-i18next';
import AdminLayout from '@/pages/layouts/AdminLayout';
import * as Routes from '@/routes.js';
import type { Language, SurveyCrud, SurveyItemCrud } from '@/types/serializers';
import Form from './shared/form';
import { Menu } from './shared/menu';

type Props = {
  surveyDto: SurveyCrud;
  courses: Language[];
  surveyItems: SurveyItemCrud[];
};

export default function Edit({ surveyDto, surveyItems }: Props) {
  const { t } = useTranslation();
  // console.log(surveyItems)

  return (
    <AdminLayout
      header={t('admin.surveys.edit.header', { id: surveyDto.survey.id })}
    >
      <Menu data={surveyDto} />
      <Form
        surveyItems={surveyItems}
        method="patch"
        data={surveyDto}
        url={Routes.admin_survey_path(surveyDto.survey.id)}
      />
    </AdminLayout>
  );
}
