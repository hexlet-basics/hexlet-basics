import type { PropsWithChildren } from "react";
import { Container } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import RootLayout from "./RootLayout.tsx";
import NavbarBlock from "./blocks/NavbarBlock.tsx";

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
    <RootLayout>
      <Container fluid className="mb-2">
        <NavbarBlock />
      </Container>
      {children}
    </RootLayout>
  );
}
