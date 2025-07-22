import { useTranslation } from 'react-i18next';
import AdminLayout from '@/pages/layouts/AdminLayout';
import * as Routes from '@/routes.js';
import type { LanguageCategoryCrudWithAttrs } from '@/types';
import Form from './shared/form';
import { Menu } from './shared/menu';

type Props = {
  categoryDto: LanguageCategoryCrudWithAttrs;
};

export default function New({ categoryDto }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout header={t('admin.language_categories.new.header')}>
      <Menu />
      <Form data={categoryDto} url={Routes.admin_language_categories_path()} />
    </AdminLayout>
  );
}
