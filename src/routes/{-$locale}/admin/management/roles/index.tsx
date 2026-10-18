import { Group, Stack, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminDeleteRoleMutation,
  adminListRolesOptions,
  adminListRolesQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zAdminListRolesQuery } from "@/client/zod.gen";
import type { StaffRole } from "@/client/types.gen";
import { EditLink } from "@/components/admin/EditLink";
import { ButtonLink } from "@/components/RouterLink";
import { DataTable } from "mantine-datatable";
import type { DataTableColumn } from "mantine-datatable";
import { RowActions } from "@/components/admin/RowActions";
import { useQuery } from "@tanstack/react-query";
import type { ReactElement } from "react";

export const Route = createFileRoute("/{-$locale}/admin/management/roles/")({
  component: RolesList,
  validateSearch: zAdminListRolesQuery,
});

function RolesList(): ReactElement {
  const { t } = useTranslation();
  const search = Route.useSearch();
  const navigate = Route.useNavigate();
  const { data, isFetching } = useQuery(
    adminListRolesOptions({ query: { page: search.page ?? 1, perPage: search.perPage } }),
  );

  const columns: DataTableColumn<StaffRole>[] = [
    {
      accessor: "name",
      title: t(($) => $.models.attributes.base.name),
    },
    {
      accessor: "description",
      title: t(($) => $.models.attributes.base.description),
    },
    {
      accessor: "permissionsCount",
      title: t(($) => $.models.attributes.staff_role.permissions_count),
    },
    {
      accessor: "actions",
      render: (record) => (
        <RowActions
          describe={record.name}
          id={record.id}
          listKey={adminListRolesQueryKey()}
          remove={adminDeleteRoleMutation()}
        >
          <EditLink
            to="/{-$locale}/admin/management/roles/$id"
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
        <Title order={2}>{t(($) => $.admin.resources.roles)}</Title>
        <ButtonLink to="/{-$locale}/admin/management/roles/new">
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
