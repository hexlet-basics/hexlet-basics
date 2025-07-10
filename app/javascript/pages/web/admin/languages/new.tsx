import { Grid } from '@mantine/core';
import { useTranslation } from 'react-i18next';

import AdminLayout from '@/pages/layouts/AdminLayout';
import * as Routes from '@/routes.js';
import type LanguageCrud from '@/types/serializers/LanguageCrud';
import Form from './shared/form';
import { Menu } from './shared/menu';

type Props = {
  courseDto: LanguageCrud;
};

export default function New({ courseDto }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout header={t('admin.languages.new.header')}>
      <Menu />
      <Grid>
        <Grid.Col span={7}>
          <Form data={courseDto} url={Routes.admin_languages_path()} />
        </Grid.Col>
      </Grid>
    </AdminLayout>
  );
}
