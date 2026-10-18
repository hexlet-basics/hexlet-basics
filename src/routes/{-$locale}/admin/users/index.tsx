import { Group, Stack, TextInput, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminDeleteUserMutation,
  adminListUsersOptions,
  adminListUsersQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zAdminListUsersQuery } from "@/client/zod.gen";
import type { UserCrud } from "@/client/types.gen";
import { EditLink } from "@/components/admin/EditLink";
import { ButtonLink } from "@/components/RouterLink";
import { DataTable } from "mantine-datatable";
import type { DataTableColumn } from "mantine-datatable";
import { RowActions } from "@/components/admin/RowActions";
import { SEARCH_DEBOUNCE_MS } from "@/lib/duration";
import { useState } from "react";
import { useDebouncedCallback } from "@mantine/hooks";
import { useQuery } from "@tanstack/react-query";
import type { ReactElement } from "react";

// Users admin (legacy `/admin/api/users` — the `api/` segment was a Rails
// namespacing artifact and is dropped from the UI path). The search box is the
// list's own `q` filter, so a search pages and totals like the list does.
export const Route = createFileRoute("/{-$locale}/admin/users/")({
  component: UsersList,
  validateSearch: zAdminListUsersQuery,
});

function UsersList(): ReactElement {
  const { t } = useTranslation();
  const search = Route.useSearch();
  const navigate = Route.useNavigate();
  const { data, isFetching } = useQuery(
    adminListUsersOptions({
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
    { accessor: "email", title: t(($) => $.models.attributes.user.email) },
    { accessor: "firstName", title: t(($) => $.models.attributes.user.first_name) },
    { accessor: "lastName", title: t(($) => $.models.attributes.user.last_name) },
    {
      accessor: "admin",
      render: (record) => (record.admin ? "✓" : ""),
      title: t(($) => $.models.attributes.user.admin),
    },
    {
      accessor: "actions",
      render: (record) => (
        <RowActions
          describe={record.email ?? ""}
          id={record.id}
          listKey={adminListUsersQueryKey()}
          remove={adminDeleteUserMutation()}
        >
          <EditLink to="/{-$locale}/admin/users/$id" params={{ id: String(record.id) }} />
        </RowActions>
      ),
      textAlign: "right",
      title: "",
    },
  ];

  return (
    <Stack>
      <Group justify="space-between">
        <Title order={2}>{t(($) => $.admin.resources.users)}</Title>
        <ButtonLink to="/{-$locale}/admin/users/new">{t(($) => $.admin.crud.new)}</ButtonLink>
      </Group>

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
