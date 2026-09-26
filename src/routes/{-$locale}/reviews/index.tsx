import {
  Alert,
  Anchor,
  Card,
  Container,
  Group,
  SimpleGrid,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { IconUserCircle } from "@tabler/icons-react";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { Trans, useTranslation } from "react-i18next";
import { z } from "zod";
import { listPublicReviewsOptions } from "@/client/@tanstack/react-query.gen";
import Breadcrumbs, { CurrentCrumb } from "@/components/Breadcrumbs";
import Pager from "@/components/Pager";
import { longDate } from "@/lib/time-ago";

// Where a ru reader is sent to leave a review (legacy reviews/index).
const ADD_REVIEW_URL = "https://taplink.cc/codebasics_reviews";

// Student reviews, at their legacy URL, ported from legacy reviews/index. The
// loader prefetches the requested page so the reviews are in the
// server-rendered HTML (ADR-0008); a malformed `page` falls back to page one.
export const Route = createFileRoute("/{-$locale}/reviews/")({
  validateSearch: z.object({
    page: z.number().int().min(1).optional().catch(undefined),
  }),
  loaderDeps: ({ search }) => ({ page: search.page }),
  loader: ({ context, deps }) =>
    context.queryClient.ensureQueryData(listPublicReviewsOptions({ query: { page: deps.page } })),
  // Legacy meta: the header as title (behind the site name, as meta-tags
  // rendered it) and og:title (legacy asked for a `.title` key that never
  // existed), the reviews description, and a summary Twitter card. The
  // canonical link is left out until the frontend knows its own public origin.
  head: ({ match }) => {
    const { t } = match.context.i18n;
    const title = t(($) => $.reviews.index.header);
    const description = t(($) => $.reviews.index.meta.description);
    return {
      meta: [
        { title: `CodeBasics | ${title}` },
        { name: "description", content: description },
        { property: "og:title", content: title },
        { property: "og:description", content: description },
        { name: "twitter:card", content: "summary" },
        { name: "twitter:site", content: t(($) => $.links.hexlet_twitter_handle) },
      ],
    };
  },
  component: Index,
});

function Index() {
  const { t, i18n } = useTranslation();
  const { page } = Route.useSearch();
  const { data } = useQuery(listPublicReviewsOptions({ query: { page } }));
  const header = t(($) => $.reviews.index.header);
  const reviews = data?.items ?? [];

  return (
    <Container size="lg" my="xl">
      <Stack mb="xl">
        <Breadcrumbs homeLabel={t(($) => $.blog_posts.show.to_home_title)}>
          <CurrentCrumb>{header}</CurrentCrumb>
        </Breadcrumbs>
        <Title order={1}>{header}</Title>
      </Stack>

      {i18n.language === "ru" && (
        <Alert mb="xl">
          <Trans
            t={t}
            i18nKey={($) => $.reviews.index.add_review}
            components={{
              a: <Anchor href={ADD_REVIEW_URL} target="_blank" rel="noopener noreferrer" />,
            }}
          />
        </Alert>
      )}

      <SimpleGrid cols={{ base: 1, xs: 2 }}>
        {reviews.map((review) => (
          <Card key={review.id} shadow="sm" padding="lg" radius="md" withBorder>
            <Group mb="md">
              <IconUserCircle size={20} />
              <Text fw={500} size="lg">
                {review.fullName}
              </Text>
            </Group>
            <Text mb="sm">{review.body}</Text>
            <Group justify="space-between" c="dimmed" mt="auto">
              <Text size="sm">{longDate(review.createdAt, i18n.language)}</Text>
            </Group>
          </Card>
        ))}
      </SimpleGrid>

      {data && <Pager total={data.total} page={data.page} perPage={data.perPage} />}
    </Container>
  );
}
