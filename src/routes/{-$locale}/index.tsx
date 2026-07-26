import { createFileRoute } from "@tanstack/react-router";
import { listCoursesOptions } from "@/client/@tanstack/react-query.gen";
import CourseCatalog from "@/components/CourseCatalog";

// Catalog at `/`. The loader prefetches into the request-scoped QueryClient so
// the course list is dehydrated into the SSR HTML.
export const Route = createFileRoute("/")({
  loader: ({ context }) =>
    context.queryClient.ensureQueryData(listCoursesOptions()),
  component: CourseCatalog,
});
