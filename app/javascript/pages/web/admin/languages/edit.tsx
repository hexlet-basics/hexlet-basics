import { DTDateTemplate } from "@/components/dtTemplates";
import AdminLayout from "@/pages/layouts/AdminLayout";
import * as Routes from "@/routes.js";
import type {
  LanguageCrud,
  LanguageLandingPage,
  LanguageVersion,
} from "@/types/serializers";
import { Link } from "@inertiajs/react";
import { DataTable } from 'mantine-datatable';
import { Grid, Button, Group, Title } from '@mantine/core';
import { useTranslation } from "react-i18next";

import Form from "./shared/form";
import { Menu } from "./shared/menu";

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
      header={t("admin.languages.edit.header", { id: courseDto.meta.slug })}
    >
      <Menu data={courseDto} landingPage={landingPage} />
      <Grid>
        <Grid.Col span={4}>
          <Form
            method="patch"
            data={courseDto}
            url={Routes.admin_language_path(courseDto.language.id)}
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
              href={Routes.admin_language_versions_path(courseDto.language.id)}
            >
              Load New Version
            </Button>
          </Group>
          <DataTable
            records={courseVersions}
            columns={[
              { accessor: 'id', title: 'id' },
              { accessor: 'result', title: 'result' },
              { accessor: 'created_at', title: 'created_at', render: DTDateTemplate },
            ]}
          />
        </Grid.Col>
      </Grid>
    </AdminLayout>
  );
}
