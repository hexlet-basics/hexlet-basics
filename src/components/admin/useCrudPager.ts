import type { OnChangeFn, PaginationState } from "@tanstack/react-table";
import { useState } from "react";

// Pagination state of an admin list. The page size is the server's: the contract
// declares the default once and every response echoes the size it used in
// `perPage`, so the table's initial state reads it from the first response and
// requests omit `perPage` until the user picks a size of their own.
export function useCrudPager() {
  const [pageIndex, setPageIndex] = useState(0);
  const [chosenSize, setChosenSize] = useState<number>();

  return {
    // Spread into the list hook's `query`.
    query: { page: pageIndex + 1, perPage: chosenSize },
    // The size is 0 until the first response says what the server used; the
    // table renders no pager before there is data anyway.
    props: (serverPerPage: number | undefined) => {
      const pagination: PaginationState = {
        pageIndex,
        pageSize: chosenSize ?? serverPerPage ?? 0,
      };
      const onPaginationChange: OnChangeFn<PaginationState> = (updater) => {
        const next = typeof updater === "function" ? updater(pagination) : updater;
        setPageIndex(next.pageIndex);
        if (next.pageSize !== pagination.pageSize) setChosenSize(next.pageSize);
      };
      return { pagination, onPaginationChange };
    },
  };
}
