import { Container, Group, SimpleGrid, Stack, Text, Title } from "@mantine/core";
import { IconArrowRight } from "@tabler/icons-react";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { listPublicCourseCategoriesOptions } from "@/client/@tanstack/react-query.gen";
import Breadcrumbs, { CurrentCrumb } from "@/components/Breadcrumbs";
import { CardLink } from "@/components/RouterLink";
import { seoHead } from "@/lib/seo-head";

// The course categories, at their legacy URL, ported from legacy
// language_categories/index. The loader prefetches the list so the cards are in
// the server-rendered HTML (ADR-0008).
export const Route = createFileRoute("/{-$locale}/language_categories/")({
  loader: ({ context }) =>
    context.queryClient.query({ ...listPublicCourseCategoriesOptions(), staleTime: "static" }),
  // Legacy meta: the header as title, the index description, the canonical link
  // and the social tags.
  head: ({ match }) => {
    const { i18n } = match.context;
    return seoHead({
      i18n,
      title: i18n.t(($) => $.course_categories.index.header),
      description: i18n.t(($) => $.course_categories.index.meta.description),
      canonicalPath: match.pathname,
    });
  },
  component: Index,
});

function Index() {
  const { t } = useTranslation();
  const { data } = useQuery(listPublicCourseCategoriesOptions());
  const header = t(($) => $.course_categories.index.header);
  const categories = data ?? [];

  return (
    <Container size="lg" my="xl">
      <Stack mb="xl">
        <Breadcrumbs homeLabel={t(($) => $.courses.show.to_home_title)}>
          <CurrentCrumb>{header}</CurrentCrumb>
        </Breadcrumbs>
        <Title order={1}>{header}</Title>
      </Stack>

      <SimpleGrid py="md" cols={{ base: 1, md: 2, lg: 3 }}>
        {categories.map((category) => (
          <CardLink
            key={category.id}
            to="/{-$locale}/language_categories/$slug"
            params={{ slug: category.slug ?? "" }}
            padding="xl"
            radius="md"
            withBorder
            h="100%"
            td="none"
          >
            <Stack h="100%">
              <Text size="xl" fw={700}>
                {category.header}
              </Text>
              <Group mt="auto" c="blue">
                <Text>{t(($) => $.course_categories.index.link)}</Text>
                <IconArrowRight size={16} />
              </Group>
            </Stack>
          </CardLink>
        ))}
      </SimpleGrid>
    </Container>
  );
}
