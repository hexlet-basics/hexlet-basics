import {
  Accordion,
  Alert,
  Box,
  Button,
  Card,
  Center,
  Container,
  Flex,
  Grid,
  Group,
  Image,
  NumberFormatter,
  Paper,
  Progress,
  SimpleGrid,
  Skeleton,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import type { ButtonProps } from "@mantine/core";
import { IconClock, IconUsers } from "@tabler/icons-react";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute, redirect } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { getActiveBannerOptions, getCourseOptions } from "@/client/@tanstack/react-query.gen";
import type { CourseView } from "@/client/types.gen";
import codeIllustration from "@/assets/code.svg";
import learningVideo from "@/assets/course-landing-page/learning_ru.mp4";
import { Banner } from "@/components/Banner";
import Breadcrumbs, { CurrentCrumb } from "@/components/Breadcrumbs";
import MarkdownViewer from "@/components/MarkdownViewer";
import LessonMark from "@/components/lesson/LessonMark";
import { useEnterLesson } from "@/components/lesson/useEnterLesson";
import QnaBlock from "@/components/QnaBlock";
import { NavLink } from "@/components/RouterLink";
import { seoHead } from "@/lib/seo-head";
import { rethrowMissingAsNotFound } from "@/lib/api-error";
import { orNull } from "@/lib/lenient";
import type { ReactElement } from "react";

// The Course page, at its legacy URL (ADR-0002): the course's landing copy, its
// current lessons, and one button that puts the learner in it.
//
// Ported from legacy `languages/show`, module accordion and Q&A included. What
// did not come across is what would branch on having an account or add a way
// in: the lead and sign-up forms, the promo video, and the two extra calls to
// action, which would each be a second way in that skips the start command.
//
// The loader prefetches into the request-scoped QueryClient, so the landing
// copy is in the server-rendered HTML (ADR-0008). Loading this page starts
// nothing: the router preloads on hover, so progress begins only when the
// button is pressed (ADR-0012).
//
// `unfinished` is set by the completion page when it bounces a learner whose
// Enrollment is not finished: legacy's flash warning, carried in the URL so it
// survives a server-side redirect.
export const Route = createFileRoute("/{-$locale}/languages/$slug/")({
  validateSearch: z.object({ unfinished: z.boolean().optional() }),

  loader: async ({ context, params }) => {
    // The banner is prefetched so it is in the server-rendered HTML; fetched on
    // the client only, it would pop in after hydration. It is decoration, so
    // failing to read it must not take the page down.
    const banner = orNull(
      context.queryClient.query({ ...getActiveBannerOptions(), staleTime: "static" }),
    );
    try {
      const [view] = await Promise.all([
        // Carries the visitor's progress, so it takes the client's default
        // staleness and is dropped when who is visiting changes.
        context.queryClient.query(getCourseOptions({ path: { slug: params.slug } })),
        banner,
      ]);
      // An archived landing page moved: legacy answered it with a permanent
      // redirect, which keeps old links and search results pointing somewhere.
      if (view.movedTo) {
        throw redirect({
          params: { ...params, slug: view.movedTo },
          statusCode: 301,
          to: "/{-$locale}/languages/$slug",
        });
      }
      return view;
    } catch (error) {
      // The slug is a landing page's, in this locale's site: an unknown one and
      // one of another locale are the same to a visitor, there is no such page.
      // Thrown here, not left to the error page, so the answer is a 404.
      return rethrowMissingAsNotFound(error);
    }
  },

  // Legacy languages#show meta: the landing page's meta title and description,
  // the canonical link, the social block with the course cover, and Open Graph
  // typed as a website in the page's locale.
  //
  // og:url is the canonical. Legacy built it from the landing page's slug, the
  // slug it resolved the page by; this page is resolved by the course's slug,
  // which is also what legacy's canonical used, so the page's own URL is both.
  head: ({ loaderData, match }) => {
    if (!loaderData) {
      return {};
    }
    const { i18n } = match.context;
    const { course, landingPage } = loaderData;
    return seoHead({
      canonicalPath: match.pathname,
      description: landingPage?.metaDescription ?? "",
      i18n,
      image: course.coverListVariant,
      openGraph: { locale: i18n.language, type: "website" },
      title: landingPage?.metaTitle ?? course.name ?? course.slug,
    });
  },

  component: Show,
});

