import { DataTable } from 'mantine-datatable';
import type { PropsWithChildren } from 'react';

import * as Routes from '@/routes.js';
import { useTranslation } from 'react-i18next';

import AdminLayout from '@/pages/layouts/AdminLayout';
import AppAnchor from '@/components/AppAnchor';
import type { LanguageLessonMemberMessage, Grid } from '@/types/serializers';
import useDataTableProps from '@/hooks/useDataTableProps';
import { Modal, Button } from '@mantine/core';
import { useState } from 'react';
import MarkdownViewer from '@/components/MarkdownViewer';
import dayjs from 'dayjs';

type Props = PropsWithChildren & {
  messages: LanguageLessonMemberMessage[];
  grid: Grid;
};

function renderLesson(item: LanguageLessonMemberMessage) {
  return (
    <AppAnchor href={Routes.language_lesson_path(item.language_slug, item.language_lesson_slug)}>
      {item.language_lesson_name}
    </AppAnchor>
  );
}

function renderBody(item: LanguageLessonMemberMessage) {
  const [visible, setVisible] = useState(false);
  return (
    <>
      <Button variant="subtle" onClick={() => setVisible(true)}>
        {item.content?.slice(0, 50)}
      </Button>
      <Modal
        opened={visible}
        onClose={() => setVisible(false)}
        size="50vw"
      >
        <MarkdownViewer>{item.content}</MarkdownViewer>
      </Modal>
    </>
  );
}

export default function Index({ grid, messages }: Props) {
  const { t } = useTranslation();
  const { gridProps } = useDataTableProps<LanguageLessonMemberMessage, {}>(grid);

  return (
    <AdminLayout header={t('admin.messages.index.header')}>
      <DataTable
        records={messages}
        columns={[
          { accessor: 'id' },
          { accessor: 'role' },
          { accessor: 'user_id' },
          { accessor: 'lesson', title: 'Lesson Url', render: renderLesson },
          { accessor: 'body', title: 'body', render: renderBody },
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
