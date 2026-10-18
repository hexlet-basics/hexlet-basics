import {
  Box,
  Button,
  Container,
  Divider,
  Grid,
  Group,
  Image,
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
import { ButtonLink } from "@/components/RouterLink";
import bookToc from "@/lib/book";
import { seoHead } from "@/lib/seo-head";
import type { ReactElement } from "react";

// The download is a browser navigation, not a client call: the API answers it
// with a redirect into the bucket, which an XHR would follow cross-origin. The
// generated type pins the path to the contract.
const downloadHref: DownloadBookData["url"] = "/api/book/download";

// The free book page, ported from legacy books/show. Legacy routed it under the
// ru suffix only, so `/book` and `/es/book` are not pages at all.
export const Route = createFileRoute("/{-$locale}/book")({
  beforeLoad: ({ params }) => {
    if (params.locale !== "ru") {
      throw notFound();
    }
  },
  component: Show,
  // Legacy meta: the header as title, the description cut to 160 characters as
  // Rails' truncate did, the canonical link and the social tags.
  head: ({ match }) => {
    const { i18n } = match.context;
    return seoHead({
      canonicalPath: match.pathname,
      description: truncate(
        i18n.t(($) => $.books.show.description),
        { length: 160 },
      ),
      i18n,
      title: i18n.t(($) => $.books.show.header).trim(),
    });
  },
  loader: ({ context }) => context.queryClient.query({ ...getBookOptions(), staleTime: "static" }),
});

function Show(): ReactElement {
  const { t } = useTranslation();
  const { user } = Route.useRouteContext();
  const { data } = useQuery(getBookOptions());

  const features = [
    {
      explanation: t(($) => $.books.show.features.direction_explanation),
      icon: IconCompass,
      key: "direction",
      title: t(($) => $.books.show.features.direction),
    },
    {
      explanation: t(($) => $.books.show.features.plan_explanation),
      icon: IconList,
      key: "plan",
      title: t(($) => $.books.show.features.plan),
    },
    {
      explanation: t(($) => $.books.show.features.resume_explanation),
      icon: IconFileText,
      key: "resume",
      title: t(($) => $.books.show.features.resume),
    },
    {
      explanation: t(($) => $.books.show.features.interview_explanation),
      icon: IconUsers,
      key: "interview",
      title: t(($) => $.books.show.features.interview),
    },
  ] as const;

  function renderAction(): ReactElement {
    if (data?.requested) {
      return (
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
      );
    }
    if (user) {
      return <RequestButton />;
    }
    // Legacy's request posted through its auth guard, which sent a visitor to
    // sign in.
    return (
      <ButtonLink to="/{-$locale}/session/new" mt="md" variant="outline" size="lg">
        {t(($) => $.books.show.request)}
      </ButtonLink>
    );
  }
  const action = renderAction();

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

          {action}

          <Grid mt={40} gap="md">
            {features.map((feature) => (
              <Grid.Col span={{ base: 12, md: 6 }} key={feature.key}>
                <Feature {...feature} />
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
              <Box component="ul" m={0} pl={20}>
                {item.subsections.map((subsection) => (
                  <li key={subsection}>{subsection}</li>
                ))}
              </Box>
            </Grid.Col>
          </Grid>
          {index !== bookToc.length - 1 && <Divider my="md" />}
        </Box>
      ))}
    </Container>
  );
}

// Legacy posted the request and redirected back to the page with a flash, for
// a repeat request too. Refetching the page state swaps in the download button.
function RequestButton(): ReactElement {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const request = useMutation({
    ...createBookRequestMutation(),
    onError: () => notifications.show({ color: "red", message: t(($) => $.common.errors.network) }),
    onSuccess: async () => {
      notifications.show({
        color: "green",
        message: t(($) => $.flash.books.create_request.success),
      });
      await queryClient.invalidateQueries({ queryKey: getBookQueryKey() });
    },
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

// One thing the book teaches: an icon and a title, then what it means.
function Feature({
  title,
  explanation,
  icon: Icon,
}: {
  title: string;
  explanation: string;
  icon: typeof IconCompass;
}): ReactElement {
  return (
    <Stack gap={4}>
      <Group gap={8} wrap="nowrap">
        <Icon size={20} />
        <Text fw={700} fz="lg">
          {title}
        </Text>
      </Group>
      <Text>{explanation}</Text>
    </Stack>
  );
}
