import type { SharedProps } from "@/types";
import { usePage } from "@inertiajs/react";
import { usePostHog } from "posthog-js/react";
import { type PropsWithChildren, useEffect } from "react";

// import { usePostHog } from "posthog-js/react";
import analytics from "@/analytics";

import { useEffect, type PropsWithChildren } from "react";

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
            analytics.identify(event.data.id.toString(), {
              email: event.data.email
            });
            break;
          case "UserSignedUpEvent":
<<<<<<< HEAD
            posthog?.capture("signed_up");
            posthog?.identify(event.data.id, {
=======
            analytics.identify(event.data.id.toString(), {
>>>>>>> wip
              email: event.data.email,
              first_name: event.data.first_name,
            });
            break;
          case "CourseStartedEvent":
<<<<<<< HEAD
            posthog?.capture("course_started", {
=======
            // posthog?.capture('course_started', {
            //   slug: event.data.slug,
            //   locale: event.data.locale,
            // });
            analytics.track('course_started', {
>>>>>>> wip
              slug: event.data.slug,
              locale: event.data.locale,
            });
            break;
          case "CourseFinishedEvent":
            break;
          case "LessonStartedEvent":
<<<<<<< HEAD
            posthog?.capture("lesson_started", {
=======
            analytics.track('lesson_started', {
>>>>>>> wip
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
