import { useTranslation } from 'react-i18next';
import AdminLayout from '@/pages/layouts/AdminLayout';
import * as Routes from '@/routes.js';
import type { LanguageCategoryCrudWithAttrs } from '@/types';
import Form from './shared/form';
import { Menu } from './shared/menu';

type Props = {
  categoryDto: LanguageCategoryCrudWithAttrs;
  // courseVersions: LanguageVersion[];
};

export default function Edit({ categoryDto }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout
      header={t('admin.language_categories.edit.header', {
        id: categoryDto.data.name,
      })}
    >
      <Menu data={categoryDto} />
      <Form
        method="patch"
        data={categoryDto}
        url={Routes.admin_language_category_path(categoryDto.data.id)}
      />
    </AdminLayout>
  );
}
