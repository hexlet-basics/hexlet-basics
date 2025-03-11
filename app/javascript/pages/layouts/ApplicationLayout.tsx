import type { PropsWithChildren } from "react";
import { Container } from "react-bootstrap";

import XFlash from "@/components/XFlash.tsx";
import { XBreadcrumb } from "@/components/breadcrumbs.tsx";
import * as Routes from "@/routes.js";
import type { SharedProps } from "@/types";
import type { BreadcrumbItem } from "@/types/index.js";
import { usePage } from "@inertiajs/react";
import RootLayout from "./RootLayout.tsx";
import FooterBlock from "./blocks/FooterBlock.tsx";
import NavbarBlock from "./blocks/NavbarBlock.tsx";
import TgContestBanner from "./tg_contest_banner/TgContestBanner.tsx";

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
  const { locale } = usePage<SharedProps>().props;
  const currentPath = window.location.pathname;

  const isHomePage = currentPath === Routes.root_path();
  const isCoursePage = currentPath.startsWith("/ru/languages/");

  return (
    <RootLayout>
      {locale === "ru" && (isHomePage || isCoursePage) && <TgContestBanner />}
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
