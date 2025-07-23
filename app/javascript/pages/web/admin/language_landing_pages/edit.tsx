import { useTranslation } from 'react-i18next';
import AdminLayout from '@/pages/layouts/AdminLayout';
import * as Routes from '@/routes.js';
import type {
  Language,
  LanguageLandingPage,
  LanguageLandingPageCrud,
} from '@/types';

import Form from './shared/form';
import { Menu } from './shared/menu';

type Props = {
  landingPageDto: LanguageLandingPageCrud;
  landingPages: LanguageLandingPage[];
  languages: Language[];
  // courseVersions: LanguageVersion[];
};

export default function Edit({
  landingPageDto,
  landingPages,
  languages,
}: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout
      header={t('admin.language_landing_pages.edit.header', {
        id: landingPageDto.data.header,
      })}
    >
      <Menu data={landingPageDto} />
      <Form
        languages={languages}
        landingPages={landingPages}
        method="patch"
        data={landingPageDto}
        url={Routes.admin_language_landing_page_path(landingPageDto.data.id)}
      />
    </AdminLayout>
  );
}
