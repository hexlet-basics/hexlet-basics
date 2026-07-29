import { createFileRoute } from "@tanstack/react-router";
import { listCoursesOptions } from "@/client/@tanstack/react-query.gen";
import CourseCatalog from "@/components/CourseCatalog";

// Legacy catalog URL `/languages`, kept for backward compatibility (ADR-0002).
export const Route = createFileRoute("/{-$locale}/languages/")({
  loader: ({ context }) => context.queryClient.ensureQueryData(listCoursesOptions()),
  component: CourseCatalog,
});
