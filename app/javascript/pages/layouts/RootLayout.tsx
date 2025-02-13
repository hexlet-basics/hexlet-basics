import type { SharedProps } from "@/types";
import { usePage } from "@inertiajs/react";
import { usePostHog } from "posthog-js/react";
import { useEffect, type PropsWithChildren } from "react";

type Props = PropsWithChildren & {};
export default (props: Props) => {
  const { events } = usePage<SharedProps>().props;
  const posthog = usePostHog();

  // biome-ignore lint/correctness/useExhaustiveDependencies: <explanation>
  useEffect(() => {
    if (events) {
      for (const event of events) {
        switch (event.type) {
          case "UserSignedInEvent":
            // add carrot quest, ym, ga
            posthog?.identify(event.data.id, {
              email: event.data.email,
            });
            break;
          case "UserSignedUpEvent":
            posthog?.capture("signed_up");
            posthog?.identify(event.data.id, {
              email: event.data.email,
            });
            break;
          case "CourseStartedEvent":
            posthog?.capture("course_started", {
              slug: event.data.slug,
              locale: event.data.locale,
            });
            break;
          case "CourseFinishedEvent":
            break;
          case "LessonStartedEvent":
            posthog?.capture("lesson_started", {
              course_slug: event.data.course_slug,
              lesson_slug: event.data.lesson_slug,
              locale: event.data.locale,
            });
            break;
          case "LessonFinishedEvent":
            break;
          default:
            break;
        }
      }
    }
  }, []);

  return props.children;
};
