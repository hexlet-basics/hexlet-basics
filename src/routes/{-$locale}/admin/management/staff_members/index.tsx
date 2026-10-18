import { Group, Stack, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminDeleteStaffMemberMutation,
  adminListStaffMembersOptions,
  adminListStaffMembersQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zAdminListStaffMembersQuery } from "@/client/zod.gen";
import type { StaffMember } from "@/client/types.gen";
import { EditLink } from "@/components/admin/EditLink";
import { ButtonLink } from "@/components/RouterLink";
import { DataTable } from "mantine-datatable";
import type { DataTableColumn } from "mantine-datatable";
import { RowActions } from "@/components/admin/RowActions";
import { useQuery } from "@tanstack/react-query";
import type { ReactElement } from "react";

export const Route = createFileRoute("/{-$locale}/admin/management/staff_members/")({
  component: StaffMembersList,
  validateSearch: zAdminListStaffMembersQuery,
});

function StaffMembersList(): ReactElement {
  const { t } = useTranslation();
  const search = Route.useSearch();
  const navigate = Route.useNavigate();
  const { data, isFetching } = useQuery(
    adminListStaffMembersOptions({ query: { page: search.page ?? 1, perPage: search.perPage } }),
  );

  const columns: DataTableColumn<StaffMember>[] = [
    {
      accessor: "user",
      render: (record) => record.user.email ?? String(record.userId),
      title: t(($) => $.models.attributes.staff_member.user_id),
    },
    {
      accessor: "role",
      render: (record) => record.role.name,
      title: t(($) => $.models.attributes.staff_member.role_id),
    },
    {
      accessor: "allowedLocales",
      render: (record) => record.allowedLocales.join(", "),
      title: t(($) => $.models.attributes.staff_member.allowed_locales),
    },
    {
      accessor: "actions",
      render: (record) => (
        <RowActions
          describe={record.user.email ?? ""}
          id={record.id}
          listKey={adminListStaffMembersQueryKey()}
          remove={adminDeleteStaffMemberMutation()}
        >
          <EditLink
            to="/{-$locale}/admin/management/staff_members/$id"
            params={{ id: String(record.id) }}
          />
        </RowActions>
      ),
      textAlign: "right",
      title: "",
    },
  ];

  return (
    <Stack>
      <Group justify="space-between">
        <Title order={2}>{t(($) => $.admin.resources.staffMembers)}</Title>
        <ButtonLink to="/{-$locale}/admin/management/staff_members/new">
          {t(($) => $.admin.crud.new)}
        </ButtonLink>
      </Group>

      <DataTable
        withTableBorder
        striped
        highlightOnHover
        fetching={isFetching}
        records={data?.items}
        totalRecords={data?.total}
        recordsPerPage={data?.perPage ?? search.perPage}
        page={search.page ?? 1}
        onPageChange={(page) => {
          void navigate({ search: (prev) => ({ ...prev, page }) });
        }}
        noRecordsText={t(($) => $.admin.crud.empty)}
        columns={columns}
      />
    </Stack>
  );
}
