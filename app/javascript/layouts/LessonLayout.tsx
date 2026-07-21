import { AppShell } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import type { PropsWithChildren } from "react";
import { useLessonStore } from "@/pages/web/languages/lessons/show/store.tsx";
import NavbarBlock from "./blocks/NavbarBlock.tsx";

type Props = PropsWithChildren & {};

export default function LessonLayout({ children }: Props) {
  const [opened, { toggle }] = useDisclosure();
  const mobileNavOpened = useLessonStore((state) => state.mobileNavOpened);

  return (
    <AppShell
      h="100%"
      navbar={{
        width: { base: "100%", sm: 400, md: 450, lg: 550 },
        breakpoint: "sm",
        collapsed: { mobile: !mobileNavOpened },
      }}
      header={{ height: 60 }}
    >
      <AppShell.Header>
        <NavbarBlock onToggle={toggle} opened={opened} />
      </AppShell.Header>

      {children}
    </AppShell>
  );
}
