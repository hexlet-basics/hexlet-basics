import { Center, Pagination } from "@mantine/core";
import { Link } from "@tanstack/react-router";
import { pageCount } from "@/lib/page-count";

// The public pager, ported from legacy XPaging: Mantine Pagination whose pages
// And arrows are real links to `?page=N` on the current page, so every page is
// Reachable (and server-rendered) by URL. Hidden when there is only one page.
//
// The controls are rendered by Pagination itself, which only lets a caller
// Swap their root through getItemProps/getControlProps — so this is the one
// Place a router Link is passed as `component` rather than built with
// CreateLink. The link is relative to the current location and changes only
// The search, so nothing untyped (no path params) is threaded through.
// Page one is the bare URL, as legacy canonicalised the list.
const linkTo = (target: number) => ({
  component: Link,
  to: ".",
  search: target === 1 ? {} : { page: target },
});

export default function Pager({
  total,
  page,
  perPage,
}: {
  total: number;
  page: number;
  perPage: number;
}) {
  const last = pageCount(total, perPage);
  if (last <= 1) {
    return null;
  }

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
