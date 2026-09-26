import {
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
import { notifications } from "@mantine/notifications";
import { IconClock, IconUsers } from "@tabler/icons-react";
import { useMutation, useQuery } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { getCourseOptions, startLessonMutation } from "@/client/@tanstack/react-query.gen";
import type { CourseView } from "@/client/types.gen";
import codeIllustration from "@/assets/code.svg";
import Breadcrumbs, { CurrentCrumb } from "@/components/Breadcrumbs";
import LessonMark from "@/components/lesson/LessonMark";
import { NavLink } from "@/components/RouterLink";

// The Course page, at its legacy URL (ADR-0002): the course's landing copy, its
// current lessons, and one button that puts the learner in it.
//
// Ported from legacy `languages/show`. What did not come across is what the
// contract does not carry or what would branch on having an account: the
// module accordion (the payload's lesson list is flat), the Q&A block, the
// lead and sign-up forms, the promo video, and the two extra calls to action,
// which would each be a second way in that skips the start command.
//
// The loader prefetches into the request-scoped QueryClient, so the landing
// copy is in the server-rendered HTML (ADR-0008). Loading this page starts
// nothing: the router preloads on hover, so progress begins only when the
// button is pressed (ADR-0012).
export const Route = createFileRoute("/{-$locale}/languages/$slug/")({
  loader: ({ context, params }) =>
    context.queryClient.ensureQueryData(getCourseOptions({ path: { slug: params.slug } })),
  component: Show,
});

function Show() {
  const { slug } = Route.useParams();
  const { t, i18n } = useTranslation();
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

  // Checks and locks come from `progress`, names and order from `lessons`,
  // joined by slug — the same pair the player's list renders.
  const stateBySlug = new Map(progress?.lessons.map((item) => [item.slug, item]) ?? []);

  return (
    <Container size="lg">
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
              {t(($) => $.courses.show.completion, { completion: progress.completion })}
            </Text>
            <Progress value={progress.completion} aria-hidden="true" />
          </Stack>
        )}

        {/* A locked lesson is still a link: theory is public, and the lock
            says "not yet", never "you cannot read this". */}
        {view.lessons.map((item) => (
          <NavLink
            key={item.slug}
            to="/{-$locale}/languages/$slug/lessons/$lessonSlug"
            params={{ slug: course.slug, lessonSlug: item.slug }}
            label={item.name}
            leftSection={<LessonMark state={stateBySlug.get(item.slug)} />}
          />
        ))}
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
    </Container>
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
  const { locale } = Route.useParams();
  const navigate = useNavigate();
  const start = useMutation({
    ...startLessonMutation(),
    // A refusal or a network failure leaves the visitor where they are; landing
    // in a lesson that was never started would be worse than staying put.
    onError: () => notifications.show({ message: t(($) => $.common.errors.network) }),
  });

  const { course, progress } = view;
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
            href={`${course.hexletProgramLandingPage}?utm_source=code-basics&utm_medium=referral`}
            target="_blank"
            rel="noopener noreferrer"
          >
            {t(($) => $.courses.show.hexlet_program_link)}
          </Button>
        )}
      </Group>
    );
  }

  const lessonSlug = progress.nextLessonSlug;
  const lesson = view.lessons.find((item) => item.slug === lessonSlug);
  if (!lesson) return null;

  // Navigation waits on the command's success: the lesson is entered started or
  // not at all.
  const enter = () =>
    start.mutate(
      { path: { id: lesson.id } },
      {
        onSuccess: () =>
          navigate({
            to: "/{-$locale}/languages/$slug/lessons/$lessonSlug",
            params: { locale, slug: course.slug, lessonSlug },
          }),
      },
    );

  return (
    <Button size="lg" loading={start.isPending} onClick={enter}>
      {progress.furthestFinishedPosition === 0
        ? t(($) => $.courses.show.try)
        : t(($) => $.courses.show.continue)}
    </Button>
  );
}
