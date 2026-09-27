import {
  ActionIcon,
  Alert,
  Anchor,
  AspectRatio,
  Box,
  Center,
  Container,
  Group,
  Image,
  SimpleGrid,
  Stack,
  Text,
  Title,
  Typography,
} from "@mantine/core";
import { useIntersection } from "@mantine/hooks";
import { notifications } from "@mantine/notifications";
import {
  IconArrowRight,
  IconClockHour7,
  IconMessageCircle,
  IconThumbUp,
  IconUser,
} from "@tabler/icons-react";
import { useInfiniteQuery, useMutation, useQuery } from "@tanstack/react-query";
import { createFileRoute, notFound, useRouter } from "@tanstack/react-router";
import { isAxiosError } from "axios";
import { useEffect, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { getBlogPostOptions, likeBlogPostMutation } from "@/client/@tanstack/react-query.gen";
import { getNextBlogPost } from "@/client/sdk.gen";
import type { BlogPost, CourseCatalogItem } from "@/client/types.gen";
import BlogPostBlock from "@/components/BlogPostBlock";
import Breadcrumbs, { CurrentCrumb } from "@/components/Breadcrumbs";
import CourseBlock from "@/components/CourseBlock";
import { ActionIconLink, TextLink } from "@/components/RouterLink";
import { seoHead } from "@/lib/seo-head";
import { timeAgo } from "@/lib/time-ago";

// A blog post, at its legacy URL, ported from legacy blog_posts/show. Reading on
// scrolls into the next older post (legacy useInfiniteItems), and the address
// bar follows whichever post is being read.
export const Route = createFileRoute("/{-$locale}/blog_posts/$slug")({
  loader: async ({ context, params }) => {
    try {
      return await context.queryClient.query({
        ...getBlogPostOptions({ path: { slug: params.slug } }),
        staleTime: "static",
      });
    } catch (error) {
      // A draft, a post in another language and a missing slug are all the
      // same to a visitor: there is no such page.
      if (isAxiosError(error) && error.response?.status === 404) throw notFound();
      throw error;
    }
  },
  // Legacy meta: the post's name and description, its canonical URL, the cover
  // as og:image, and a schema.org Article carrying the like count.
  head: ({ loaderData, match }) => {
    if (!loaderData) return {};
    const { post } = loaderData;
    const article = {
      "@context": "https://schema.org",
      "@type": "Article",
      author: post.creator.name,
      name: post.name,
      datePublished: post.createdAt,
      headline: post.description,
      image: post.coverMainVariant,
      interactionStatistic: [
        {
          "@type": "InteractionCounter",
          interactionType: { "@type": "LikeAction" },
          userInteractionCount: post.likesCount,
        },
      ],
    };
    return {
      ...seoHead({
        i18n: match.context.i18n,
        title: post.name ?? "",
        description: post.description ?? "",
        image: post.coverMainVariant,
        canonicalPath: match.pathname,
      }),
      scripts: [
        {
          type: "application/ld+json",
          // `<` escaped so a post name cannot close the script element.
          children: JSON.stringify(article).replace(/</g, "\\u003c"),
        },
      ],
    };
  },
  component: Show,
});

// The transparent pixel legacy fell back to for a post without a cover.
const NO_COVER = "data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///ywAAAAAAQABAAACAUwAOw==";

function Show() {
  const { t } = useTranslation();
  const { slug } = Route.useParams();
  const { data } = useQuery(getBlogPostOptions({ path: { slug } }));
  if (!data) return null;
  const { post, recommendedPosts, relatedLandingPages } = data;

  return (
    <>
      <Container size="lg" my="xl">
        <Stack>
          <Breadcrumbs homeLabel={t(($) => $.blog_posts.show.to_home_title)}>
            <TextLink to="/{-$locale}/blog_posts" size="sm">
              {t(($) => $.blog_posts.index.header)}
            </TextLink>
            <CurrentCrumb>{post.name}</CurrentCrumb>
          </Breadcrumbs>
          <Title order={1} ta="center">
            {post.name}
          </Title>
        </Stack>
      </Container>

      {/* Keyed by post so reading on from a fresh slug starts a fresh chain. */}
      <PostChain
        key={post.id}
        first={post}
        recommendedPosts={recommendedPosts}
        relatedLandingPages={relatedLandingPages}
      />
    </>
  );
}

// The post and the older posts appended under it as the reader reaches the
// bottom. As in legacy, every appended post shows the first post's related
// courses and recommendations (they are page-level data).
function PostChain({
  first,
  recommendedPosts,
  relatedLandingPages,
}: {
  first: BlogPost;
  recommendedPosts: BlogPost[];
  relatedLandingPages: CourseCatalogItem[];
}) {
  const next = useInfiniteQuery({
    queryKey: ["blogPostChain", first.id],
    initialPageParam: first.id,
    queryFn: async ({ pageParam, signal }) => {
      const { data } = await getNextBlogPost({
        path: { id: pageParam },
        signal,
        throwOnError: true,
      });
      return data;
    },
    getNextPageParam: (last) => last.id,
    // Pages load only when the reader nears the end, never on mount.
    enabled: false,
    retry: false,
  });

  // Legacy watched a bottom marker with a 600px top margin. A failed fetch — the
  // 404 after the oldest post — ends the chain.
  const { ref: markerRef, entry } = useIntersection({ rootMargin: "600px 0px 0px 0px" });
  const nearEnd = entry?.isIntersecting ?? false;
  const { fetchNextPage, isFetching, isError } = next;
  useEffect(() => {
    if (nearEnd && !isFetching && !isError) void fetchNextPage();
  }, [nearEnd, isFetching, isError, fetchNextPage]);

  const posts = [first, ...(next.data?.pages ?? [])];
  const track = useSyncedUrl();

  return (
    <Container size="sm">
      {posts.map((post, index) => (
        <Stack key={post.id} mb="xl" ref={(element) => track(post, element)}>
          {index !== 0 && (
            <Title order={1} mt="xl" mb="sm">
              {post.name}
            </Title>
          )}
          <PostBody
            post={post}
            recommendedPosts={recommendedPosts}
            relatedLandingPages={relatedLandingPages}
          />
        </Stack>
      ))}
      <div ref={markerRef} />
    </Container>
  );
}

// Keeps the address bar on the post crossing a line a quarter down the
// viewport, ported from legacy useInfiniteItems (the issue #587 geometry: live
// bounding boxes on scroll, symmetric up and down). The URL is replaced behind
// the router's back on purpose — it records where the reader is, it is not a
// navigation, so no loader runs and the chain is not torn down.
function useSyncedUrl() {
  // Insertion order is chain order: posts only ever append.
  const boxes = useRef(new Map<number, { post: BlogPost; element: HTMLElement }>());
  // The chain is in the page's locale, so each post's path is built under the
  // page's own locale segment, by the router that owns the prefix rule.
  const router = useRouter();
  const { locale } = Route.useParams();

  useEffect(() => {
    let frame = 0;

    const syncUrl = () => {
      frame = 0;
      const referenceLine = window.innerHeight * 0.25;

      let active: { post: BlogPost; top: number } | null = null;
      for (const { post, element } of boxes.current.values()) {
        const { top, bottom } = element.getBoundingClientRect();
        if (top <= referenceLine && bottom > referenceLine) {
          active = { post, top };
          break;
        }
        // Fallback for gaps between posts: the nearest post above the line.
        if (top <= referenceLine && (!active || top > active.top)) {
          active = { post, top };
        }
      }

      if (!active) return;
      const { pathname } = router.buildLocation({
        to: "/{-$locale}/blog_posts/$slug",
        params: { locale, slug: active.post.slug ?? "" },
      });
      if (window.location.pathname !== pathname) {
        window.history.replaceState(window.history.state, "", pathname);
      }
    };

    const onScroll = () => {
      if (frame) return;
      frame = requestAnimationFrame(syncUrl);
    };

    window.addEventListener("scroll", onScroll, { passive: true });
    return () => {
      window.removeEventListener("scroll", onScroll);
      if (frame) cancelAnimationFrame(frame);
    };
  }, [router, locale]);

  return (post: BlogPost, element: HTMLElement | null) => {
    if (element) boxes.current.set(post.id, { post, element });
    else boxes.current.delete(post.id);
  };
}

function PostBody({
  post,
  recommendedPosts,
  relatedLandingPages,
}: {
  post: BlogPost;
  recommendedPosts: BlogPost[];
  relatedLandingPages: CourseCatalogItem[];
}) {
  const { t, i18n } = useTranslation();

  return (
    <>
      <AspectRatio ratio={2 / 1}>
        <Image
          fit="cover"
          w="100%"
          h="100%"
          radius="md"
          fallbackSrc={NO_COVER}
          src={post.coverMainVariant}
          alt={post.name ?? ""}
        />
      </AspectRatio>

      <Typography>
        {/* Trusted editor HTML, stored and served as written (see blog_posts.go). */}
        <Box
          className="blog-post-content"
          dangerouslySetInnerHTML={{ __html: post.richBodyHtml }}
        />
      </Typography>

      {relatedLandingPages.length > 0 && (
        <SimpleGrid cols={{ base: 1, xs: 2 }} spacing="xl">
          {relatedLandingPages.map((item) => (
            <CourseBlock key={item.id} lazy item={item} />
          ))}
        </SimpleGrid>
      )}

      <Box>
        <Group mb="lg">
          <Group me="auto">
            <IconUser size={18} />
            <Text fw="bold">{post.creator.name}</Text>
            <Text fw="bold">{timeAgo(post.createdAt, i18n.language)}</Text>
          </Group>
          <LikeButton post={post} />
          <Center>
            <Center me="xs">
              <IconClockHour7 size={18} />
            </Center>
            {t(($) => $.common.time.minutes, { count: 5 })}
          </Center>
        </Group>

        {i18n.language === "ru" && (
          <Alert
            radius="lg"
            p="xl"
            mb="xl"
            title={t(($) => $.blog_posts.show.join_community)}
            icon={<IconMessageCircle />}
          >
            <Text fz="lg" lh="sm" mb="md">
              {t(($) => $.blog_posts.show.discuss)}
            </Text>
            <Anchor
              href={t(($) => $.common.community_url)}
              target="_blank"
              rel="noreferrer"
              underline="never"
            >
              <Group gap={0}>
                <Text component="span" mr="sm">
                  {t(($) => $.blog_posts.show.link)}
                </Text>
                <IconArrowRight />
              </Group>
            </Anchor>
          </Alert>
        )}

        <SimpleGrid cols={{ base: 1, xs: 2 }}>
          {recommendedPosts.map((recommended) => (
            <BlogPostBlock key={recommended.id} post={recommended} />
          ))}
        </SimpleGrid>
      </Box>
    </>
  );
}

// Legacy posted the like and redirected back with a flash: "thanks" for a new
// like, "already counted" for a repeat one. The API answers both with the post,
// so a count that did not move is the repeat. A visitor is sent to sign in, as
// the legacy auth guard did.
function LikeButton({ post }: { post: BlogPost }) {
  const { t } = useTranslation();
  const { user } = Route.useRouteContext();
  const [likesCount, setLikesCount] = useState(post.likesCount);
  const like = useMutation({
    ...likeBlogPostMutation(),
    onSuccess: (liked) => {
      notifications.show({
        color: "green",
        message:
          liked.likesCount > likesCount
            ? t(($) => $.flash.blog_posts.likes.create.success)
            : t(($) => $.flash.blog_posts.likes.create.notice),
      });
      setLikesCount(liked.likesCount);
    },
    onError: () => notifications.show({ color: "red", message: t(($) => $.common.errors.network) }),
  });

  return (
    <Group gap={0} me="lg">
      {user ? (
        <ActionIcon
          variant="subtle"
          me="xs"
          aria-label="like"
          loading={like.isPending}
          onClick={() => like.mutate({ path: { id: post.id } })}
        >
          <IconThumbUp size={18} />
        </ActionIcon>
      ) : (
        <ActionIconLink to="/{-$locale}/session/new" variant="subtle" me="xs" aria-label="like">
          <IconThumbUp size={18} />
        </ActionIconLink>
      )}
      {likesCount}
    </Group>
  );
}
