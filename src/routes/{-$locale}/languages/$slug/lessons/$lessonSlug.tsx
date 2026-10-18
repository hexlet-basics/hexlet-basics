import { createFileRoute } from "@tanstack/react-router";
import { truncate } from "es-toolkit/compat";
import { getCourseLessonOptions } from "@/client/@tanstack/react-query.gen";
import type { CourseLessonView } from "@/client/types.gen";
import LessonPage, { LessonMissing } from "@/components/lesson/LessonPage";
import { seoHead } from "@/lib/seo-head";

// The lesson player, at its legacy URL (ADR-0002) under the optional locale
// Prefix.
//
// The loader prefetches the payload into the request-scoped QueryClient, so the
// Theory — the page's indexable content and the first thing a learner reads — is
// In the server-rendered HTML rather than fetched after hydration (ADR-0008).
//
// Loading this page starts nothing. The router preloads on hover, so a read with
// That effect would enroll a learner in every lesson they merely pointed at;
// Progress begins only through the start command (ADR-0012).
export const Route = createFileRoute("/{-$locale}/languages/$slug/lessons/$lessonSlug")({
  staticData: { chrome: "bare" },
  loader: ({ context, params }) =>
    // The payload carries the visitor's progress: default staleness, not static.
    context.queryClient.query(
      getCourseLessonOptions({ path: { courseSlug: params.slug, slug: params.lessonSlug } }),
    ),
  // Legacy lessons#show meta, composed from the payload: a title naming the
  // Lesson and the course's landing copy, a description drawn from the theory,
  // The canonical link, and Open Graph as an article with the course cover.
  // Legacy gave this page no og:description and no Twitter card, so the shared
  // Social block is off and only the Open Graph block is asked for.
  head: ({ loaderData, match }) => {
    if (!loaderData) {
      return {};
    }
    const { i18n } = match.context;
    const { lesson, landingPage } = loaderData;
    // The copy is a YAML block scalar ending in a newline, which legacy squished.
    const title = squish(
      i18n.t(($) => $.courses.lessons.show.title, {
        lesson_name: lesson.name ?? "",
        language_name: landingPage?.name ?? "",
      }),
    );
    return seoHead({
      i18n,
      title,
      description: lessonDescription(loaderData),
      canonicalPath: match.pathname,
      image: lesson.course.coverListVariant,
      social: false,
      openGraph: { type: "article", locale: i18n.language },
    });
  },
  // A slug that resolves to nothing rejects in the loader, above the page, so
  // The apology for it lives here rather than in a branch the page cannot reach.
  errorComponent: LessonMissing,
  component: LessonRoute,
});

function LessonRoute() {
  const { slug, lessonSlug } = Route.useParams();
  const { user } = Route.useRouteContext();
  return <LessonPage courseSlug={slug} lessonSlug={lessonSlug} signedIn={Boolean(user)} />;
}

// Legacy's description: `[<version name>] — <lesson> — <theory>`, cut by Rails'
// `truncate(length: 220)` — 217 characters and an ellipsis, mid-word — and then
// Whitespace-squashed by meta-tags. The theory is the raw markdown, as legacy
// Used it. es-toolkit's truncate counts code points once a string holds any
// Astral or combining character, the way Ruby counts characters.
function lessonDescription({ lesson }: CourseLessonView): string {
  const version = lesson.course.currentVersion?.name ?? "";
  const text = `[${version}] — ${lesson.name ?? ""} — ${lesson.theory ?? ""}`;
  return squish(truncate(text, { length: 220 }));
}

// Rails' `squish`: every whitespace run to one space, both ends trimmed.
// Es-toolkit has no counterpart, so this one line stays.
function squish(text: string): string {
  return text.replaceAll(/\s+/g, " ").trim();
}
