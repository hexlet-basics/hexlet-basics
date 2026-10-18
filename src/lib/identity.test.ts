import { QueryClient } from "@tanstack/react-query";
import { expect, test } from "vitest";
import { getCurrentUserQueryKey } from "@/client/@tanstack/react-query.gen";
import type { User } from "@/client/types.gen";
import { onIdentityChanged } from "@/lib/identity";

const user = { id: 7, name: "Ada" } as User;

// A router that only records what the cache held when it was asked to reload.
function routerOver(queryClient: QueryClient) {
  const seen: { cacheSizeAtInvalidate?: number } = {};
  return {
    seen,
    router: {
      invalidate: () => {
        seen.cacheSizeAtInvalidate = queryClient.getQueryCache().getAll().length;
        return Promise.resolve();
      },
    },
  };
}

test("drops everything read as the previous visitor before the pages are read again", async () => {
  const queryClient = new QueryClient();
  queryClient.setQueryData(["a course read as a guest"], { progress: 0 });
  const { router, seen } = routerOver(queryClient);

  await onIdentityChanged(queryClient, router);

  expect(queryClient.getQueryData(["a course read as a guest"])).toBeUndefined();
  expect(seen.cacheSizeAtInvalidate).toBe(0);
});

test("keeps the user the flow already has, so the page does not fetch them again", async () => {
  const queryClient = new QueryClient();
  const { router, seen } = routerOver(queryClient);

  await onIdentityChanged(queryClient, router, user);

  expect(queryClient.getQueryData(getCurrentUserQueryKey())).toStrictEqual({ user });
  expect(seen.cacheSizeAtInvalidate).toBe(1);
});

test("can be told there is no router to reload, for a flow that is itself a navigation", async () => {
  const queryClient = new QueryClient();
  queryClient.setQueryData(["stale"], 1);

  await onIdentityChanged(queryClient, null, user);

  expect(queryClient.getQueryData(["stale"])).toBeUndefined();
  expect(queryClient.getQueryData(getCurrentUserQueryKey())).toStrictEqual({ user });
});
