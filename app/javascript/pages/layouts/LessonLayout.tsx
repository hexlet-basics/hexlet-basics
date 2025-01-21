import { type PropsWithChildren, useEffect } from "react";
import { Container } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import * as Routes from "@/routes.js";

import type { SharedProps } from "@/types/index.js";
import { usePage } from "@inertiajs/react";
import i18next from "i18next";
import NavbarBlock from "./NavbarBlock.tsx";

type Props = PropsWithChildren & {};

export default function LessonLayout({ children }: Props) {
  const { t: tLayouts } = useTranslation("layouts");
  const { t: tCommon } = useTranslation("common");
  // const { locale, suffix } = usePage<SharedProps>().props;
  //
  // useEffect(() => {
  //   i18next.changeLanguage(locale);
  //   Routes.configure({ default_url_options: { suffix } });
  // }, [suffix, locale]);

  return (
    <>
      <Container fluid className="mb-2">
        <NavbarBlock />
      </Container>
      {children}
    </>
  );
}
