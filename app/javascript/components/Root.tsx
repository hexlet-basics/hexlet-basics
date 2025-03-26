import * as Routes from "@/routes.js";
import * as Sentry from "@sentry/react";
import "dayjs/locale/ru";
import { PrimeReactProvider } from "primereact/api";
import "react-bootstrap";
import { dayjs } from "@/lib/utils";
import type { RootProps } from "@/types";
import i18next from "i18next";
import type { PropsWithChildren } from "react";
import { Container } from "react-bootstrap";
import { initReactI18next, useTranslation } from "react-i18next";
import locales from "../locales.json";

const resources = locales;
const defaultNS = "web";

i18next.use(initReactI18next);
i18next.init({
  resources,
  defaultNS,
  ns: Object.keys(resources.ru),
  // lng: locale,
  interpolation: {
    prefix: "%{",
    suffix: "}",
    escapeValue: false, // react already safes from xss => https://www.i18next.com/translation-function/interpolation#unescape
  },
});

function FallbackComponent() {
  const { t: tLayouts } = useTranslation("layouts");
  return (
    <Container className="d-flex h-100">
      <div className="m-auto">
        <h1>{tLayouts("web.root.fallback.header")}</h1>
        <div>{tLayouts("web.root.fallback.description")}</div>
      </div>
    </Container>
  );
}

function Root(props: PropsWithChildren) {
  const typedProps = props as RootProps;
  // const { locale, suffix } = usePage<SharedProps>().props;
  const { locale, suffix } = typedProps.initialPage.props;

  // useEffect?
  i18next.changeLanguage(locale);
  dayjs.locale(locale);
  console.log(locale, dayjs.locale());
  Routes.configure({ default_url_options: { suffix } });

  return (
    <PrimeReactProvider>
      <Sentry.ErrorBoundary fallback={FallbackComponent} showDialog>
        {props.children}
      </Sentry.ErrorBoundary>
    </PrimeReactProvider>
  );
}

export default Sentry.withProfiler(Root);
