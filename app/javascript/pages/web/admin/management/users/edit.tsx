import { Button, Checkbox, Grid, Stack, TextInput, Title } from '@mantine/core';
import { DataTable } from 'mantine-datatable';
import { useTranslation } from 'react-i18next';
import { useAppForm } from '@/hooks/useAppForm';
import AdminLayout from '@/pages/layouts/AdminLayout';
import * as Routes from '@/routes.js';
import type { UserCrud } from '@/types';
import { Menu } from './shared/menu';

type Props = {
  user: UserCrud;
  progress: Record<string, number>[];
};

export default function Edit({ user, progress }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation('helpers');

  const {
    getInputProps,
    submit,
    formState: { isSubmitting },
  } = useAppForm<UserCrud>({
    url: Routes.admin_management_user_path(user.data.id),
    method: 'patch',
    container: user,
  });

  return (
    <AdminLayout header={t(($) => $.admin.management.users.edit.header)}>
      <Menu data={user} />
      <Grid>
        <Grid.Col span={5}>
          <form onSubmit={submit}>
            <TextInput
              {...getInputProps('email')}
              autoComplete="email"
              disabled
            />
            <Checkbox {...getInputProps('admin')} />
            <TextInput {...getInputProps('first_name')} autoComplete="name" />
            <TextInput {...getInputProps('last_name')} autoComplete="name" />
            <Button type="submit" loading={isSubmitting}>
              {tHelpers(($) => $.submit.save)}
            </Button>
          </form>
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
