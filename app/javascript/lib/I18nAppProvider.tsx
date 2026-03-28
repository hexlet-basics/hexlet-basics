import type { PageProps } from "@inertiajs/core";
import type { ReactElement, ReactNode } from "react";
import { I18nextProvider, useSSR } from "react-i18next";
import {
  buildInitialI18nStore,
  createServerI18n,
  getClientI18n,
} from "@/lib/i18n";
import type { Locale } from "@/types";

type InertiaAppElementProps = {
  initialPage: {
    props: PageProps;
  };
};

type Props = {
  app: ReactElement;
  ssr: boolean;
};

function I18nHydrator({ app, locale }: { app: ReactNode; locale: Locale }) {
  useSSR(buildInitialI18nStore(locale), locale);

  return app;
}

export default function I18nAppProvider({ app, ssr }: Props) {
  const pageProps = (app.props as InertiaAppElementProps).initialPage.props;
  const locale = pageProps.locale ?? "ru";

  const i18n = ssr ? createServerI18n(locale) : getClientI18n(locale);

  return (
    <I18nextProvider i18n={i18n}>
      {ssr ? app : <I18nHydrator app={app} locale={locale} />}
    </I18nextProvider>
  );
}
