import { Stack, Text } from '@mantine/core';
import dayjs from 'dayjs';
import { DataTable } from 'mantine-datatable';
import type { PropsWithChildren } from 'react';
import { useTranslation } from 'react-i18next';
import AppAnchor from '@/components/Elements/AppAnchor';
import useDataTableProps from '@/hooks/useDataTableProps';
import AdminLayout from '@/pages/layouts/AdminLayout';
import * as Routes from '@/routes.js';
import type { Grid, LanguageLessonMember } from '@/types/serializers';

type Props = {
  languageLessonMembers: LanguageLessonMember[];
  grid: Grid;
};

function renderMessages(member: LanguageLessonMember) {
  return (
    <AppAnchor
      href={Routes.admin_messages_path({
        fields: { language_lesson_member_id_eq: member.id },
      })}
    >
      {member.messages_count}
    </AppAnchor>
  );
}

function renderLesson(member: LanguageLessonMember) {
  return (
    <Stack>
      <AppAnchor
        href={Routes.language_lesson_path(
          member.language_slug,
          member.language_lesson_slug,
        )}
      >
        {member.language_lesson_name}
      </AppAnchor>
      <Text c="dimmed" size="sm">
        {member.language_slug}
      </Text>
    </Stack>
  );
}

export default function Index({
  grid,
  languageLessonMembers,
}: PropsWithChildren<Props>) {
  const { t } = useTranslation();
  const { gridProps } = useDataTableProps<LanguageLessonMember, {}>(grid);

  return (
    <AdminLayout
      header={t(($) => $.admin.language_lesson_members.index.header)}
    >
      <DataTable
        records={languageLessonMembers}
        columns={[
          { accessor: 'id' },
          { accessor: 'lesson', title: 'Lesson Url', render: renderLesson },
          { accessor: 'user_id' },
          { accessor: 'state' },
          { accessor: 'openai_thread_id' },
          {
            accessor: 'messages_count',
            title: 'messages_count',
            render: renderMessages,
          },
          {
            accessor: 'created_at',
            sortable: true,
            render: (r) => dayjs(r.created_at).format('LL'),
          },
        ]}
        {...gridProps}
      />
    </AdminLayout>
  );
}
