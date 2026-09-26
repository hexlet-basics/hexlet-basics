import {
  Accordion,
  Alert,
  Box,
  Button,
  Center,
  Container,
  Grid,
  Group,
  Image,
  Loader,
  NumberFormatter,
  Progress,
  SimpleGrid,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { IconClock, IconUsers } from "@tabler/icons-react";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { getCourseOptions } from "@/client/@tanstack/react-query.gen";
import type { CourseView } from "@/client/types.gen";
import codeIllustration from "@/assets/code.svg";
import Breadcrumbs, { CurrentCrumb } from "@/components/Breadcrumbs";
import LessonMark from "@/components/lesson/LessonMark";
import { useEnterLesson } from "@/components/lesson/useEnterLesson";
import QnaBlock from "@/components/QnaBlock";
import { NavLink } from "@/components/RouterLink";
import { seoHead } from "@/lib/seo-head";

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
  loader: ({ context, params }) =>
    context.queryClient.ensureQueryData(getCourseOptions({ path: { slug: params.slug } })),
  // Legacy languages#show meta: the landing page's meta title and description,
  // the canonical link, the social block with the course cover, and Open Graph
  // typed as a website in the page's locale.
  //
  // og:url is the canonical. Legacy built it from the landing page's slug, the
  // slug it resolved the page by; this page is resolved by the course's slug,
  // which is also what legacy's canonical used, so the page's own URL is both.
  head: ({ loaderData, match }) => {
    if (!loaderData) return {};
    const { i18n } = match.context;
    const { course, landingPage } = loaderData;
    return seoHead({
      i18n,
      title: landingPage?.metaTitle ?? course.name ?? course.slug,
      description: landingPage?.metaDescription ?? "",
      canonicalPath: match.pathname,
      image: course.coverListVariant,
      openGraph: { type: "website", locale: i18n.language },
    });
  },
  component: Show,
});

function Show() {
  const { slug } = Route.useParams();
  const { t, i18n } = useTranslation();
  const { unfinished } = Route.useSearch();
  const { data: view, isPending, isError } = useQuery(getCourseOptions({ path: { slug } }));

  if (isPending) {
    return (
      <Center h="50vh">
        <Loader />
      </Center>
    );
  }

  if (isError) {
    return (
      <Center h="50vh">
        <Text c="red">{t(($) => $.common.errors.server)}</Text>
      </Center>
    );
  }

  const { course, landingPage, progress } = view;
  const header = landingPage?.header ?? course.name ?? course.slug;
  const name = landingPage?.name ?? course.name ?? course.slug;
  const updatedAt = course.currentVersion?.createdAt;

  return (
    <Container size="lg">
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
          <Text size="lg" mb="sm">
            {landingPage?.description}
          </Text>
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
            <Text>{landingPage.usedInDescription}</Text>
          </Grid.Col>
        </Grid>
      )}

      {landingPage?.outcomesHeader && (
        <SimpleGrid cols={{ base: 1, lg: 2 }} py="xl">
          {landingPage.outcomesImage && (
            <Image
              src={landingPage.outcomesImage}
              width="100%"
              height="auto"
              loading="lazy"
              alt={t(($) => $.courses.show.learning_preview)}
              style={{
                borderRadius: "var(--mantine-radius-xl)",
                boxShadow: "var(--mantine-shadow-lg)",
              }}
            />
          )}
          <Box>
            <Title order={2} size="h1" mb="md">
              {landingPage.outcomesHeader}
            </Title>
            <Text>{landingPage.outcomesDescription}</Text>
          </Box>
        </SimpleGrid>
      )}

      <Box my="xl">
        <Title order={2} size="h1" mb="md">
          {t(($) => $.courses.show.learning_program, { name })}
        </Title>

        {/* Completion as the server computed it: a share of the current
            version's lessons, which this page has no business recounting. */}
        {progress && (
          <Stack gap={4} mb="md" maw={320}>
            <Text size="sm">
              {t(($) => $.courses.show.progress.completion, { completion: progress.completion })}
            </Text>
            <Progress value={progress.completion} aria-hidden="true" />
          </Stack>
        )}

        <LearningProgram view={view} />
      </Box>

      <Box my={{ base: "lg", sm: "xxl" }}>
        <Title order={2} fz="h1" mb="xl">
          {t(($) => $.courses.show.about_learning)}
        </Title>
        <Text fw="bold">{t(($) => $.courses.show["convenient format"])}</Text>
        <Text mb="md">{t(($) => $.courses.show.learning_conveniently)}</Text>
        <Text fw="bold">{t(($) => $.courses.show.browser_practice)}</Text>
        <Text mb="md">{t(($) => $.courses.show.real_life_challenges)}</Text>
        <Text fw="bold">{t(($) => $.courses.show.ai_without_limits)}</Text>
        <Text mb="md">{t(($) => $.courses.show.ai_explanation)}</Text>
      </Box>

      <QnaBlock items={view.qnaItems} />
    </Container>
  );
}

