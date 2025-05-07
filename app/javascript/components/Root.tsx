import * as Sentry from "@sentry/react";
import { SnackbarProvider } from "notistack";
import { PrimeReactProvider } from "primereact/api";
import "react-bootstrap";
import { type PropsWithChildren } from "react";
import { Container } from "react-bootstrap";
import { useTranslation } from "react-i18next";
import { Locale } from "@/types";
import configure from "@/lib/configure";

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

type Props = PropsWithChildren & {
  // locale: Locale;
  // suffix: string | null
}

function Root(props: Props) {

  // configure(props.locale, props.suffix)

  return (
    <PrimeReactProvider>
      <SnackbarProvider>
        <Sentry.ErrorBoundary fallback={FallbackComponent} showDialog>
          {props.children}
        </Sentry.ErrorBoundary>
      </SnackbarProvider>
    </PrimeReactProvider>
  );
}

export default Sentry.withProfiler(Root);