function Show(): ReactElement {
  const { slug } = Route.useParams();
  const { t, i18n } = useTranslation();
  const { unfinished } = Route.useSearch();
  const { user } = Route.useRouteContext();
  const { data: view, isPending, isError } = useQuery(getCourseOptions({ path: { slug } }));

  if (isPending) {
    return (
      <Container my="xl">
        <Skeleton h={48} w="50%" mb="md" />
        <Skeleton h={16} w="70%" mb="xl" />
        <Skeleton h={280} />
      </Container>
    );
  }

  if (isError) {
    return (
      <Center h="50vh">
        <Text c="red">{t(($) => $.common.errors.server)}</Text>
      </Center>
    );
  }

  const { course, landingPage, progress, enrollment } = view;
  const header = landingPage?.header ?? course.name ?? course.slug;
  const name = landingPage?.name ?? course.name ?? course.slug;
  const updatedAt = course.currentVersion?.createdAt;

  return (
    <>
      <Banner />
      <Container size="lg" my="xl">
        {unfinished && (
          <Alert color="yellow" mb="md">
            {t(($) => $.flash.courses.success.warning)}
          </Alert>
        )}
        <Box mb="lg">
          <Breadcrumbs homeLabel={t(($) => $.courses.show.to_home_title)}>
            <CurrentCrumb>{header}</CurrentCrumb>
          </Breadcrumbs>
        </Box>

        <Grid>
          <Grid.Col span={{ base: 12, sm: 7 }}>
            <Text size="sm" c="dimmed">
              {t(($) => $.courses.show.free_course, { name: header })}
            </Text>
            <Title order={1} mb="lg" fz={48}>
              {header}
            </Title>
            <Box fz="lg" mb="sm">
              <MarkdownViewer>{landingPage?.description ?? ""}</MarkdownViewer>
            </Box>
            <Group mb="xl">
              <Group me="lg">
                <IconUsers size={16} />
                <NumberFormatter thousandSeparator value={course.enrollmentsCount} />
              </Group>
              {updatedAt && (
                <Group>
                  <IconClock size={16} />
                  {t(($) => $.courses.show.updated_at, {
                    // UTC on both sides, so the server's date and the browser's
                    // hydrated one cannot disagree across a midnight.
                    date: new Date(updatedAt).toLocaleDateString(i18n.language, {
                      dateStyle: "long",
                      timeZone: "UTC",
                    }),
                  })}
                </Group>
              )}
            </Group>

            <CourseAction view={view} />
          </Grid.Col>
          <Grid.Col visibleFrom="sm" span={{ base: 12, sm: 5 }}>
            <Image src={codeIllustration} fit="cover" alt={t(($) => $.courses.show.cover_image)} />
          </Grid.Col>
        </Grid>

        {landingPage?.usedInHeader && (
          <Grid my="xl">
            <Grid.Col span={{ base: 12, lg: 9 }}>
              <Title order={2} size="h1" mb="md">
                {landingPage.usedInHeader}
              </Title>
              <MarkdownViewer>{landingPage.usedInDescription ?? ""}</MarkdownViewer>
            </Grid.Col>
          </Grid>
        )}

        <OutcomesSection landingPage={landingPage} />

        <Card withBorder my={{ base: "lg", sm: "xxl" }} p="xl">
          <Group justify="space-between">
            <Text fz="h3">{t(($) => $.courses.show.course_graduates)}</Text>
            <EnterButton view={view} bg="dark" size="xl">
              {t(($) => $.courses.show.start, { name })}
            </EnterButton>
          </Group>
        </Card>

        <Box my="xl">
          <Title order={2} size="h1" mb="md">
            {t(($) => $.courses.show.learning_program, { name })}
          </Title>

          {/* Completion as the server computed it: a share of the current
            version's lessons, which this page has no business recounting. A
            visitor who has not started has nothing to show, so only an
            enrolled learner sees it. */}
          {progress && enrollment && (
            <Stack gap={4} mb="md" maw={320}>
              <Text size="sm">
                {t(($) => $.courses.show.progress.completion, { completion: progress.completion })}
              </Text>
              <Progress value={progress.completion} aria-hidden="true" />
            </Stack>
          )}

          <LearningProgram view={view} />
        </Box>

        <AboutLearning />

        {/* The demo pitch is for visitors: a signed-in learner already has the
          hero button, and is not being invited to try anything. */}
        {!user && (
          <Card bg="indigo.0" p="xl">
            <Group justify="space-between">
              <Text fz="h3">{t(($) => $.courses.show.demo_description)}</Text>
              <Flex h="100%" justify="end">
                <EnterButton view={view} my="auto" size="xl">
                  {t(($) => $.courses.show.demo_start)}
                </EnterButton>
              </Flex>
            </Group>
          </Card>
        )}

        <QnaBlock items={view.qnaItems} />
      </Container>
    </>
  );
}

