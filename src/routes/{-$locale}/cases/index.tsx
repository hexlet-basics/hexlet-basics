import { Card, Container, Group, SimpleGrid, Stack, Text, Title } from "@mantine/core";
import { IconArrowRight } from "@tabler/icons-react";
import { createFileRoute, notFound } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { TextLink } from "@/components/RouterLink";

// The cases list, ported from legacy cases/index. Legacy served it for ru only
// (`require_russian_locale`, whose RoutingError is a 404).
export const Route = createFileRoute("/{-$locale}/cases/")({
  beforeLoad: ({ params }) => {
    if (params.locale !== "ru") throw notFound();
  },
  // Legacy meta: title and description only.
  head: ({ match }) => {
    const { t } = match.context.i18n;
    return {
      meta: [
        { title: `CodeBasics | ${t(($) => $.cases.index.title)}` },
        { name: "description", content: t(($) => $.cases.index.meta.description) },
      ],
    };
  },
  component: Index,
});

function Index() {
  const { t } = useTranslation();

  return (
    <Container py="md" mih="100%">
      <SimpleGrid cols={{ base: 1, md: 2, lg: 3 }} py="md">
        <Card bg="gray.0" p="xl" radius="xl" shadow="sm" h="100%">
          <Stack h="100%">
            <Title order={4} mb="xs" fw="bold">
              {t(($) => $.cases.index.for_teachers)}
            </Title>
            <TextLink to="/{-$locale}/cases/for_teachers" c="anchor" td="none" mt="auto">
              <Group gap={6}>
                <Text span>{t(($) => $.cases.index.link)}</Text>
                <IconArrowRight size={16} />
              </Group>
            </TextLink>
          </Stack>
        </Card>
      </SimpleGrid>
    </Container>
  );
}
