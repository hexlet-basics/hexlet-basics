import analytics from "@/lib/analytics";
import type { SharedProps } from "@/types";
import { Head, usePage } from "@inertiajs/react";
import parseHtml from "html-react-parser";
import { type PropsWithChildren, useEffect } from "react";

type Props = PropsWithChildren & {};
export default (props: Props) => {
  const page = usePage<SharedProps>();
  const { url } = page;
  const { auth, happendEvents: events, carrotQuestUserHash, metaTagsHTMLString } = page.props;

  const user = auth.user;

  // biome-ignore lint/correctness/useExhaustiveDependencies: <explanation>
  useEffect(() => {
    analytics.page({
      $pathname: url,
    });

    if (!user.guest) {
      // NOTE: This is a hack to access the Carrotquest plugin. Type Plugins contains only enable and disable properties,
      // but according to the documentation, we can call custom plugins methods this way.
      // https://github.com/DavidWells/analytics/issues/266
      // https://getanalytics.io/plugins/writing-plugins/#adding-custom-methods
      // @ts-expect-error
      analytics.plugins.carrotquest.auth(user.id, carrotQuestUserHash);
    }
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
