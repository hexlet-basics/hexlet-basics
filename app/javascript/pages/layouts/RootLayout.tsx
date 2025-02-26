import analytics from "@/analytics";
import type { SharedProps } from "@/types";
import { usePage } from "@inertiajs/react";
import { type PropsWithChildren, useEffect } from "react";

type Props = PropsWithChildren & {};
export default (props: Props) => {
  const { auth, events } = usePage<SharedProps>().props;
  const user = auth.user;

  // biome-ignore lint/correctness/useExhaustiveDependencies: <explanation>
  useEffect(() => {
    if (!user.guest) {
      // biome-ignore lint/suspicious/noExplicitAny: This is a hack to access the Carrotquest plugin
      const plugins = analytics.plugins as any;
      plugins.carrotquest.auth(user.id, user.carrotquest_hash);
    }
    if (events) {
      for (const event of events) {
        switch (event.type) {
          case "UserSignedInEvent":
            analytics.track("signed_in", event.data);
            analytics.identify(event.data.id.toString(), event.data);
            break;
          case "UserSignedUpEvent":
            analytics.track("signed_up", event.data);
            analytics.identify(event.data.id.toString(), event.data);
            break;
          case "CourseStartedEvent":
            analytics.track("course_started", event.data);
            break;
          case "CourseFinishedEvent":
            break;
          case "LessonStartedEvent":
            analytics.track("lesson_started", event.data);
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
