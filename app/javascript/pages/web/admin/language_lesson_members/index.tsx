import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import { DTDateTemplate } from "@/components/dtTemplates";
import useDataTable from "@/hooks/useDataTable";
import { fieldsToFilters } from "@/lib/utils";
import AdminLayout from "@/pages/layouts/AdminLayout";
import type { Grid, LanguageLessonMember } from "@/types/serializers";
import { Link } from "@inertiajs/react";
import { Column } from "primereact/column";
import { DataTable } from "primereact/datatable";
import { Dialog } from "primereact/dialog";
import { useState } from "react";
// import { Menu } from "./shared/menu";

type Props = {
  languageLessonMembers: LanguageLessonMember[];
  grid: Grid;
};

export function MessagesUrlTemplate(member: LanguageLessonMember) {
  const url = Routes.admin_messages_path({
    fields: { language_lesson_member_id_eq: member.id },
  });
  return <Link href={url}>{member.messages_count}</Link>;
}

export function LessonUrlTemplate(member: LanguageLessonMember) {
  const url = Routes.language_lesson_path(
    member.language_slug,
    member.language_lesson_slug,
  );
  return (
    <div>
      <Link href={url}>{member.language_lesson_name}</Link>
      <div className="text-muted small">{member.language_slug}</div>
    </div>
  );
}

export default function Index({ grid, languageLessonMembers }: Props) {
  const { t } = useTranslation();

  const handleDataTable = useDataTable();

  return (
    <AdminLayout header={t("admin.language_lesson_members.index.header")}>
      {/* <Menu /> */}
      <DataTable
        lazy
        first={grid.first}
        paginator
        totalRecords={grid.tr}
        rows={grid.per}
        sortField={grid.sf}
        sortOrder={grid.so}
        onSort={handleDataTable}
        onFilter={handleDataTable}
        onPage={handleDataTable}
        filters={fieldsToFilters(grid.fields)}
        value={languageLessonMembers}
      >
        <Column field="id" header="id" />
        <Column body={LessonUrlTemplate} header="Lesson Url" />
        <Column field="user_id" header="user_id" />
        <Column field="state" header="state" />
        <Column field="openai_thread_id" header="openai_thread_id" />
        <Column body={MessagesUrlTemplate} header="messages_count" />
        <Column
          field="created_at"
          header="created_at"
          sortable
          body={DTDateTemplate}
        />
        {/* <Column header="actions" body={actionBodyTemplate} /> */}
      </DataTable>
    </AdminLayout>
  );
}
