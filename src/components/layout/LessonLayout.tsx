import { AppShell } from "@mantine/core";
import type { PropsWithChildren } from "react";
import Header from "./Header";

// The lesson player's chrome: the site header and nothing else — no footer, no
// floating chat button. The player is a two-pane workspace that has to fill the
// viewport exactly, and page furniture below it would either push the panes off
// screen or make the page scroll behind them.
//
// The main area is a full-viewport box (Mantine sets `box-sizing: border-box`
// globally, so the header offset AppShell adds as padding comes out of the same
// 100dvh) — which is what gives the splitter a definite height to divide.
export default function LessonLayout({ children }: PropsWithChildren) {
  return (
    <AppShell header={{ height: 60 }}>
      <AppShell.Header>
        <Header />
      </AppShell.Header>

      <AppShell.Main h="100dvh" p={0} pt={60}>
        {children}
      </AppShell.Main>
    </AppShell>
  );
}
