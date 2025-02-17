import type { SharedProps } from "@/types";
import { usePage } from "@inertiajs/react";
import { type PropsWithChildren, useEffect } from "react";

// import { usePostHog } from "posthog-js/react";
import analytics from "@/analytics";

type Props = PropsWithChildren & {};
export default (props: Props) => {
  const { events } = usePage<SharedProps>().props;
  // const posthog = usePostHog();

  // biome-ignore lint/correctness/useExhaustiveDependencies: <explanation>
  useEffect(() => {
    if (events) {
      for (const event of events) {
        switch (event.type) {
          case "UserSignedInEvent":
            // add carrot quest, ym, ga
            analytics.identify(event.data.id, {
              email: event.data.email,
            });
            break;
          case "UserSignedUpEvent":
            analytics.identify(event.data.id, {
              email: event.data.email,
              first_name: event.data.first_name,
            });
            break;
          case "CourseStartedEvent":
            // posthog?.capture('course_started', {
            //   slug: event.data.slug,
            //   locale: event.data.locale,
            // });
            analytics.track("course_started", {
              slug: event.data.slug,
              locale: event.data.locale,
            });
            break;
          case "CourseFinishedEvent":
            break;
          case "LessonStartedEvent":
            analytics.track("lesson_started", {
              course_slug: event.data.course_slug,
              lesson_slug: event.data.lesson_slug,
              locale: event.data.locale,
            });
            // posthog?.capture('lesson_started', {
            //   course_slug: event.data.course_slug,
            //   lesson_slug: event.data.lesson_slug,
            //   locale: event.data.locale,
            // });
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
