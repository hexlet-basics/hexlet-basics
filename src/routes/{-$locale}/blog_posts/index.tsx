import { Container, SimpleGrid, Stack, Text, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { listBlogPostsOptions } from "@/client/@tanstack/react-query.gen";
import BlogPostBlock from "@/components/BlogPostBlock";
import Breadcrumbs, { CurrentCrumb } from "@/components/Breadcrumbs";
import { seoHead } from "@/lib/seo-head";

// The blog, at its legacy URL, ported from legacy blog_posts/index. The list is
// the first page only: legacy paginated server-side but its pager was commented
// out, so page one is all a visitor could reach. The loader prefetches it so the
// posts are in the server-rendered HTML (ADR-0008).
export const Route = createFileRoute("/{-$locale}/blog_posts/")({
  loader: ({ context }) =>
    context.queryClient.query({ ...listBlogPostsOptions(), staleTime: "static" }),
  // Legacy meta: the header as title, the blog description, the canonical link
  // and the social tags.
  head: ({ match }) => {
    const { i18n } = match.context;
    return seoHead({
      i18n,
      title: i18n.t(($) => $.blog_posts.index.header),
      description: i18n.t(($) => $.blog_posts.index.meta.description),
      canonicalPath: match.pathname,
    });
  },
  component: Index,
});

function Index() {
  const { t } = useTranslation();
  const { data } = useQuery(listBlogPostsOptions());
  const header = t(($) => $.blog_posts.index.header);
  const posts = data?.items ?? [];

  return (
    <Container size="lg" my="xl">
      <Stack mb="xl">
        <Breadcrumbs homeLabel={t(($) => $.blog_posts.show.to_home_title)}>
          <CurrentCrumb>{header}</CurrentCrumb>
        </Breadcrumbs>
        <Title order={1}>{header}</Title>
      </Stack>

      {posts.length === 0 ? (
        <Text c="dimmed">{t(($) => $.blog_posts.index.empty)}</Text>
      ) : (
        <SimpleGrid cols={{ base: 1, xs: 2 }} spacing="md">
          {posts.map((post) => (
            <BlogPostBlock key={post.id} post={post} />
          ))}
        </SimpleGrid>
      )}
    </Container>
  );
}
