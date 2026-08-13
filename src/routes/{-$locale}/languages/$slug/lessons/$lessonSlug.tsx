import { createFileRoute } from "@tanstack/react-router";
import { getCourseLessonOptions } from "@/client/@tanstack/react-query.gen";
import LessonPage from "@/components/lesson/LessonPage";

// The lesson player, at its legacy URL (ADR-0002) under the optional locale
// prefix.
//
// The loader prefetches the payload into the request-scoped QueryClient, so the
// theory — the page's indexable content and the first thing a learner reads — is
// in the server-rendered HTML rather than fetched after hydration (ADR-0008).
//
// Loading this page starts nothing. The router preloads on hover, so a read with
// that effect would enroll a learner in every lesson they merely pointed at;
// progress begins only through the start command (ADR-0012).
export const Route = createFileRoute("/{-$locale}/languages/$slug/lessons/$lessonSlug")({
  staticData: { chrome: "bare" },
  loader: ({ context, params }) =>
    context.queryClient.ensureQueryData(
      getCourseLessonOptions({ path: { courseSlug: params.slug, slug: params.lessonSlug } }),
    ),
  component: LessonRoute,
});

function LessonRoute() {
  const { slug, lessonSlug } = Route.useParams();
  return <LessonPage courseSlug={slug} lessonSlug={lessonSlug} />;
}
