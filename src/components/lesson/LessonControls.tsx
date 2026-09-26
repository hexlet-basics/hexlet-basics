import { ActionIcon, Box, Button, Divider, Group, Stack, Text } from "@mantine/core";
import { modals } from "@mantine/modals";
import { IconPlayerPlay, IconRepeat } from "@tabler/icons-react";
import { useLocation, useRouteContext } from "@tanstack/react-router";
import { Trans, useTranslation } from "react-i18next";
import type { CourseLessonView } from "@/client/types.gen";
import { useEnterLesson } from "@/components/lesson/useEnterLesson";
import { ButtonLink, TextLink } from "@/components/RouterLink";

// The bar under the workspace, where the buttons that act on the exercise live:
// reset, Previous, Run, Next — legacy's order. It sits below the tabs rather than
// inside the editor pane so that running a solution is one press away from
// whichever tab the learner is on.
export default function LessonControls({
  view,
  passed,
  onReset,
  onRun,
  running,
}: {
  view: CourseLessonView;
  // Passed in this visit or finished before it — the same fact that opens the
  // reference solution, so Next and the solution can never disagree.
  passed: boolean;
  onReset: () => void;
  onRun: () => void;
  running: boolean;
}) {
  const { t } = useTranslation();
  const user = useRouteContext({ from: "__root__", select: (context) => context.user });
  const courseSlug = view.lesson.course.slug;

  // Neighbours come from the ordered lesson list the payload already carries,
  // not from dedicated fields: the list is the course order, and deriving both
  // ends from it keeps Previous, Next and the navigation tab in agreement.
  const index = view.lessons.findIndex((item) => item.slug === view.lesson.slug);
  const prevLesson = index > 0 ? view.lessons[index - 1] : undefined;
  const nextLesson = index >= 0 ? view.lessons[index + 1] : undefined;

  // Resetting throws away whatever the learner has written, and there is no undo
  // once the buffer is overwritten — so it asks first.
  const confirmReset = () =>
    modals.openConfirmModal({
      title: t(($) => $.common.confirm),
      labels: { confirm: t(($) => $.common.boolean.yes), cancel: t(($) => $.common.boolean.no) },
      onConfirm: onReset,
    });

  return (
    <Stack gap={0}>
      <Divider />
      <Box py="sm">
        <Group justify="center">
          <ActionIcon
            variant="light"
            size="lg"
            onClick={confirmReset}
            aria-label={t(($) => $.helpers.reset)}
          >
            <IconRepeat size={18} />
          </ActionIcon>

          {/* A plain link: going back to re-read starts nothing. On the first
              lesson there is nowhere to go, so it is a disabled button rather
              than an anchor to nowhere. */}
          {prevLesson ? (
            <ButtonLink
              variant="outline"
              color="green"
              to="/{-$locale}/languages/$slug/lessons/$lessonSlug"
              params={{ slug: courseSlug, lessonSlug: prevLesson.slug }}
            >
              {t(($) => $.courses.lessons.show.prev)}
            </ButtonLink>
          ) : (
            <Button variant="outline" color="green" disabled>
              {t(($) => $.courses.lessons.show.prev)}
            </Button>
          )}

          {/* While a check is in flight the button says so and refuses a second
              press: one solution is running, and submitting it twice would tell
              the learner nothing new. */}
          <Button
            leftSection={<IconPlayerPlay size={18} />}
            onClick={onRun}
            loading={running}
            disabled={running}
          >
            {t(($) => $.courses.lessons.show.controls.run)}
          </Button>

          {user ? (
            <ForwardButton view={view} nextLesson={nextLesson} passed={passed} />
          ) : (
            <GuestSignUp passed={passed} />
          )}
        </Group>

        {/* The one thing this port adds rather than copies: legacy carried the
            string and rendered it nowhere. The merge that keeps a guest's
            progress on sign-up is built, so the promise is now true. */}
        {!user && passed && <GuestPrompt />}
      </Box>
    </Stack>
  );
}

// Next, for a signed-in learner — or, on the last lesson, the way back to the
// course. Unavailable until the lesson is passed, as in legacy, so the button
// tells the learner what is expected; `passed` already covers a lesson finished
// on an earlier visit, because revisiting is not the same as being stuck.
function ForwardButton({
  view,
  nextLesson,
  passed,
}: {
  view: CourseLessonView;
  nextLesson: CourseLessonView["lessons"][number] | undefined;
  passed: boolean;
}) {
  const { t } = useTranslation();
  const courseSlug = view.lesson.course.slug;
  // Next starts the following lesson and only then goes there.
  const { enter, isPending } = useEnterLesson(courseSlug);

  // The last lesson reads as completion and returns to the course page. The
  // dedicated completion page has no contract operation yet.
  if (!nextLesson) {
    return passed ? (
      <ButtonLink
        variant="outline"
        color="green"
        to="/{-$locale}/languages/$slug"
        params={{ slug: courseSlug }}
      >
        {t(($) => $.courses.lessons.show.finish)}
      </ButtonLink>
    ) : (
      <Button variant="outline" color="green" disabled>
        {t(($) => $.courses.lessons.show.finish)}
      </Button>
    );
  }

  return (
    <Button
      variant="outline"
      color="green"
      disabled={!passed}
      loading={isPending}
      onClick={() => enter(nextLesson)}
    >
      {t(($) => $.courses.lessons.show.next)}
    </Button>
  );
}

// Where Next sits for a guest, as in legacy: a link to sign up, carrying the
// way back to this lesson, and unavailable until the lesson is passed.
function GuestSignUp({ passed }: { passed: boolean }) {
  const { t } = useTranslation();
  const redirect = useLocation({ select: (location) => location.href });

  return passed ? (
    <ButtonLink variant="outline" color="green" to="/{-$locale}/users/new" search={{ redirect }}>
      {t(($) => $.courses.lessons.show.next)}
    </ButtonLink>
  ) : (
    <Button variant="outline" color="green" disabled>
      {t(($) => $.courses.lessons.show.next)}
    </Button>
  );
}

// The line a guest reads after passing: signing up keeps what they have done.
function GuestPrompt() {
  const { t } = useTranslation();
  const redirect = useLocation({ select: (location) => location.href });

  return (
    <Text ta="center" size="sm" mt="xs">
      <Trans
        t={t}
        i18nKey={($) => $.courses.lessons.show.sign_up_for_tracking_progress_html}
        components={{ a: <TextLink to="/{-$locale}/users/new" search={{ redirect }} inherit /> }}
      />
    </Text>
  );
}
