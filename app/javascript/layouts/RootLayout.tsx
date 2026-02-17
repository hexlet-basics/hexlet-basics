import { Head, router, usePage } from "@inertiajs/react";
import {
  Center,
  Container,
  MantineProvider,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { ModalsProvider } from "@mantine/modals";
import * as Sentry from "@sentry/react";
// import { ColorSchemeScript } from '@mantine/core';
import parseHtml from "html-react-parser";
import { type PropsWithChildren, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { useDebugAppData } from "@/hooks/useDebugAppData";
import ahoy from "@/lib/ahoy";
import analytics, { processHappendEvents } from "@/lib/analytics";
import { resolver, theme } from "@/lib/mantine";
import { fromWindow } from "@/lib/utils";

type Props = PropsWithChildren & {};

function FallbackComponent() {
  const { t } = useTranslation();

  return (
    <Container>
      <Center h="100vh">
        <Stack align="center">
          <Title order={1}>{t(($) => $.layouts.root.fallback.header)}</Title>
          <Text>{t(($) => $.layouts.root.fallback.description)}</Text>
        </Stack>
      </Center>
    </Container>
  );
}

function DebugAppData() {
  useDebugAppData();

  return null;
}

function RootLayout(props: Props) {
  const {
    props: { metaTagsHTMLString },
    flash,
  } = usePage();

  useEffect(() => {
    const ymClientId = fromWindow("ymClientId");
    ahoy.trackView({ ym_client_id: ymClientId });
  }, []);

  useEffect(() => {
    const unlisten = router.on("navigate", () => {
      analytics.page();
    });

    return () => unlisten();
  }, []);

  useEffect(() => {
    if (flash.events) {
      processHappendEvents(flash.events);
    }
  }, [flash.events]);

  return (
    <MantineProvider cssVariablesResolver={resolver} theme={theme}>
      <ModalsProvider>
        <Sentry.ErrorBoundary fallback={FallbackComponent} showDialog>
          <Head>{parseHtml(metaTagsHTMLString, { trim: true })}</Head>
          {import.meta.env.DEV && <DebugAppData />}
          {props.children}
        </Sentry.ErrorBoundary>
      </ModalsProvider>
    </MantineProvider>
  );
}

export default Sentry.withProfiler(RootLayout);
