import { createFileRoute } from "@tanstack/react-router";
import { truncate } from "es-toolkit/compat";
import { getCourseLessonOptions } from "@/client/@tanstack/react-query.gen";
import type { CourseLessonView } from "@/client/types.gen";
import LessonPage, { LessonMissing } from "@/components/lesson/LessonPage";
import { seoHead } from "@/lib/seo-head";

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
  // Legacy lessons#show meta, composed from the payload: a title naming the
  // lesson and the course's landing copy, a description drawn from the theory,
  // the canonical link, and Open Graph as an article with the course cover.
  // Legacy gave this page no og:description and no Twitter card, so the shared
  // social block is off and only the Open Graph block is asked for.
  head: ({ loaderData, match }) => {
    if (!loaderData) return {};
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
  // the apology for it lives here rather than in a branch the page cannot reach.
  errorComponent: LessonMissing,
  component: LessonRoute,
});

function LessonRoute() {
  const { slug, lessonSlug } = Route.useParams();
  return <LessonPage courseSlug={slug} lessonSlug={lessonSlug} />;
}

// Legacy's description: `[<version name>] — <lesson> — <theory>`, cut by Rails'
// `truncate(length: 220)` — 217 characters and an ellipsis, mid-word — and then
// whitespace-squashed by meta-tags. The theory is the raw markdown, as legacy
// used it. es-toolkit's truncate counts code points once a string holds any
// astral or combining character, the way Ruby counts characters.
function lessonDescription({ lesson }: CourseLessonView): string {
  const version = lesson.course.currentVersion?.name ?? "";
  const text = `[${version}] — ${lesson.name ?? ""} — ${lesson.theory ?? ""}`;
  return squish(truncate(text, { length: 220 }));
}

// Rails' `squish`: every whitespace run to one space, both ends trimmed.
// es-toolkit has no counterpart, so this one line stays.
function squish(text: string): string {
  return text.replace(/\s+/g, " ").trim();
}
