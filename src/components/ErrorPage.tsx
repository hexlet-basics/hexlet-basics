import { Container, Stack, Text, Title } from "@mantine/core";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";

// The HTTP errors that have their own wording; any other code reads the
// generic one.
type ErrorCode = 403 | 404 | 500;

// Typed selectors cannot take a runtime key, so each code names its strings.
// The generic description falls back to the generic *header*, as legacy
// errors#show did (`default: codes.other.header`); `codes.other.description`
// was never shown.
function errorText(t: TFunction, code: number): { header: string; description: string } {
  switch (code as ErrorCode) {
    case 403:
      return {
        header: t(($) => $.errors.show.codes["403"].header),
        description: t(($) => $.errors.show.codes["403"].description),
      };
    case 404:
      return {
        header: t(($) => $.errors.show.codes["404"].header),
        description: t(($) => $.errors.show.codes["404"].description),
      };
    case 500:
      return {
        header: t(($) => $.errors.show.codes["500"].header),
        description: t(($) => $.errors.show.codes["500"].description),
      };
    default:
      return {
        header: t(($) => $.errors.show.codes.other.header),
        description: t(($) => $.errors.show.codes.other.header),
      };
  }
}

// The error page, ported from legacy errors/show: the code, a header and what
// to do about it. The status line itself is the router's: a not-found match
// answers 404 and a failed load 500.
export default function ErrorPage({ code }: { code: number }) {
  const { t } = useTranslation();
  const { header, description } = errorText(t, code);

  return (
    <Container ta="center" my="xl" py="xl">
      <Stack align="center" gap="md">
        <Text fw={700}>{code}</Text>
        <Title order={1}>{header}</Title>
        <Text size="lg">{description}</Text>
      </Stack>
    </Container>
  );
}

// Legacy errors/show had its own 403 wording. No ported page answers 403:
// legacy's web pages redirected a visitor without access (the admin guards
// still do), and its API's `head :forbidden` carried no page. The page is here
// for the first surface that refuses rather than redirects.
export function ForbiddenPage() {
  return <ErrorPage code={403} />;
}

export function NotFoundPage() {
  return <ErrorPage code={404} />;
}

export function ServerErrorPage() {
  return <ErrorPage code={500} />;
}
