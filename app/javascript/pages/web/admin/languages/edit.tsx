import { Link } from '@inertiajs/react';
import { Button, Grid, Group, Title } from '@mantine/core';
import dayjs from 'dayjs';
import { DataTable } from 'mantine-datatable';
import { useTranslation } from 'react-i18next';
import AdminLayout from '@/pages/layouts/AdminLayout';
import * as Routes from '@/routes.js';
import type {
  LanguageCrud,
  LanguageLandingPage,
  LanguageVersion,
} from '@/types/serializers';
import Form from './shared/form';
import { Menu } from './shared/menu';

type Props = {
  courseDto: LanguageCrud;
  landingPage?: LanguageLandingPage;
  courseVersions: LanguageVersion[];
};

export default function Edit({
  courseDto,
  courseVersions,
  landingPage,
}: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout
      header={t('admin.languages.edit.header', { id: courseDto.meta.slug })}
    >
      <Menu data={courseDto} landingPage={landingPage} />
      <Grid>
        <Grid.Col span={4}>
          <Form
            method="patch"
            data={courseDto}
            url={Routes.admin_language_path(courseDto.data.id)}
          />
        </Grid.Col>
        <Grid.Col>
          <Group justify="space-between" mb="xs">
            <Title order={2}>Versions</Title>
            <Button
              component={Link}
              method="post"
              variant="outline"
              size="sm"
              href={Routes.admin_language_versions_path(courseDto.data.id)}
            >
              Load New Version
            </Button>
          </Group>
          <DataTable
            records={courseVersions}
            columns={[
              { accessor: 'id', title: 'id' },
              { accessor: 'result', title: 'result' },
              {
                accessor: 'created_at',
                sortable: true,
                render: (r) => dayjs(r.created_at).format('LL'),
              },
            ]}
          />
        </Grid.Col>
      </Grid>
    </AdminLayout>
  );
}
