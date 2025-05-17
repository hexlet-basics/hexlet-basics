import analytics from "@/lib/analytics";
import type { SharedProps } from "@/types";
import { Head, usePage } from "@inertiajs/react";
import parseHtml from "html-react-parser";
import { type PropsWithChildren, useEffect } from "react";

type Props = PropsWithChildren & {};
export default (props: Props) => {
  const page = usePage<SharedProps>();
  const { url } = page;
  const { auth, happendEvents, metaTagsHTMLString } = page.props;

  // const user = auth.user;

  useEffect(() => {

    analytics.page()

    if (happendEvents) {
      for (const event of happendEvents) {
        switch (event.type) {
          case "UserSignedInEvent":
            analytics.identify(event.data.id.toString(), event.data);
            analytics.track("signed_in", event.data);
            break;
          case "UserSignedUpEvent":
            analytics.identify(event.data.id.toString(), event.data);
            analytics.track("signed_up", event.data);
            break;
          case "CourseStartedEvent":
            analytics.track("course_started", event.data);
            break;
          case "LeadCreatedEvent":
            analytics.track("lead_created", event.data);
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
