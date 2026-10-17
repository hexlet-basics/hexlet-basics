import { ActionIcon, Group, Tooltip } from "@mantine/core";
import {
  type QueryKey,
  type UseMutationOptions,
  type UseQueryOptions,
  useQuery,
} from "@tanstack/react-query";
import type { OnChangeFn, PaginationState } from "@tanstack/react-table";
import { type ReactNode, useState } from "react";
import { useTranslation } from "react-i18next";
import { useDeleteConfirmation } from "@/hooks/useDeleteConfirmation";
import { useResourceMutation } from "@/hooks/useResourceMutation";
import type { CrudColumnDef } from "./CrudList";

// The envelope of every admin list response (`XxxPage`).
interface Page {
  items: { id: number }[];
  total: number;
  perPage: number;
}

interface ListQuery {
  page: number;
  // Left out until the user picks a size: the server's default is declared once in
  // the contract and echoed back in every response's `perPage`.
  perPage?: number;
}

// What the generated delete mutation takes: the row's id in the path.
interface DeleteVariables {
  path: { id: number };
}

interface CrudListOptions<
  TPage extends Page,
  TListError,
  TListKey extends QueryKey,
  TDeleteData,
  TDeleteError,
  TDeleteVars extends DeleteVariables,
> {
  // The generated list query options, e.g. `adminListBannersOptions`.
  list: (options: { query: ListQuery }) => UseQueryOptions<TPage, TListError, TPage, TListKey>;
  // The generated list query key, e.g. `adminListBannersQueryKey`: what a delete
  // invalidates.
  listKey: () => QueryKey;
  // The generated delete mutation, e.g. `adminDeleteBannerMutation`. Left out for
  // a resource that cannot be deleted; its rows then carry no delete button.
  remove?: () => UseMutationOptions<TDeleteData, TDeleteError, TDeleteVars>;
}

interface ActionsColumnOptions<TItem> {
  // The row's edit link: a typed `<EditLink to=... params=... />`. The route path
  // stays with the resource, where the router can check it.
  edit: (row: TItem) => ReactNode;
  // A resource's own row actions (a build, a review), shown between edit and
  // delete.
  extra?: (row: TItem) => ReactNode;
  // What the delete confirmation names, e.g. the banner's text.
  describe: (row: TItem) => ReactNode;
}

// Stands in for the delete mutation of a resource that has none; no button ever
// runs it.
const noMutation = {};

// The list engine of an admin index route: it owns the pagination state, runs the
// generated list query, and (when the resource can be deleted) the confirmed
// delete with its toasts and the list's invalidation. A route keeps only its
// title and columns, and hands in the generated hooks it uses.
//
// `props` is spread into `CrudList`; `actionsColumn` builds the shared last
// column (edit link and delete button); `query` is the list's query state, for a
// route that mixes in a second source of rows (the users search).
export function useCrudList<
  TPage extends Page,
  TListError,
  TListKey extends QueryKey,
  TDeleteData,
  TDeleteError,
  TDeleteVars extends DeleteVariables,
>({
  list,
  listKey,
  remove,
}: CrudListOptions<TPage, TListError, TListKey, TDeleteData, TDeleteError, TDeleteVars>) {
  type TItem = TPage["items"][number];
  const { t } = useTranslation();
  const confirmDelete = useDeleteConfirmation();

  const [pageIndex, setPageIndex] = useState(0);
  const [chosenSize, setChosenSize] = useState<number>();

  const query = useQuery(list({ query: { page: pageIndex + 1, perPage: chosenSize } }));

  const deleteMutation = useResourceMutation({
    mutation: remove?.() ?? noMutation,
    invalidate: [listKey()],
    successMessage: t(($) => $.admin.crud.deleted),
    errorMessage: t(($) => $.admin.crud.deleteError),
  });

  // The size is 0 until the first response says what the server used; the table
  // renders no pager before there is data anyway.
  const pagination: PaginationState = {
    pageIndex,
    pageSize: chosenSize ?? query.data?.perPage ?? 0,
  };
  const onPaginationChange: OnChangeFn<PaginationState> = (updater) => {
    const next = typeof updater === "function" ? updater(pagination) : updater;
    setPageIndex(next.pageIndex);
    if (next.pageSize !== pagination.pageSize) setChosenSize(next.pageSize);
  };

  const actionsColumn = ({
    edit,
    extra,
    describe,
  }: ActionsColumnOptions<TItem>): CrudColumnDef<TItem> => ({
    id: "actions",
    header: "",
    cell: ({ row }) => (
      <Group gap="xs" justify="flex-end" wrap="nowrap">
        {edit(row.original)}
        {extra?.(row.original)}
        {remove && (
          <Tooltip label={t(($) => $.admin.crud.delete)}>
            <ActionIcon
              color="red"
              variant="light"
              aria-label={t(($) => $.admin.crud.delete)}
              onClick={() =>
                confirmDelete({
                  description: describe(row.original),
                  onConfirm: () =>
                    deleteMutation.mutate({ path: { id: row.original.id } } as TDeleteVars),
                })
              }
            >
              ✕
            </ActionIcon>
          </Tooltip>
        )}
      </Group>
    ),
  });

  return {
    query,
    props: {
      data: (query.data?.items ?? []) as TItem[],
      total: query.data?.total ?? 0,
      pagination,
      onPaginationChange,
      isLoading: query.isLoading,
    },
    actionsColumn,
  };
}
