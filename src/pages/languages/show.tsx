import { Container, Text, Title } from "@mantine/core";
import { getRouteApi } from "@tanstack/react-router";

const route = getRouteApi("/languages/$slug");

// Stub course page — to be ported from legacy next.
export default function Show() {
  const { slug } = route.useParams();
  return (
    <Container size="lg" my="xl">
      <Title order={1}>{slug}</Title>
      <Text c="dimmed">Course page — coming soon.</Text>
    </Container>
  );
}
