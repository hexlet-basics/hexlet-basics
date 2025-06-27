import type { PropsWithChildren } from "react";
import { Anchor, Box, Container, Grid, Image, Stack, Title, Text, Card, Group, Center, Alert, SimpleGrid, TypographyStylesProvider } from '@mantine/core';

import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type { BreadcrumbItem } from "@/types";
import type { BlogPost } from "@/types/serializers";
import i18next from "i18next";

import BlogPostBlock from "@/components/BlogPostBlock";
import MarkdownViewer from "@/components/MarkdownViewer.tsx";
import * as Routes from "@/routes.js";
import { Head } from "@inertiajs/react";
import { useTranslation } from "react-i18next";
import type { Article, WithContext } from "schema-dts";
import dayjs from "dayjs";
import { Clock7, MessageCircleMore, MoveRight, ThumbsUp, User } from "lucide-react";
import AppAnchor from "@/components/AppAnchor";

type Props = PropsWithChildren & {
  blogPost: BlogPost;
  recommendedBlogPosts: BlogPost[];
};

export default function Show({ blogPost, recommendedBlogPosts }: Props) {
  const { t } = useTranslation();
  const { t: tCommon } = useTranslation("common");

  const postUrl = Routes.blog_post_url(blogPost.slug!);
  const items: BreadcrumbItem[] = [
    {
      name: t("blog_posts.index.header"),
      url: Routes.blog_posts_url(),
    },
    {
      name: blogPost.name!,
      url: postUrl,
    },
  ];

  const article: WithContext<Article> = {
    "@context": "https://schema.org",
    "@type": "Article",
    author: blogPost.creator.name!,
    name: blogPost.name!,
    datePublished: blogPost.created_at,
    headline: blogPost.description!,
    image: blogPost.cover_main_variant!,
  };

  return (
    <>
      <Head>
        <script type="application/ld+json">{JSON.stringify(article)}</script>
      </Head>

      <ApplicationLayout items={items} center header={blogPost.name!}>
        <Container size="sm">
          <Stack>
            <Image
              className="img-fluid"
              fetchPriority="high"
              radius="md"
              src={blogPost.cover_main_variant!}
              mb="xl"
            />
            <MarkdownViewer allowHtml>{blogPost.body || ""}</MarkdownViewer>

            <Group mb="lg">
              <Group fw="bold" me="auto">
                <User size={18} />
                {blogPost.creator.name}
                {dayjs().to(blogPost.created_at)}
              </Group>
              <Group gap={0} me="lg">
                <AppAnchor href={postUrl} me="xs" display="flex">
                  <ThumbsUp size={18} />
                </AppAnchor>
                {blogPost.likes_count}
              </Group>
              <Center>
                <Center me="xs">
                  <Clock7 size={18} />
                </Center>
                {tCommon("time.minutes", { count: 5 })}
              </Center>
            </Group>

            {i18next.language === "ru" && (
              <Alert
                radius="lg"
                p="xl"
                mb="xl"
                pos="relative"
                title={t('blog_posts.show.join_community')}
                icon={<MessageCircleMore />}
              >
                <Text fz="lg" lh="sm" mb="md">
                  {t('blog_posts.show.discuss')}
                </Text>
                <Group gap={0}>
                  <Text component="span" mr="sm">
                    {t('blog_posts.show.link')}
                  </Text>
                  <MoveRight />
                </Group>



                <AppAnchor
                  pos="absolute"
                  inset={0}
                  href="https://t.me/HexletLearningBot"
                  external
                />

              </Alert>
            )}
          </Stack>


          <SimpleGrid cols={{ base: 1, xs: 2 }}>
            {recommendedBlogPosts.map((post) => (
              <BlogPostBlock key={post.id} post={post} />
            ))}
          </SimpleGrid>

        </Container>
      </ApplicationLayout>
    </>
  );
}
