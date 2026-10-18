import { ActionIcon, Affix, AppShell, Tooltip } from "@mantine/core";
import { IconSend } from "@tabler/icons-react";
import type { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import Footer from "./Footer";
import Header from "./Header";

// Application shell, ported from legacy ApplicationLayout. Wraps every public
// page with a fixed-height header and a footer; the floating community-chat
// button mirrors the legacy Affix. Mounted once in the `{-$locale}` layout
// route so it covers the catalog, course, and auth pages without per-page
// duplication.
export default function ApplicationLayout({ children }: PropsWithChildren) {
  const { t } = useTranslation();

  return (
    <>
      <Affix position={{ bottom: 20, right: 20 }}>
        <Tooltip label={t(($) => $.common.community_chat)}>
          <ActionIcon
            component="a"
            href={t(($) => $.common.community_url)}
            aria-label={t(($) => $.common.community_chat)}
            variant="filled"
            size="xl"
            radius="xl"
            target="_blank"
            rel="noreferrer"
          >
            <IconSend aria-hidden="true" />
          </ActionIcon>
        </Tooltip>
      </Affix>

      <AppShell header={{ height: 60 }}>
        <AppShell.Header>
          <Header />
        </AppShell.Header>

        <AppShell.Main>
          {children}
          <Footer />
        </AppShell.Main>
      </AppShell>
    </>
  );
}
