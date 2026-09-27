import { notifications } from "@mantine/notifications";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { getCourseLessonQueryKey, startLessonMutation } from "@/client/@tanstack/react-query.gen";

// The one way into a lesson that begins progress: issue the start command, and
// only once it succeeds go there. The Course page's Try/Continue and the
// player's Next both enter through it, so the two cannot drift apart.
//
// Progress begins through this command and nothing else: the router preloads on
// hover, so a lesson page that started itself on load would enroll a learner in
// every lesson they pointed at (ADR-0012). That is also why callers render a
// button rather than a link — a link would be preloaded, and could be opened
// around the command.
export function useEnterLesson(courseSlug: string) {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const start = useMutation({
    ...startLessonMutation(),
    // A refusal or a network failure leaves the visitor where they are; landing
    // in a lesson that was never started would be worse than staying put.
    onError: () => notifications.show({ message: t(($) => $.common.errors.network) }),
  });

  const enter = (lesson: { id: number; slug: string }) =>
    start.mutate(
      { path: { id: lesson.id } },
      {
        onSuccess: () => {
          // The lesson's payload may already be cached from a hover, taken
          // before this start — and before the pass that unlocked it. The
          // loader's queryClient.query would serve that stale copy, lock and all,
          // and an inactive query is not refetched by a plain invalidation, so
          // the entry is dropped and the loader reads it afresh.
          queryClient.removeQueries({
            queryKey: getCourseLessonQueryKey({ path: { courseSlug, slug: lesson.slug } }),
          });
          void navigate({
            to: "/{-$locale}/languages/$slug/lessons/$lessonSlug",
            // Keeps the URL locale of whichever page the learner came from.
            params: (prev) => ({ ...prev, slug: courseSlug, lessonSlug: lesson.slug }),
          });
        },
      },
    );

  return { enter, isPending: start.isPending };
}
