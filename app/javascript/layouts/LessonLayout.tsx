import { AppShell } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import type { PropsWithChildren } from "react";
import NavbarBlock from "./blocks/NavbarBlock.tsx";

type Props = PropsWithChildren & {};

export default function LessonLayout({ children }: Props) {
  const [opened, { toggle }] = useDisclosure();

  return (
    <AppShell h="100%" header={{ height: 60 }}>
      <AppShell.Header>
        <NavbarBlock onToggle={toggle} opened={opened} />
      </AppShell.Header>

      {children}
    </AppShell>
  );
}
