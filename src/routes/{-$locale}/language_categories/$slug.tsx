import { Container, Grid, SimpleGrid, Stack, Text, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { getPublicCourseCategoryOptions } from "@/client/@tanstack/react-query.gen";
import Breadcrumbs, { CurrentCrumb } from "@/components/Breadcrumbs";
import CourseBlock from "@/components/CourseBlock";
import QnaBlock from "@/components/QnaBlock";
import { AnchorLink } from "@/components/RouterLink";
import { seoHead } from "@/lib/seo-head";
import { rethrowMissingAsNotFound } from "@/lib/api-error";

// A course category, at its legacy URL, ported from legacy
// language_categories/show: the category's courses, a consultation form for
// signed-in ru visitors, and the category's questions.
export const Route = createFileRoute("/{-$locale}/language_categories/$slug")({
  loader: async ({ context, params }) => {
    try {
      return await context.queryClient.query({
        ...getPublicCourseCategoryOptions({ path: { slug: params.slug } }),
        staleTime: "static",
      });
    } catch (error) {
      // A missing slug and a category of another language are the same to a
      // visitor: there is no such page.
      return rethrowMissingAsNotFound(error);
    }
  },

  // Legacy meta: title and description built around the category header, the
  // canonical link and the social tags.
  head: ({ loaderData, match }) => {
    if (!loaderData) {
      return {};
    }
    const { i18n } = match.context;
    const name = loaderData.category.header ?? "";
    return seoHead({
      canonicalPath: match.pathname,
      description: i18n.t(($) => $.course_categories.show.meta.description, { name }),
      i18n,
      // The ru copy ends in a newline, which legacy's title tag trimmed away.
      title: i18n.t(($) => $.course_categories.show.header, { name }).trim(),
    });
  },

  component: Show,
});

function Show() {
  const { t } = useTranslation();
  const { slug } = Route.useParams();
  const { data } = useQuery(getPublicCourseCategoryOptions({ path: { slug } }));
  if (!data) {
    return null;
  }
  const { category, landingPages, qnaItems } = data;

  return (
    <Container size="lg" my="xl">
      <Stack mb="xl">
        <Breadcrumbs homeLabel={t(($) => $.courses.show.to_home_title)}>
          <AnchorLink to="/{-$locale}/language_categories" size="sm">
            {t(($) => $.course_categories.index.header)}
          </AnchorLink>
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

      <SimpleGrid cols={{ base: 2, sm: 4, xs: 3 }} mb="xl">
        {landingPages.map((item) => (
          <CourseBlock key={item.id} item={item} />
        ))}
      </SimpleGrid>

      <QnaBlock items={qnaItems} />
    </Container>
  );
}
