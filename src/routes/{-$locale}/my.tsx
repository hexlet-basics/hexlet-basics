import {
  Alert,
  Anchor,
  Card,
  Container,
  EmptyState,
  Group,
  Image,
  Progress,
  SimpleGrid,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { IconBooks, IconMessage2 } from "@tabler/icons-react";
import { createFileRoute, Link } from "@tanstack/react-router";
import { Trans, useTranslation } from "react-i18next";
import { getMyDashboardOptions } from "@/client/@tanstack/react-query.gen";
import type { CourseCatalogItem, Enrollment } from "@/client/types.gen";
import CourseBlock from "@/components/CourseBlock";
import { requireAuth } from "@/lib/auth";

// The signed-in learner's dashboard, ported from legacy web/my#show: the
// courses in progress with where to resume, and the finished ones. Anonymous
// visitors are sent to sign in and back (ADR-0008: the guard runs on the server).
export const Route = createFileRoute("/{-$locale}/my")({
  beforeLoad: ({ context, location }) => {
    requireAuth(context.user, location.href);
  },
  loader: ({ context }) => context.queryClient.query(getMyDashboardOptions()),
  head: ({ match }) => ({
    meta: [{ title: match.context.i18n.t(($) => $.layouts.shared.nav.my) }],
  }),
  component: My,
});

function StartedCourse({
  enrollment,
  landingPage,
}: {
  enrollment: Enrollment;
  landingPage: CourseCatalogItem;
}) {
  const title = landingPage.header ?? landingPage.name ?? landingPage.course.slug;

  return (
    <Card
      renderRoot={(props) => (
        <Link
          to="/{-$locale}/languages/$slug"
          params={(prev) => ({ ...prev, slug: landingPage.slug })}
          {...props}
        />
      )}
      p={0}
      radius="md"
      withBorder
      td="none"
    >
      <Group wrap="nowrap" align="stretch">
        <Image
          visibleFrom="sm"
          alt={title}
          maw={150}
          w="auto"
          fit="contain"
          src={landingPage.coverUrl}
        />

        <Stack w="100%" p="sm" gap="xs">
          <Title order={3}>{title}</Title>

          {enrollment.nextLessonName && <Text>{enrollment.nextLessonName} →</Text>}

          <Progress.Root mt="auto">
            <Progress.Section
              aria-label={`${enrollment.completion}%`}
              value={enrollment.completion}
            />
          </Progress.Root>
        </Stack>
      </Group>
    </Card>
  );
}

function My() {
  const { t, i18n } = useTranslation();
  const { startedEnrollments, finishedEnrollments, landingPagesByCourseId } = Route.useLoaderData();

  return (
    <Container>
      <Title order={2} my="xl">
        {t(($) => $.my.show.started)}
      </Title>
      {startedEnrollments.length === 0 && (
        <EmptyState icon={<IconBooks />} title={t(($) => $.common.empty)} />
      )}

      <SimpleGrid cols={{ base: 1, xs: 2 }}>
        {startedEnrollments.map((enrollment) => {
          const landingPage = landingPagesByCourseId[enrollment.courseId];
          return landingPage ? (
            <StartedCourse key={enrollment.id} enrollment={enrollment} landingPage={landingPage} />
          ) : null;
        })}
      </SimpleGrid>

      <Title my="xl" order={2}>
        {t(($) => $.my.show.finished)}
      </Title>

      {finishedEnrollments.length > 0 && i18n.language === "ru" && (
        <Alert icon={<IconMessage2 />} mb="xl" px="xl">
          <Trans
            t={t}
            i18nKey={($) => $.my.show.add_review}
            components={{
              a: (
                <Anchor
                  href="https://taplink.cc/codebasics_reviews"
                  target="_blank"
                  rel="noopener noreferrer"
                />
              ),
            }}
          />
        </Alert>
      )}

      {finishedEnrollments.length === 0 && (
        <EmptyState icon={<IconBooks />} title={t(($) => $.common.empty)} />
      )}

      <SimpleGrid cols={{ base: 2, xs: 3, sm: 4 }}>
        {finishedEnrollments.map((enrollment) => {
          const landingPage = landingPagesByCourseId[enrollment.courseId];
          return landingPage ? <CourseBlock key={enrollment.id} item={landingPage} /> : null;
        })}
      </SimpleGrid>
    </Container>
  );
}
