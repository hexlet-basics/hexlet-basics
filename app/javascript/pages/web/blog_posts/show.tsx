import { Head } from '@inertiajs/react';
import {
  Alert,
  AspectRatio,
  Box,
  Center,
  Container,
  Group,
  Image,
  SimpleGrid,
  Skeleton,
  Stack,
  Text,
  Title,
} from '@mantine/core';
import axios from 'axios';
import dayjs from 'dayjs';
import i18next from 'i18next';
import {
  Clock7,
  MessageCircleMore,
  MoveRight,
  ThumbsUp,
  User,
} from 'lucide-react';
import { type PropsWithChildren, useState } from 'react';
import { useTranslation } from 'react-i18next';
import type { Article, WithContext } from 'schema-dts';
import AppAnchor from '@/components/AppAnchor';
import BlogPostBlock from '@/components/BlogPostBlock';
import MarkdownViewer from '@/components/MarkdownViewer.tsx';
import CoursesList from '@/components/ProgramsList';
import useInfiniteItems from '@/hooks/useInfiniteItems';
import ApplicationLayout from '@/pages/layouts/ApplicationLayout';
import * as Routes from '@/routes.js';
import type {
  BlogPost,
  BreadcrumbItem,
  LanguageLandingPageForLists,
} from '@/types';

type Props = PropsWithChildren & {
  blogPost: BlogPost;
  recommendedBlogPosts: BlogPost[];
  relatedLandingPages: LanguageLandingPageForLists[];
};

export default function Show({
  blogPost,
  recommendedBlogPosts,
  relatedLandingPages,
}: Props) {
  const { t } = useTranslation();
  const { t: tCommon } = useTranslation('common');

  const postUrl = Routes.blog_post_url(blogPost.slug!);
  const items: BreadcrumbItem[] = [
    {
      name: t('blog_posts.index.header'),
      url: Routes.blog_posts_url(),
    },
    {
      name: blogPost.name!,
      url: postUrl,
    },
  ];

  const article: WithContext<Article> = {
    '@context': 'https://schema.org',
    '@type': 'Article',
    author: blogPost.creator.name!,
    name: blogPost.name!,
    datePublished: blogPost.created_at,
    headline: blogPost.description!,
    image: blogPost.cover_main_variant!,
  };

  const loadNext = async (lastPostId: number): Promise<BlogPost> => {
    const res = await axios.get<BlogPost>(
      Routes.next_api_blog_post_path(lastPostId),
    );
    return res.data;
  };

  const {
    items: posts,
    setContainerRef,
    markerRef,
  } = useInfiniteItems<BlogPost>(blogPost, loadNext);

  const components = {
    '::courses': () => <CoursesList landingPages={relatedLandingPages} />,
  };

  return (
    <>
      <Head>
        <script type="application/ld+json">{JSON.stringify(article)}</script>
      </Head>

      <ApplicationLayout items={items} center header={blogPost.name!}>
        <Container size="sm">
          {posts.map((post, index) => (
            <Stack
              ref={(ref) => setContainerRef(ref, post)}
              key={post.id}
              mb="xl"
            >
              {index !== 0 && (
                <Title order={2} mt="xl" mb="sm">
                  {post.name}
                </Title>
              )}
              <AspectRatio ratio={2 / 1}>
                <Image
                  fit="cover"
                  w="100%"
                  fallbackSrc="data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///ywAAAAAAQABAAACAUwAOw=="
                  h="100%"
                  radius="md"
                  src={post.cover_main_variant!}
                />
              </AspectRatio>
              <MarkdownViewer components={components} allowHtml>
                {post.body || ''}
              </MarkdownViewer>

              {index === 0 && (
                <Box>
                  <Group mb="lg">
                    <Group fw="bold" me="auto">
                      <User size={18} />
                      {post.creator.name}
                      {dayjs().to(post.created_at)}
                    </Group>
                    <Group gap={0} me="lg">
                      <AppAnchor href={postUrl} me="xs" display="flex">
                        <ThumbsUp size={18} />
                      </AppAnchor>
                      {post.likes_count}
                    </Group>
                    <Center>
                      <Center me="xs">
                        <Clock7 size={18} />
                      </Center>
                      {tCommon('time.minutes', { count: 5 })}
                    </Center>
                  </Group>

                  {i18next.language === 'ru' && (
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
                        <AppAnchor
                          href="https://t.me/HexletLearningBot"
                          className="after:absolute after:inset-0"
                          external
                        >
                          <Text component="span" mr="sm">
                            {t('blog_posts.show.link')}
                          </Text>
                        </AppAnchor>
                        <MoveRight />
                      </Group>
                    </Alert>
                  )}

                  <SimpleGrid cols={{ base: 1, xs: 2 }}>
                    {recommendedBlogPosts.map((post) => (
                      <BlogPostBlock key={post.id} post={post} />
                    ))}
                  </SimpleGrid>
                </Box>
              )}
            </Stack>
          ))}
          <div ref={markerRef}></div>
        </Container>
      </ApplicationLayout>
    </>
  );
}
