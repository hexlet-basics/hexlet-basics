import { Button, Checkbox, Grid, Stack, TextInput, Title } from "@mantine/core";
import { DataTable } from "mantine-datatable";
import { useTranslation } from "react-i18next";
import { useAppForm } from "@/hooks/useAppForm";
import AdminLayout from "@/layouts/AdminLayout";
import * as Routes from "@/routes.js";
import type { UserCrud } from "@/types";
import { Menu } from "./shared/menu";

type Props = {
  user: UserCrud;
  progress: Record<string, number>[];
};

export default function Edit({ user, progress }: Props) {
  const { t } = useTranslation();

  const payload = user;

  const { onSubmit, processing, form } = useAppForm(payload, {
    url: Routes.admin_management_user_path(user.id),
    method: "patch",
  });

  return (
    <AdminLayout header={t(($) => $.admin.management.users.edit.header)}>
      <Menu data={user} />
      <Grid>
        <Grid.Col span={5}>
          <form onSubmit={onSubmit}>
            <TextInput
              {...form.getInputProps("email")}
              autoComplete="email"
              disabled
            />
            <Checkbox {...form.getCheckboxProps("admin")} />
            <TextInput
              {...form.getInputProps("first_name")}
              autoComplete="name"
            />
            <TextInput
              {...form.getInputProps("last_name")}
              autoComplete="name"
            />
            <Button type="submit" loading={processing}>
              {t(($) => $.helpers.submit.save)}
            </Button>
          </form>
        </Grid.Col>
        <Grid.Col>
          <Stack>
            <Title order={2}>
              {t(($) => $.admin.management.users.edit.progress)}
            </Title>
            <DataTable
              records={progress}
              columns={[
                {
                  accessor: "language",
                  title: t(($) => $.admin.management.users.edit.language),
                },
                {
                  accessor: "count",
                  title: t(($) => $.admin.management.users.edit.count),
                },
              ]}
            />
          </Stack>
        </Grid.Col>
      </Grid>
    </AdminLayout>
  );
}