// The learning program, as legacy's accordion laid it out: a panel per module,
// the first one open, each with its lessons beside the module's description.
//
// A lesson in the panel is a link carrying its mark. Checks and locks come from
// `progress`, names and order from `lessons`, joined by slug — the same pair
// the player's list renders; a module names its lessons by slug the same way.
// A lesson no module in this locale claims is listed after the panels rather
// than dropped, and a course with no modules is the flat list alone.
function LearningProgram({ view }: { view: CourseView }) {
  const lessonBySlug = new Map(view.lessons.map((item) => [item.slug, item]));
  const stateBySlug = new Map(view.progress?.lessons.map((item) => [item.slug, item]) ?? []);
  const claimed = new Set(view.modules.flatMap((module) => module.lessonSlugs));
  const unclaimed = view.lessons.filter((item) => !claimed.has(item.slug));

  // A locked lesson is still a link: theory is public, and the lock says "not
  // yet", never "you cannot read this".
  const lessonLinks = (lessons: CourseView["lessons"]) =>
    lessons.map((item) => (
      <NavLink
        key={item.slug}
        to="/{-$locale}/languages/$slug/lessons/$lessonSlug"
        params={{ slug: view.course.slug, lessonSlug: item.slug }}
        label={item.name}
        leftSection={<LessonMark state={stateBySlug.get(item.slug)} />}
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
                  <Grid.Col span={{ base: 12, xs: 4 }}>
                    {lessonLinks(
                      module.lessonSlugs.flatMap((slug) => lessonBySlug.get(slug) ?? []),
                    )}
                  </Grid.Col>
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
function CourseAction({ view }: { view: CourseView }) {
  const { t } = useTranslation();
  const { course, progress } = view;
  const { enter, isPending } = useEnterLesson(course.slug);

  if (!progress) return null;

  if (progress.nextLessonSlug === null) {
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

  const lesson = view.lessons.find((item) => item.slug === progress.nextLessonSlug);
  if (!lesson) return null;

  // Navigation waits on the command's success: the lesson is entered started or
  // not at all.
  return (
    <Button size="lg" loading={isPending} onClick={() => enter(lesson)}>
      {progress.furthestFinishedPosition === 0
        ? t(($) => $.courses.show.try)
        : t(($) => $.courses.show.continue)}
    </Button>
  );
}

// The Hexlet program the course leads on to, tagged as a referral from here as
// legacy tagged it. The parameters are set on the parsed URL rather than
// appended, so a landing address that already carries a query stays well formed.
function programLink(landingPage: string): string {
  // An address the admin form let through malformed is linked as written
  // rather than taking the page down with it.
  if (!URL.canParse(landingPage)) return landingPage;
  const url = new URL(landingPage);
  url.searchParams.set("utm_source", "code-basics");
  url.searchParams.set("utm_medium", "referral");
  return url.toString();
}
