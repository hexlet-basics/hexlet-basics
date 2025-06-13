import type { PropsWithChildren } from "react";

import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import { XCheck, XForm, XInput } from "@/components/forms";
import AdminLayout from "@/pages/layouts/AdminLayout";
import type { User } from "@/types/serializers";
import { DataTable } from 'mantine-datatable';
import { Grid, Button, Title, Stack } from '@mantine/core';
import { Submit } from "use-inertia-form";
import { Menu } from "./shared/menu";

type Props = PropsWithChildren & {
  user: User;
  progress: Record<string, number>[];
};

export default function Edit({ user, progress }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");

  return (
    <AdminLayout header={t("admin.management.users.edit.header")}>
      <Menu data={user} />
      <Grid>
        <Grid.Col span={5}>
          <XForm
            method="patch"
            model="user"
            data={{ user }}
            to={Routes.admin_management_user_path(user.id)}
          >
            <XInput field="email" autoComplete="email" disabled />
            <XCheck field="admin" />
            <XInput field="first_name" autoComplete="name" />
            <XInput field="last_name" autoComplete="name" />
            <XInput field="last_name" autoComplete="name" />
            <Button type="submit">
              {tHelpers("submit.save")}
            </Button>
          </XForm>
        </Grid.Col>
        <Grid.Col>
          <Stack>
            <Title order={2}>Progress</Title>
            <DataTable
              records={progress}
              columns={[
                { accessor: 'language', title: 'Language' },
                { accessor: 'count', title: 'Count' },
              ]}
            />
          </Stack>
        </Grid.Col>
      </Grid>
    </AdminLayout>
  );
}
