import { type PropsWithChildren } from "react";
import { AppShell } from "@mantine/core";

import RootLayout from "./RootLayout.tsx";
import NavbarBlock from "./blocks/NavbarBlock.tsx";
import { useDisclosure } from "@mantine/hooks";

type Props = PropsWithChildren & {};

export default function LessonLayout({ children }: Props) {
  const [opened, { toggle }] = useDisclosure();
  const [mobileOpened, ] = useDisclosure();

  return (
    <RootLayout>
      <AppShell
        h="100%"
        navbar={{
          width: { xs: 300, sm: 400, md: 450, lg: 550 },
          breakpoint: "sm",
          collapsed: { mobile: !mobileOpened }
        }}
        header={{ height: 60 }}
      >
        <AppShell.Header>
          <NavbarBlock onToggle={toggle} opened={opened} />
        </AppShell.Header>

        {children}

      </AppShell>
    </RootLayout>
  );
}