// What the learner takes away from the course: a header and description, with the
// preview image beside them when the landing page has one.
function OutcomesSection({
  landingPage,
}: {
  landingPage: CourseView["landingPage"];
}): ReactElement | null {
  const { t } = useTranslation();

  if (!landingPage?.outcomesHeader) {
    return null;
  }

  return (
    <SimpleGrid cols={{ base: 1, lg: 2 }} py="xl">
      {landingPage.outcomesImage && (
        <Paper shadow="lg" radius="xl" bg="transparent">
          <Image
            src={landingPage.outcomesImage}
            width="100%"
            height="auto"
            loading="lazy"
            alt={t(($) => $.courses.show.learning_preview)}
            radius="xl"
          />
        </Paper>
      )}
      <Box>
        <Title order={2} size="h1" mb="md">
          {landingPage.outcomesHeader}
        </Title>
        <MarkdownViewer>{landingPage.outcomesDescription ?? ""}</MarkdownViewer>
      </Box>
    </SimpleGrid>
  );
}

// Why the course is built the way it is; the Russian site also shows the video.
function AboutLearning(): ReactElement {
  const { t, i18n } = useTranslation();

  return (
    <Box my={{ base: "lg", sm: "xxl" }}>
      <Title order={2} fz="h1" mb="xl">
        {t(($) => $.courses.show.about_learning)}
      </Title>
      <SimpleGrid cols={{ base: 1, lg: i18n.language === "ru" ? 2 : 1 }}>
        <Box>
          <Text fw="bold">{t(($) => $.courses.show["convenient format"])}</Text>
          <Text mb="md">{t(($) => $.courses.show.learning_conveniently)}</Text>
          <Text fw="bold">{t(($) => $.courses.show.browser_practice)}</Text>
          <Text mb="md">{t(($) => $.courses.show.real_life_challenges)}</Text>
          <Text fw="bold">{t(($) => $.courses.show.ai_without_limits)}</Text>
          <Text mb="md">{t(($) => $.courses.show.ai_explanation)}</Text>
        </Box>
        {i18n.language === "ru" && (
          <Box
            component="video"
            h="100%"
            w="100%"
            src={learningVideo}
            autoPlay
            loop
            muted
            playsInline
          />
        )}
      </SimpleGrid>
    </Box>
  );
}

