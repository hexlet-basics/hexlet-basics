import { useTranslation } from 'react-i18next';
import AdminLayout from '@/pages/layouts/AdminLayout';
import * as Routes from '@/routes.js';
import type { SurveyCrud, SurveyCrudWithAttrs } from '@/types';
import Form from './shared/form';
import { Menu } from './shared/menu';

type Props = {
  surveyDto: SurveyCrudWithAttrs;
};

export default function New({ surveyDto }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout header={t('admin.surveys.new.header')}>
      <Menu />
      <Form data={surveyDto} url={Routes.admin_surveys_path()} />
    </AdminLayout>
  );
}
