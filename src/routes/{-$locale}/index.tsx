import { createFileRoute, redirect } from "@tanstack/react-router";
import { listCoursesOptions } from "@/client/@tanstack/react-query.gen";
import CourseCatalog from "@/components/CourseCatalog";
import { detectRootLocale } from "@/lib/detect-locale";

// Catalog at `/` (and `/ru`, `/es`). The loader prefetches into the
// request-scoped QueryClient so the course list is dehydrated into the SSR HTML.
export const Route = createFileRoute("/{-$locale}/")({
  beforeLoad: ({ params }) => {
    // Only the unprefixed root auto-detects; `/ru`/`/es` are explicit.
    if (params.locale === undefined) {
      const target = detectRootLocale();
      if (target) {
        throw redirect({ to: "/{-$locale}", params: { locale: target } });
      }
    }
  },
  loader: ({ context }) =>
    context.queryClient.ensureQueryData(listCoursesOptions()),
  component: CourseCatalog,
});