// The learning program, as legacy's accordion laid it out: a panel per module,
// the first one open, each with its lessons beside the module's description.
//
// A lesson in the panel is a link carrying its mark; each row already knows
// whether it is finished, so nothing is joined here. A lesson no module in this
// locale claims is listed after the panels rather than dropped, and a course
// with no modules is the flat list alone.
function LearningProgram({ view }: { view: CourseView }): ReactElement {
  const claimed = new Set(
    view.modules.flatMap((module) => module.lessons.map((item) => item.slug)),
  );
  const unclaimed = view.lessons.filter((item) => !claimed.has(item.slug));

  const lessonLinks = (lessons: CourseView["lessons"]): ReactElement[] =>
    lessons.map((item) => (
      <NavLink
        key={item.slug}
        to="/{-$locale}/languages/$slug/lessons/$lessonSlug"
        params={{ lessonSlug: item.slug, slug: view.course.slug }}
        label={item.name}
        leftSection={<LessonMark finished={item.finished} />}
      />
    ));

  return (
    <>
      {view.modules.length > 0 && (
        <Accordion defaultValue={String(view.modules[0]?.id)}>
          {view.modules.map((module) => (
            <Accordion.Item key={module.id} value={String(module.id)} py="lg">
              <Accordion.Control>
                <Title order={3}>{module.name}</Title>
              </Accordion.Control>
              <Accordion.Panel>
                <Grid>
                  <Grid.Col span={{ base: 12, xs: 4 }}>{lessonLinks(module.lessons)}</Grid.Col>
                  <Grid.Col span={{ base: 12, xs: 8 }}>{module.description}</Grid.Col>
                </Grid>
              </Accordion.Panel>
            </Accordion.Item>
          ))}
        </Accordion>
      )}
      {lessonLinks(unclaimed)}
    </>
  );
}

// The one way into the course. Its label and target are read from the progress
// payload and nothing else — not from the enrollment, not from whether there is
// an account — so a guest, whose position is the cookie, gets the same button a
// signed-in learner does:
//
// - nothing finished: Try, into the Next Lesson (the first one);
// - something finished: Continue, into the Next Lesson;
// - no Next Lesson: the course is finished and no lesson is offered.
//
// It is a button rather than a link because pressing it issues the start command
// before navigating. A link would enter the lesson without starting it, and
// would be preloaded on hover besides.
function CourseAction({ view }: { view: CourseView }): ReactElement | null {
  const { t } = useTranslation();
  const { course, progress } = view;
  const { enter, isPending } = useEnterLesson(course.slug);

  if (!progress) {
    return null;
  }

  if (progress.nextLesson === null) {
    // No lesson to offer is the end of the course only when it is finished: a
    // next lesson with no page in this locale is no button, not a finish.
    if (progress.state !== "finished") {
      return null;
    }
    return (
      <Group>
        <Text size="lg" fw="bold">
          {t(($) => $.shared.courses.course_finished)}
        </Text>
        {course.hexletProgramLandingPage && (
          <Button
            size="lg"
            component="a"
            href={programLink(course.hexletProgramLandingPage)}
            target="_blank"
            rel="noopener noreferrer"
          >
            {t(($) => $.courses.show.hexlet_program_link)}
          </Button>
        )}
      </Group>
    );
  }

  const lesson = progress.nextLesson;

  // Navigation waits on the command's success: the lesson is entered started or
  // not at all.
  return (
    <Button
      size="lg"
      loading={isPending}
      onClick={() => {
        enter(lesson);
      }}
    >
      {view.lessons.some((item) => item.finished)
        ? t(($) => $.courses.show.continue)
        : t(($) => $.courses.show.try)}
    </Button>
  );
}

// A secondary way into the course, as legacy's page had two more beside the
// hero's: it goes through the same start command as the hero button, into the
// learner's Next Lesson (the first one for a visitor), so none of them is a
// way in that skips the start.
function EnterButton({ view, ...props }: { view: CourseView } & ButtonProps): ReactElement | null {
  const { enter, isPending } = useEnterLesson(view.course.slug);
  const lesson = view.progress?.nextLesson ?? view.lessons[0];
  if (!lesson) {
    return null;
  }
  return (
    <Button
      {...props}
      loading={isPending}
      onClick={() => {
        enter(lesson);
      }}
    />
  );
}

// The Hexlet program the course leads on to, tagged as a referral from here as
// legacy tagged it. The parameters are set on the parsed URL rather than
// appended, so a landing address that already carries a query stays well formed.
function programLink(landingPage: string): string {
  // An address the admin form let through malformed is linked as written
  // rather than taking the page down with it.
  if (!URL.canParse(landingPage)) {
    return landingPage;
  }
  const url = new URL(landingPage);
  url.searchParams.set("utm_source", "code-basics");
  url.searchParams.set("utm_medium", "referral");
  return url.toString();
}
