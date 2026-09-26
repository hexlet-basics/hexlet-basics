import { Box, Card, Center, Container, Grid, SimpleGrid, Stack, Text, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute, notFound } from "@tanstack/react-router";
import { isAxiosError } from "axios";
import { useTranslation } from "react-i18next";
import { getPublicCourseCategoryOptions } from "@/client/@tanstack/react-query.gen";
import Breadcrumbs, { CurrentCrumb } from "@/components/Breadcrumbs";
import CourseBlock from "@/components/CourseBlock";
import LeadFormBlock from "@/components/LeadFormBlock";
import MarkdownViewer from "@/components/MarkdownViewer";
import { TextLink } from "@/components/RouterLink";

// A course category, at its legacy URL, ported from legacy
// language_categories/show: the category's courses, a consultation form for
// signed-in ru visitors, and the category's questions.
export const Route = createFileRoute("/{-$locale}/language_categories/$slug")({
  loader: async ({ context, params }) => {
    try {
      return await context.queryClient.ensureQueryData(
        getPublicCourseCategoryOptions({ path: { slug: params.slug } }),
      );
    } catch (error) {
      // A missing slug and a category of another language are the same to a
      // visitor: there is no such page.
      if (isAxiosError(error) && error.response?.status === 404) throw notFound();
      throw error;
    }
  },
  // Legacy meta: title and description built around the category header, the
  // page's canonical URL (the API builds the absolute legacy URL) and a summary
  // Twitter card.
  head: ({ loaderData, match }) => {
    if (!loaderData) return {};
    const { t } = match.context.i18n;
    const name = loaderData.category.header ?? "";
    // The ru copy ends in a newline, which legacy's title tag trimmed away.
    const title = t(($) => $.course_categories.show.header, { name }).trim();
    const description = t(($) => $.course_categories.show.meta.description, { name });
    return {
      meta: [
        { title: `CodeBasics | ${title}` },
        { name: "description", content: description },
        { property: "og:title", content: title },
        { property: "og:description", content: description },
        { name: "twitter:card", content: "summary" },
        { name: "twitter:site", content: t(($) => $.links.hexlet_twitter_handle) },
      ],
      links: [{ rel: "canonical", href: loaderData.url }],
    };
  },
  component: Show,
});

function Show() {
  const { t, i18n } = useTranslation();
  const { user } = Route.useRouteContext();
  const { slug } = Route.useParams();
  const { data } = useQuery(getPublicCourseCategoryOptions({ path: { slug } }));
  if (!data) return null;
  const { category, landingPages, qnaItems } = data;

  return (
    <Container size="lg" my="xl">
      <Stack mb="xl">
        <Breadcrumbs homeLabel={t(($) => $.courses.show.to_home_title)}>
          <TextLink to="/{-$locale}/language_categories" size="sm">
            {t(($) => $.course_categories.index.header)}
          </TextLink>
          <CurrentCrumb>{category.header}</CurrentCrumb>
        </Breadcrumbs>
        <Title order={1}>{category.header}</Title>
      </Stack>

      {category.description && (
        <Grid mb="xl">
          <Grid.Col span={{ base: 12, sm: 8 }}>
            <Text size="lg">{category.description}</Text>
          </Grid.Col>
        </Grid>
      )}

      <SimpleGrid cols={{ base: 2, xs: 3, sm: 4 }} mb="xl">
        {landingPages.map((item) => (
          <CourseBlock key={item.id} item={item} />
        ))}
      </SimpleGrid>

      {/* Legacy offered the consultation only to signed-in ru visitors. */}
      {user && i18n.language === "ru" && (
        <Grid align="center" justify="space-between" gap={0}>
          <Grid.Col span={{ base: 12, xs: 7 }}>
            <Center>
              <Text fz={40} mb="xs" fw="bold">
                {t(($) => $.home.index.consultation)}
              </Text>
            </Center>
          </Grid.Col>
          <Grid.Col span={{ base: 12, sm: 5 }}>
            <Card withBorder shadow="sm" p="xl">
              <LeadFormBlock />
            </Card>
          </Grid.Col>
        </Grid>
      )}

      {qnaItems.length > 0 && (
        <Stack py="xl">
          <Title order={2}>{t(($) => $.courses.show.sort_questions)}</Title>
          <SimpleGrid cols={{ base: 1, xs: 2 }}>
            {qnaItems.map((item) => (
              <Box key={item.id}>
                <Text size="lg" fw="bold">
                  {item.question}
                </Text>
                <MarkdownViewer>{item.answer}</MarkdownViewer>
              </Box>
            ))}
          </SimpleGrid>
        </Stack>
      )}
    </Container>
  );
}
