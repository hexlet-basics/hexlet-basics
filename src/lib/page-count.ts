// How many pages `total` rows make at `perPage` a page. An empty list is still
// one (empty) page, so a pager never shows "page 1 of 0"; a page size below one
// (the admin list's size before the first response says what the server used) is
// read as one row a page rather than dividing by zero. The public pager and the
// admin list's pager both count pages here.
export function pageCount(total: number, perPage: number): number {
  return Math.max(1, Math.ceil(total / Math.max(1, perPage)));
}
