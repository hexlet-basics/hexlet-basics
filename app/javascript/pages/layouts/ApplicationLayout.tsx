import type { PropsWithChildren } from "react";
import { Container } from "react-bootstrap";

import XFlash from "@/components/XFlash.tsx";
import { XBreadcrumb } from "@/components/breadcrumbs.tsx";
import type { BreadcrumbItem, SharedProps } from "@/types/index.js";
import RootLayout from "./RootLayout.tsx";
// import TgContestBanner from "./banners/tg_contest_banner/TgContestBanner.tsx";
import FooterBlock from "./blocks/FooterBlock.tsx";
import NavbarBlock from "./blocks/NavbarBlock.tsx";
import { usePage } from "@inertiajs/react";
import ContactMethodRequestingBlock from "./blocks/ContactMethodRequestingBlock.tsx";

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

  return (
    <RootLayout>
      <Container className="py-2">
        <NavbarBlock className="pb-3 border-bottom" />
        <XFlash />
        {shouldAddContactMethod && (
          <div className="mt-3"><ContactMethodRequestingBlock /></div>
        )}
      </Container>
      {(items || header) && (
        <Container className="my-2 my-md-4">
          {items && <XBreadcrumb items={items} />}
          {header && <h1 className={center ? "text-center" : ""}>{header}</h1>}
        </Container>
      )}
      <div className="pb-5">{children}</div>
      <FooterBlock />
    </RootLayout>
  );
}
