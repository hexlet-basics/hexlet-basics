import { Link } from '@inertiajs/react';
import { Box, Button, Grid, Group, SimpleGrid, Title } from '@mantine/core';
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
  const courseMeta = (courseDto as { meta?: { slug?: string } }).meta ?? {};

  return (
    <AdminLayout
      header={t(($) => $.admin.languages.edit.header, {
        id: courseMeta.slug ?? courseDto.data.id,
      })}
    >
      <Menu data={courseDto} landingPage={landingPage} />
      <SimpleGrid cols={{ base: 1, sm: 2 }}>
        <Box>
          <Form
            method="patch"
            data={courseDto}
            url={Routes.admin_language_path(courseDto.data.id)}
          />
        </Box>
        <Box>
          <Group justify="space-between" mb="xs">
            <Title order={2}>{t(($) => $.admin.languages.edit.versions)}</Title>
            <Button
              component={Link}
              method="post"
              variant="outline"
              size="sm"
              href={Routes.admin_language_versions_path(courseDto.data.id)}
            >
              {t(($) => $.admin.languages.edit.load_new_version)}
            </Button>
          </Group>
          <DataTable
            records={courseVersions}
            columns={[
              { accessor: 'id', title: t(($) => $.admin.languages.edit.id) },
              {
                accessor: 'result',
                title: t(($) => $.admin.languages.edit.result),
              },
              {
                accessor: 'created_at',
                sortable: true,
                render: (r) => dayjs(r.created_at).format('LL'),
              },
            ]}
          />
        </Box>
      </SimpleGrid>
    </AdminLayout>
  );
}
