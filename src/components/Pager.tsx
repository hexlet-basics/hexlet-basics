import { Center, Pagination } from "@mantine/core";
import { Link } from "@tanstack/react-router";

// The public pager, ported from legacy XPaging: Mantine Pagination whose pages
// and arrows are real links to `?page=N` on the current page, so every page is
// reachable (and server-rendered) by URL. Hidden when there is only one page.
//
// The controls are rendered by Pagination itself, which only lets a caller
// swap their root through getItemProps/getControlProps — so this is the one
// place a router Link is passed as `component` rather than built with
// createLink. The link is relative to the current location and changes only
// the search, so nothing untyped (no path params) is threaded through.
export default function Pager({
  total,
  page,
  perPage,
}: {
  total: number;
  page: number;
  perPage: number;
}) {
  const last = Math.max(1, Math.ceil(total / perPage));
  if (last <= 1) return null;

  // Page one is the bare URL, as legacy canonicalised the list.
  const linkTo = (target: number) => ({
    component: Link,
    to: ".",
    search: target === 1 ? {} : { page: target },
  });

  const controlPage = {
    first: 1,
    previous: page > 1 ? page - 1 : undefined,
    next: page < last ? page + 1 : undefined,
    last,
  };

  return (
    <Center>
      <Pagination
        my="xl"
        total={last}
        value={page}
        getItemProps={linkTo}
        getControlProps={(control) => {
          const target = controlPage[control];
          return target ? linkTo(target) : {};
        }}
        boundaries={1}
        siblings={1}
      />
    </Center>
  );
}
