import { ActionIcon, Button, Code, Modal } from '@mantine/core';
import dayjs from 'dayjs';
import { DataTable } from 'mantine-datatable';
import type { PropsWithChildren } from 'react';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import MarkdownViewer from '@/components/MarkdownViewer';
import useDataTableProps from '@/hooks/useDataTableProps';
import AdminLayout from '@/pages/layouts/AdminLayout';
import * as Routes from '@/routes.js';
import type { Grid, LanguageLesson, LanguageOriginalLesson } from '@/types';
import AppAnchor from '@/components/Elements/AppAnchor';
import { Link } from 'lucide-react';

type Props = PropsWithChildren & {
  lessons: LanguageOriginalLesson[];
  grid: Grid;
};

function DataBox({
  lesson,
  col,
}: {
  lesson: LanguageOriginalLesson;
  col: string;
}) {
  const [visible, setVisible] = useState(false);
  return (
    <>
      <Button variant="subtle" onClick={() => setVisible(true)}>
        data
      </Button>
      <Modal opened={visible} onClose={() => setVisible(false)} size="50vw">
        <MarkdownViewer>{lesson.review || ''}</MarkdownViewer>
      </Modal>
    </>
  );
}

export default function Index({ grid, lessons }: Props) {
  const { t } = useTranslation();
  const { gridProps } = useDataTableProps<LanguageOriginalLesson, {}>(grid);

  const renderActions = (item: LanguageOriginalLesson) => {
    return (
      <>
        {/* <AppAnchor */}
        {/*   me="xs" */}
        {/*   href={Routes.edit_admin_language_category_path(item.id)} */}
        {/* > */}
        {/*   <ActionIcon variant="default" size="xs"> */}
        {/*     <Edit /> */}
        {/*   </ActionIcon> */}
        {/* </AppAnchor> */}
        <AppAnchor external href={Routes.language_lesson_path(item.language!.slug!, item.slug!)}>
          <ActionIcon variant="default" size="xs">
            <Link />
          </ActionIcon>
        </AppAnchor>
        <AppAnchor external href={item.source_code_url}>
          <ActionIcon variant="default" size="xs">
            <Github />
          </ActionIcon>
        </AppAnchor>
        {/* <Link */}
        {/*   onClick={confirmDeleting} */}
        {/*   className="btn btn-link link-body-emphasis p-0 m-0" */}
        {/*   method="delete" */}
        {/*   href={Routes.admin_language_category_path(data.id)} */}
        {/* > */}
        {/*   {<i className="bi bi-file-x" />} */}
        {/* </Link> */}
      </>
    );
  };

  return (
    <AdminLayout header={t('admin.language_lessons.index.header')}>
      <DataTable
        records={lessons}
        columns={[
          { accessor: 'id' },
          { accessor: 'language.slug' },
          { accessor: 'slug' },
          {
            accessor: 'review',
            title: 'review',
            render: (rec) => <DataBox lesson={rec} col="review" />,
          },
          {
            accessor: 'created_at',
            render: (r) => dayjs(r.created_at).format('LL'),
          },
          {
            accessor: 'actions',
            title: 'actions',
            render: renderActions,
          },
        ]}
        {...gridProps}
      />
    </AdminLayout>
  );
}
