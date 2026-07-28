import { createFileRoute, notFound, Outlet } from "@tanstack/react-router";
import ApplicationLayout from "@/components/layout/ApplicationLayout";
import { DEFAULT_LOCALE, isLocale, type Locale } from "@/lib/i18n";

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
  component: () => (
    <ApplicationLayout>
      <Outlet />
    </ApplicationLayout>
  ),
});
