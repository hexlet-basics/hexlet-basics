import { createFileRoute, redirect } from "@tanstack/react-router";

// The admin root has no dashboard (the legacy home#index is not carried over):
// land on the courses list, the back office's main working surface.
export const Route = createFileRoute("/{-$locale}/admin/")({
  beforeLoad: ({ params }) => {
    throw redirect({ to: "/{-$locale}/admin/courses", params });
  },
});
