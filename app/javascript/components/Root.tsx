import * as Routes from "@/routes.js";
import "dayjs/locale/ru";
import { PrimeReactProvider } from "primereact/api";
import "react-bootstrap";
import type { RootProps } from "@/types";
import i18next from "i18next";
import type { PropsWithChildren } from "react";
import { initReactI18next } from "react-i18next";
import locales from "../locales.json";
import { dayjs } from "@/lib/utils";

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

export default function Root(props: PropsWithChildren) {
  const typedProps = props as RootProps;
  // const { locale, suffix } = usePage<SharedProps>().props;
  const { locale, suffix } = typedProps.initialPage.props;

  // useEffect?
  i18next.changeLanguage(locale);
  dayjs.locale(locale);
  console.log(locale, dayjs.locale());
  Routes.configure({ default_url_options: { suffix } });

  return <PrimeReactProvider>{props.children}</PrimeReactProvider>;
}
