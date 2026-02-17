import { Button, Modal } from "@mantine/core";
import dayjs from "dayjs";
import { DataTable } from "mantine-datatable";
import type { PropsWithChildren } from "react";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import AppAnchor from "@/components/Elements/AppAnchor";
import MarkdownViewer from "@/components/MarkdownViewer";
import useDataTableProps from "@/hooks/useDataTableProps";
import AdminLayout from "@/layouts/AdminLayout";
import * as Routes from "@/routes.js";
import type { Grid, LanguageLessonMemberMessage } from "@/types";

type Props = PropsWithChildren & {
  messages: LanguageLessonMemberMessage[];
  grid: Grid;
};

function renderLesson(item: LanguageLessonMemberMessage) {
  return (
    <AppAnchor
      href={Routes.language_lesson_path(
        item.language_slug,
        item.language_lesson_slug,
      )}
    >
      {item.language_lesson_name}
    </AppAnchor>
  );
}

type MessageBodyCellProps = {
  item: LanguageLessonMemberMessage;
};

function MessageBodyCell({ item }: MessageBodyCellProps) {
  const [visible, setVisible] = useState(false);
  return (
    <>
      <Button variant="subtle" onClick={() => setVisible(true)}>
        {item.content?.slice(0, 50)}
      </Button>
      <Modal opened={visible} onClose={() => setVisible(false)} size="50vw">
        <MarkdownViewer>{item.content}</MarkdownViewer>
      </Modal>
    </>
  );
}

export default function Index({ grid, messages }: Props) {
  const { t } = useTranslation();
  const { gridProps } = useDataTableProps<LanguageLessonMemberMessage, {}>(
    grid,
  );

  return (
    <AdminLayout header={t(($) => $.admin.messages.index.header)}>
      <DataTable
        records={messages}
        columns={[
          { accessor: "id" },
          { accessor: "role" },
          { accessor: "user_id" },
          { accessor: "lesson", title: "Lesson Url", render: renderLesson },
          {
            accessor: "body",
            title: "body",
            render: (item) => <MessageBodyCell item={item} />,
          },
          {
            accessor: "created_at",
            sortable: true,
            render: (r) => dayjs(r.created_at).format("LL"),
          },
        ]}
        {...gridProps}
      />
    </AdminLayout>
  );
}
