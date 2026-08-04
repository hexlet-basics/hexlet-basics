import { ActionIcon, Group, Stack, TextInput, Title, Tooltip } from "@mantine/core";
import { useDebouncedValue } from "@mantine/hooks";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import type { PaginationState, SortingState } from "@tanstack/react-table";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import {
  adminDeleteUserMutation,
  adminListUsersOptions,
  adminListUsersQueryKey,
  adminSearchUsersOptions,
} from "@/client/@tanstack/react-query.gen";
import type { UserCrud } from "@/client/types.gen";
import { type CrudColumnDef, CrudList } from "@/components/admin/CrudList";
import { ButtonLink } from "@/components/RouterLink";
import { useDeleteConfirmation } from "@/hooks/useDeleteConfirmation";
import { useResourceMutation } from "@/hooks/useResourceMutation";

// Users admin (legacy `/admin/api/users` — the `api/` segment was a Rails
// namespacing artifact and is dropped from the UI path). A non-empty search
// swaps the paginated list for the typeahead endpoint (legacy `#search`).
export const Route = createFileRoute("/{-$locale}/admin/users/")({
  component: UsersList,
});

function UsersList() {
  const { t } = useTranslation();
  const confirmDelete = useDeleteConfirmation();

  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize: 25,
  });
  const [sorting, setSorting] = useState<SortingState>([]);
  const [search, setSearch] = useState("");
  const [query] = useDebouncedValue(search.trim(), 300);

  const list = useQuery({
    ...adminListUsersOptions({
      query: { page: pagination.pageIndex + 1, perPage: pagination.pageSize },
    }),
    enabled: query === "",
  });
  const found = useQuery({
    ...adminSearchUsersOptions({ query: { q: query } }),
    enabled: query !== "",
  });

  const searching = query !== "";
  const items = searching ? (found.data ?? []) : (list.data?.items ?? []);
  const total = searching ? (found.data?.length ?? 0) : (list.data?.total ?? 0);

  const deleteMutation = useResourceMutation({
    mutation: adminDeleteUserMutation(),
    invalidate: [adminListUsersQueryKey()],
    successMessage: t(($) => $.admin.crud.deleted),
    errorMessage: t(($) => $.admin.crud.deleteError),
  });

  const columns: CrudColumnDef<UserCrud>[] = [
    {
      accessorKey: "email",
      header: t(($) => $.models.attributes.user.email),
      enableSorting: false,
    },
    {
      accessorKey: "firstName",
      header: t(($) => $.models.attributes.user.first_name),
      enableSorting: false,
    },
    {
      accessorKey: "lastName",
      header: t(($) => $.models.attributes.user.last_name),
      enableSorting: false,
    },
    {
      id: "admin",
      header: t(($) => $.models.attributes.user.admin),
      cell: ({ row }) => (row.original.admin ? "✓" : ""),
      enableSorting: false,
    },
    {
      id: "actions",
      header: "",
      cell: ({ row }) => (
        <Group gap="xs" justify="flex-end" wrap="nowrap">
          <ButtonLink
            to="/{-$locale}/admin/users/$id"
            params={{ id: String(row.original.id) }}
            size="xs"
            variant="light"
          >
            {t(($) => $.admin.crud.edit)}
          </ButtonLink>
          <Tooltip label={t(($) => $.admin.crud.delete)}>
            <ActionIcon
              color="red"
              variant="light"
              aria-label={t(($) => $.admin.crud.delete)}
              onClick={() =>
                confirmDelete({
                  description: row.original.email ?? "",
                  onConfirm: () => deleteMutation.mutate({ path: { id: row.original.id } }),
                })
              }
            >
              ✕
            </ActionIcon>
          </Tooltip>
        </Group>
      ),
    },
  ];

  return (
    <Stack>
      <Group justify="space-between">
        <Title order={2}>{t(($) => $.admin.resources.users)}</Title>
        <ButtonLink to="/{-$locale}/admin/users/new">{t(($) => $.admin.crud.new)}</ButtonLink>
      </Group>

      <TextInput
        placeholder={t(($) => $.admin.crud.search)}
        value={search}
        onChange={(event) => setSearch(event.currentTarget.value)}
        maw={360}
      />

      <CrudList
        columns={columns}
        data={items}
        total={total}
        pagination={pagination}
        onPaginationChange={setPagination}
        sorting={sorting}
        onSortingChange={setSorting}
        isLoading={searching ? found.isLoading : list.isLoading}
      />
    </Stack>
  );
}
