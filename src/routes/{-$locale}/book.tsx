import {
  Box,
  Button,
  Container,
  Divider,
  Grid,
  Image,
  Paper,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { IconCompass, IconFileText, IconList, IconUsers } from "@tabler/icons-react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createFileRoute, notFound } from "@tanstack/react-router";
import { truncate } from "es-toolkit/compat";
import { useTranslation } from "react-i18next";
import bookCoverImg from "@/assets/profession-developer-book-cover.webp";
import {
  createBookRequestMutation,
  getBookOptions,
  getBookQueryKey,
} from "@/client/@tanstack/react-query.gen";
import type { DownloadBookData } from "@/client/types.gen";
import LeadFormBlock from "@/components/LeadFormBlock";
import { ButtonLink } from "@/components/RouterLink";
import bookToc from "@/lib/book";
import { seoHead } from "@/lib/seo-head";

// The download is a browser navigation, not a client call: the API answers it
// with a redirect into the bucket, which an XHR would follow cross-origin. The
// generated type pins the path to the contract.
const downloadHref: DownloadBookData["url"] = "/api/book/download";

// The free book page, ported from legacy books/show. Legacy routed it under the
// ru suffix only, so `/book` and `/es/book` are not pages at all.
export const Route = createFileRoute("/{-$locale}/book")({
  beforeLoad: ({ params }) => {
    if (params.locale !== "ru") throw notFound();
  },
  loader: ({ context }) => context.queryClient.ensureQueryData(getBookOptions()),
  // Legacy meta: the header as title, the description cut to 160 characters as
  // Rails' truncate did, the canonical link and the social tags.
  head: ({ match }) => {
    const { i18n } = match.context;
    return seoHead({
      i18n,
      title: i18n.t(($) => $.books.show.header).trim(),
      description: truncate(
        i18n.t(($) => $.books.show.description),
        { length: 160 },
      ),
      canonicalPath: match.pathname,
    });
  },
  component: Show,
});

function Show() {
  const { t } = useTranslation();
  const { user } = Route.useRouteContext();
  const { data } = useQuery(getBookOptions());

  const features = [
    {
      key: "direction",
      title: t(($) => $.books.show.features.direction),
      explanation: t(($) => $.books.show.features.direction_explanation),
      icon: IconCompass,
    },
    {
      key: "plan",
      title: t(($) => $.books.show.features.plan),
      explanation: t(($) => $.books.show.features.plan_explanation),
      icon: IconList,
    },
    {
      key: "resume",
      title: t(($) => $.books.show.features.resume),
      explanation: t(($) => $.books.show.features.resume_explanation),
      icon: IconFileText,
    },
    {
      key: "interview",
      title: t(($) => $.books.show.features.interview),
      explanation: t(($) => $.books.show.features.interview_explanation),
      icon: IconUsers,
    },
  ] as const;

  return (
    <Container size="lg" my="xl">
      <Grid>
        <Grid.Col span={{ base: 12, md: 7 }}>
          <Text c="dimmed" fz="lg">
            {t(($) => $.books.show.freebook)}
          </Text>
          <Title order={1} mb="lg">
            {t(($) => $.books.show.header)}
          </Title>
          <Text fz="lg">{t(($) => $.books.show.description)}</Text>

          {data?.requested ? (
            <Button
              component="a"
              href={downloadHref}
              target="_blank"
              rel="noopener noreferrer"
              mt="md"
              variant="outline"
              size="lg"
            >
              {t(($) => $.books.show.download)}
            </Button>
          ) : user ? (
            <RequestButton />
          ) : (
            // Legacy's request posted through its auth guard, which sent a
            // visitor to sign in.
            <ButtonLink to="/{-$locale}/session/new" mt="md" variant="outline" size="lg">
              {t(($) => $.books.show.request)}
            </ButtonLink>
          )}

          <Grid mt={40} gap="md">
            {features.map(({ key, title, explanation, icon: Icon }) => (
              <Grid.Col span={{ base: 12, md: 6 }} key={key}>
                <Stack gap={4}>
                  <Box style={{ display: "flex", alignItems: "center" }}>
                    <Icon size={20} style={{ marginRight: 8 }} />
                    <Text fw={700} fz="lg">
                      {title}
                    </Text>
                  </Box>
                  <Text>{explanation}</Text>
                </Stack>
              </Grid.Col>
            ))}
          </Grid>
        </Grid.Col>

        <Grid.Col span={{ base: 12, md: 5 }}>
          <Image src={bookCoverImg} alt="Book cover" radius="md" fit="contain" />
        </Grid.Col>
      </Grid>

      <Title order={2} mt={40} mb={20}>
        {t(($) => $.books.show.toc)}
      </Title>

      {bookToc.map((item, index) => (
        <Box key={item.title} py="md">
          <Grid align="start">
            <Grid.Col span={{ base: 12, md: 2 }}>
              <Text fw={700}>{t(($) => $.books.show.chapter, { chapter: String(index + 1) })}</Text>
            </Grid.Col>
            <Grid.Col span={{ base: 12, md: 4 }}>{item.title}</Grid.Col>
            <Grid.Col span={{ base: 12, md: 6 }}>
              <ul style={{ margin: 0, paddingLeft: 20 }}>
                {item.subsections.map((subsection) => (
                  <li key={subsection}>{subsection}</li>
                ))}
              </ul>
            </Grid.Col>
          </Grid>
          {index !== bookToc.length - 1 && <Divider my="md" />}
        </Box>
      ))}

      {/* Legacy also required ru here; the page exists for ru only. */}
      {user && (
        <Grid align="center" mt={60}>
          <Grid.Col span={{ base: 12, lg: 7 }}>
            <Title order={2}>{t(($) => $.home.index.consultation)}</Title>
          </Grid.Col>
          <Grid.Col span={{ base: 12, lg: 5 }}>
            <Paper p="lg" radius="md" withBorder>
              <LeadFormBlock />
            </Paper>
          </Grid.Col>
        </Grid>
      )}
    </Container>
  );
}

// Legacy posted the request and redirected back to the page with a flash, for
// a repeat request too. Refetching the page state swaps in the download button.
function RequestButton() {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const request = useMutation({
    ...createBookRequestMutation(),
    onSuccess: async () => {
      notifications.show({
        color: "green",
        message: t(($) => $.flash.books.create_request.success),
      });
      await queryClient.invalidateQueries({ queryKey: getBookQueryKey() });
    },
    onError: () => notifications.show({ color: "red", message: t(($) => $.common.errors.network) }),
  });

  return (
    <Button
      mt="md"
      variant="outline"
      size="lg"
      loading={request.isPending}
      onClick={() => request.mutate({})}
    >
      {t(($) => $.books.show.request)}
    </Button>
  );
}
