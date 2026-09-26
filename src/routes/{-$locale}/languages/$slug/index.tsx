import { Alert, Container, Text, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { z } from "zod";

// Stub course page — to be ported from legacy next.
//
// `unfinished` is set by the completion page when it bounces a learner whose
// Enrollment is not finished: legacy's flash warning, carried in the URL so it
// survives a server-side redirect.
export const Route = createFileRoute("/{-$locale}/languages/$slug/")({
  validateSearch: z.object({ unfinished: z.boolean().optional() }),
  component: Show,
});

function Show() {
  const { t } = useTranslation();
  const { slug } = Route.useParams();
  const { unfinished } = Route.useSearch();
  return (
    <Container size="lg" my="xl">
      {unfinished && (
        <Alert color="yellow" mb="md">
          {t(($) => $.flash.courses.success.warning)}
        </Alert>
      )}
      <Title order={1}>{slug}</Title>
      <Text c="dimmed">Course page — coming soon.</Text>
    </Container>
  );
}
