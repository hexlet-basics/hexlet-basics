import type { PropsWithChildren } from "react";
import { Container } from "react-bootstrap";

import XFlash from "@/components/XFlash.tsx";
import { XBreadcrumb } from "@/components/breadcrumbs.tsx";
import type { BreadcrumbItem } from "@/types/index.js";
import RootLayout from "./RootLayout.tsx";
// import TgContestBanner from "./banners/tg_contest_banner/TgContestBanner.tsx";
import FooterBlock from "./blocks/FooterBlock.tsx";
import NavbarBlock from "./blocks/NavbarBlock.tsx";

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
  return (
    <RootLayout>
      {/* <TgContestBanner /> */}
      <Container className="py-2">
        <NavbarBlock className="pb-3 border-bottom" />
        <XFlash />
      </Container>
      {(items || header) && (
        <Container className="my-2 my-md-4">
          {items && <XBreadcrumb items={items} />}
          {header && <h1 className={center ? "text-center" : ""}>{header}</h1>}
        </Container>
      )}
      {children}
      <FooterBlock />
    </RootLayout>
  );
}
