import { AppShell } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import type { PropsWithChildren } from "react";
import NavbarBlock from "./blocks/NavbarBlock.tsx";
import RootLayout from "./RootLayout.tsx";

type Props = PropsWithChildren & {};

export default function LessonLayout({ children }: Props) {
  const [opened, { toggle }] = useDisclosure();
  const [mobileOpened] = useDisclosure();

  return (
    <AppShell
      h="100%"
      navbar={{
        width: { xs: 300, sm: 400, md: 450, lg: 550 },
        breakpoint: "sm",
        collapsed: { mobile: !mobileOpened },
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
