import * as Routes from "@/routes.js";
import { useEffect, type PropsWithChildren } from "react";
import { AppShell, Container, Title, Stack, Box } from '@mantine/core';
import * as CookieConsent from "vanilla-cookieconsent";
import cookieTranslations from '@/locales/cookie_consent.ts'

import { XBreadcrumb } from "@/components/breadcrumbs.tsx";
import type { BreadcrumbItem, SharedProps } from "@/types/index.js";
import RootLayout from "./RootLayout.tsx";
// import TgContestBanner from "./banners/tg_contest_banner/TgContestBanner.tsx";
import FooterBlock from "./blocks/FooterBlock.tsx";
import NavbarBlock from "./blocks/NavbarBlock.tsx";
import { usePage } from "@inertiajs/react";
import { useDisclosure } from "@mantine/hooks";
import XFlash from "@/components/XFlash.tsx";
import ContactMethodRequestingBlock from "./blocks/ContactMethodRequestingBlock.tsx";
import { isCurrentUrl } from "@/lib/utils.ts";

type Props = PropsWithChildren & {
  header?: string;
  center?: boolean;
  items?: BreadcrumbItem[];
};

export default function ApplicationLayout({
  children,
  header,
  items,
  center = false,
}: Props) {
  const page = usePage<SharedProps>();
  const { props: { shouldAddContactMethod } } = page
  const [opened, { toggle }] = useDisclosure();

  useEffect(() => {
    CookieConsent.run({
      categories: {
        necessary: {
          enabled: true,  // this category is enabled by default
          readOnly: true  // this category cannot be disabled
        },
        analytics: {}
      },
      language: {
        autoDetect: "document",
        default: 'ru',
        translations: cookieTranslations,
      }
    });
  }, []);

  return (
    <RootLayout>
      <AppShell
        header={{ height: 60 }}
      >
        <AppShell.Header>
          <NavbarBlock onToggle={toggle} opened={opened} />
        </AppShell.Header>

        <AppShell.Main>
          <XFlash />
          {shouldAddContactMethod && isCurrentUrl(Routes.new_lead_path()) && (
            <Box mt="lg"><ContactMethodRequestingBlock /></Box>
          )}
          {(items || header) && (
            <Container size="lg" my="xl">
              <Stack>
                {items && <XBreadcrumb items={items} />}
                {header && (
                  <Title order={1} ta={center ? "center" : "left"}>
                    {header}
                  </Title>
                )}
              </Stack>
            </Container>
          )}
          {children}
        </AppShell.Main>

        {/* <AppShell.Footer> */}
        <FooterBlock />
        {/* </AppShell.Footer> */}
      </AppShell>
    </RootLayout>
  );
}
