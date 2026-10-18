import { Stack, TextInput, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { adminListManagementUsersOptions } from "@/client/@tanstack/react-query.gen";
import { zAdminListManagementUsersQuery } from "@/client/zod.gen";
import type { UserCrud } from "@/client/types.gen";
import { EditLink } from "@/components/admin/EditLink";
import { DataTable } from "mantine-datatable";
import type { DataTableColumn } from "mantine-datatable";
import { RowActions } from "@/components/admin/RowActions";
import { SEARCH_DEBOUNCE_MS } from "@/lib/duration";
import { useState } from "react";
import { useDebouncedCallback } from "@mantine/hooks";
import { useQuery } from "@tanstack/react-query";
import type { ReactElement } from "react";

// Management users (legacy `/admin/management/users`) — list + edit only,
// no create/delete (parity with legacy: accounts are self-registered).
export const Route = createFileRoute("/{-$locale}/admin/management/users/")({
  component: ManagementUsersList,
  validateSearch: zAdminListManagementUsersQuery,
});

function ManagementUsersList(): ReactElement {
  const { t } = useTranslation();
  const search = Route.useSearch();
  const navigate = Route.useNavigate();
  const { data, isFetching } = useQuery(
    adminListManagementUsersOptions({
      query: { page: search.page ?? 1, perPage: search.perPage, q: search.q },
    }),
  );
  const [term, setTerm] = useState(search.q ?? "");
  const applyTerm = useDebouncedCallback((text: string) => {
    void navigate({
      search: (prev) => ({ ...prev, page: undefined, q: text === "" ? undefined : text }),
    });
  }, SEARCH_DEBOUNCE_MS);

  const columns: DataTableColumn<UserCrud>[] = [
    {
      accessor: "email",
      title: t(($) => $.models.attributes.user.email),
    },
    {
      accessor: "firstName",
      title: t(($) => $.models.attributes.user.first_name),
    },
    {
      accessor: "lastName",
      title: t(($) => $.models.attributes.user.last_name),
    },
    {
      accessor: "admin",
      render: (record) => (record.admin ? "✓" : ""),
      title: t(($) => $.models.attributes.user.admin),
    },
    {
      accessor: "actions",
      render: (record) => (
        <RowActions describe={record.email ?? ""} id={record.id}>
          <EditLink
            to="/{-$locale}/admin/management/users/$id"
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
      <Title order={2}>{t(($) => $.admin.resources.managementUsers)}</Title>

      <TextInput
        maw={360}
        placeholder={t(($) => $.admin.crud.search)}
        value={term}
        onChange={(event) => {
          setTerm(event.currentTarget.value);
          applyTerm(event.currentTarget.value.trim());
        }}
      />

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
