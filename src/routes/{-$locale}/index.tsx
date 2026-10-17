import { createFileRoute, redirect } from "@tanstack/react-router";
import {
  getActiveBannerOptions,
  listBlogPostsOptions,
  listCoursesOptions,
} from "@/client/@tanstack/react-query.gen";
import HomePage from "@/components/HomePage";
import { detectRootLocale, rememberLocale } from "@/lib/detect-locale";

// The home page at `/` (and `/ru`, `/es`). The loader prefetches into the
// request-scoped QueryClient so the course list is dehydrated into the SSR HTML.
export const Route = createFileRoute("/{-$locale}/")({
  beforeLoad: ({ params }) => {
    // Only the unprefixed root auto-detects; `/ru`/`/es` are explicit.
    if (params.locale === undefined) {
      const target = detectRootLocale();
      if (target) {
        throw redirect({ to: "/{-$locale}", params: { locale: target } });
      }
      rememberLocale("en");
    }
  },
  // The page shows the banner, the course grid and the latest three posts; all are
  // prefetched so they are in the server-rendered HTML.
  loader: async ({ context }) => {
    await Promise.all([
      context.queryClient.query({ ...listCoursesOptions(), staleTime: "static" }),
      // A banner is decoration: failing to read it must not take the page down.
      context.queryClient
        .query({ ...getActiveBannerOptions(), staleTime: "static" })
        .catch(() => null),
      context.queryClient.query({
        ...listBlogPostsOptions({ query: { perPage: 3 } }),
        staleTime: "static",
      }),
    ]);
  },
  component: HomePage,
});
