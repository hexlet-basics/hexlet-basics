import { Container, Text, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";

// Stub course page — to be ported from legacy next.
export const Route = createFileRoute("/languages/$slug")({
  component: Show,
});

function Show() {
  const { slug } = Route.useParams();
  return (
    <Container size="lg" my="xl">
      <Title order={1}>{slug}</Title>
      <Text c="dimmed">Course page — coming soon.</Text>
    </Container>
  );
}
