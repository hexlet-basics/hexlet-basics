import { createFileRoute, notFound, Outlet, useMatches } from "@tanstack/react-router";
import ApplicationLayout from "@/components/layout/ApplicationLayout";
import LessonLayout from "@/components/layout/LessonLayout";
import { DEFAULT_LOCALE, isLocale, type Locale } from "@/lib/i18n";

// A route can ask for chrome other than the site layout. Declared per route
// rather than matched on the URL here, so the page that wants a bare shell says
// so itself and this layout keeps knowing nothing about which pages exist.
declare module "@tanstack/react-router" {
  interface StaticDataRouteOption {
    chrome?: "bare";
  }
}

// Optional locale prefix layout: `/` → en (unprefixed), `/ru`, `/es` → prefixed
// (legacy `scope "(:suffix)", suffix: /es|ru/`). Only `ru`/`es` are valid
// prefixes; any other segment (e.g. `/xx`) is a 404, never rendered as a locale.
export const Route = createFileRoute("/{-$locale}")({
  beforeLoad: async ({ params, context }) => {
    const prefix = params.locale;
    if (prefix !== undefined && (prefix === "en" || !isLocale(prefix))) {
      // `en` is the unprefixed default, so `/en/...` is not a valid URL either.
      throw notFound();
    }

    const locale: Locale = prefix ?? DEFAULT_LOCALE;
    await context.i18n.changeLanguage(locale);

    return { locale };
  },
  component: LocaleLayout,
});

function LocaleLayout() {
  const bare = useMatches({
    select: (matches) => matches.some((match) => match.staticData.chrome === "bare"),
  });
  const Layout = bare ? LessonLayout : ApplicationLayout;

  return (
    <Layout>
      <Outlet />
    </Layout>
  );
}
