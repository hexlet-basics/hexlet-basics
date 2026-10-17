import { Group, Stack, TextInput, Title } from "@mantine/core";
import { useDebouncedValue } from "@mantine/hooks";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
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
import { EditLink } from "@/components/admin/EditLink";
import { userColumns } from "@/components/admin/generated/user";
import { useGeneratedColumns } from "@/components/admin/useGeneratedFields";
import { useCrudList } from "@/components/admin/useCrudList";
import { ButtonLink } from "@/components/RouterLink";

// Users admin (legacy `/admin/api/users` — the `api/` segment was a Rails
// namespacing artifact and is dropped from the UI path). A non-empty search
// swaps the paginated list for the typeahead endpoint (legacy `#search`).
export const Route = createFileRoute("/{-$locale}/admin/users/")({
  component: UsersList,
});

function UsersList() {
  const { t } = useTranslation();

  const list = useCrudList({
    list: adminListUsersOptions,
    listKey: adminListUsersQueryKey,
    remove: adminDeleteUserMutation,
  });
  const [search, setSearch] = useState("");
  const [term] = useDebouncedValue(search.trim(), 300);

  const searching = term !== "";
  const found = useQuery({
    ...adminSearchUsersOptions({ query: { q: term } }),
    enabled: searching,
  });

  const columns: CrudColumnDef<UserCrud>[] = [
    ...useGeneratedColumns<UserCrud>(userColumns, ["email", "firstName", "lastName"]),
    {
      id: "admin",
      header: t(($) => $.models.attributes.user.admin),
      cell: ({ row }) => (row.original.admin ? "✓" : ""),
    },
    list.actionsColumn({
      edit: (row) => <EditLink to="/{-$locale}/admin/users/$id" params={{ id: String(row.id) }} />,
      describe: (row) => row.email ?? "",
    }),
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
        {...list.props}
        {...(searching && {
          data: found.data ?? [],
          total: found.data?.length ?? 0,
          isLoading: found.isLoading,
        })}
      />
    </Stack>
  );
}
