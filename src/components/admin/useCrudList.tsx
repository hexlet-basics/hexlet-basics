import { ActionIcon, Group, Tooltip } from "@mantine/core";
import {
  type QueryKey,
  type UseMutationOptions,
  type UseQueryOptions,
  useQuery,
} from "@tanstack/react-query";
import type { OnChangeFn, PaginationState } from "@tanstack/react-table";
import { useDebouncedCallback } from "@mantine/hooks";
import { useNavigate, useSearch } from "@tanstack/react-router";
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
  // The text filter, on the lists whose operation takes one.
  q?: string;
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
  // Show a search box that filters the list on the server (`q`). Only for a list
  // whose operation declares `q`.
  searchable?: boolean;
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
  searchable = false,
}: CrudListOptions<TPage, TListError, TListKey, TDeleteData, TDeleteError, TDeleteVars>) {
  type TItem = TPage["items"][number];
  const { t } = useTranslation();
  const confirmDelete = useDeleteConfirmation();

  // The page, its size and the text filter live in the address, as the route's
  // generated list-query schema (`validateSearch`) declares them, so a list page
  // can be shared and survives a reload. The hook serves every index route, so it
  // reads the search of whichever route it is mounted in.
  const search: ListQuery = useSearch({ strict: false });
  const navigate = useNavigate();
  const setSearch = (patch: Partial<ListQuery>) =>
    navigate({ to: ".", search: (prev) => ({ ...prev, ...patch }) });

  // What the box shows follows the keystrokes; the address (and so the query)
  // follows it a beat later. A new term starts from the first page: the old page
  // may not exist in the narrowed list.
  const [text, setText] = useState(search.q ?? "");
  const applyTerm = useDebouncedCallback(
    (term: string) => setSearch({ q: term === "" ? undefined : term, page: undefined }),
    300,
  );

  const query = useQuery(
    list({
      query: {
        page: search.page ?? 1,
        perPage: search.perPage,
        q: searchable ? search.q : undefined,
      },
    }),
  );

  const deleteMutation = useResourceMutation({
    mutation: remove?.() ?? noMutation,
    invalidate: [listKey()],
    successMessage: t(($) => $.admin.crud.deleted),
    errorMessage: t(($) => $.admin.crud.deleteError),
  });

  // The size is 0 until the first response says what the server used; the table
  // renders no pager before there is data anyway.
  const pagination: PaginationState = {
    pageIndex: (search.page ?? 1) - 1,
    pageSize: search.perPage ?? query.data?.perPage ?? 0,
  };
  const onPaginationChange: OnChangeFn<PaginationState> = (updater) => {
    const next = typeof updater === "function" ? updater(pagination) : updater;
    setPageIndex(next.pageIndex);
    if (next.pageSize !== pagination.pageSize) {
      setChosenSize(next.pageSize);
    }
  };

  // Typed through the page, not read off it: a property read on a type parameter
  // widens to its constraint, and the rows would lose their own type.
  const rows: TPage["items"] | undefined = query.data?.["items"];

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
                  onConfirm: () => deleteMutation.mutate({ path: { id: row.original.id } }),
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
      data: rows,
      total: query.data?.total ?? 0,
      pagination,
      onPaginationChange,
      isLoading: query.isLoading,
      search: searchable
        ? {
            value: text,
            onChange: (value: string) => {
              setText(value);
              applyTerm(value.trim());
            },
          }
        : undefined,
    },
    actionsColumn,
  };
}
