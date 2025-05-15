import analytics from "@/lib/analytics";
import type { SharedProps } from "@/types";
import { Head, usePage } from "@inertiajs/react";
import parseHtml from "html-react-parser";
import { type PropsWithChildren, useEffect } from "react";

type Props = PropsWithChildren & {};
export default (props: Props) => {
  const page = usePage<SharedProps>();
  const { url } = page;
  const { auth, happendEvents: events, metaTagsHTMLString } = page.props;

  // const user = auth.user;

  // biome-ignore lint/correctness/useExhaustiveDependencies: <explanation>
  useEffect(() => {

    analytics.page()

    if (events) {
      for (const event of events) {
        switch (event.type) {
          case "UserSignedInEvent":
            // TODO: Add type for event name to avoid sending wrong events
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
          case "LessonStartedEvent":
            analytics.track("lesson_started", event.data);
            break;
          default:
            break;
        }
      }
    }
  }, []);

  return (
    <>
      <Head>{parseHtml(metaTagsHTMLString, { trim: true })}</Head>
      {props.children}
    </>
  );
};
