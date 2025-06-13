import analytics, { processHappendEvents } from "@/lib/analytics";

import type { SharedProps } from "@/types";
import { Head, router, usePage } from "@inertiajs/react";
import { ColorSchemeScript } from "@mantine/core";
import parseHtml from "html-react-parser";
import { type PropsWithChildren, useEffect } from "react";

type Props = PropsWithChildren & {};

export default (props: Props) => {
  const page = usePage<SharedProps>();
  const { happendEvents, metaTagsHTMLString } = page.props;

  useEffect(() => {
    const unlisten = router.on('navigate', (event) => {
      analytics.page()

      if (happendEvents) {
        processHappendEvents(happendEvents)
      }
    },)

    return () => unlisten()
  }, [happendEvents])

  return (
    <>
      <Head>{parseHtml(metaTagsHTMLString, { trim: true })}</Head>
      <Head><ColorSchemeScript /></Head>
      {props.children}
    </>
  );
};
