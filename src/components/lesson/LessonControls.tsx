import { ActionIcon, Box, Button, Divider, Group, Stack } from "@mantine/core";
import { modals } from "@mantine/modals";
import { IconPlayerPlay, IconRepeat } from "@tabler/icons-react";
import { useLocation, useRouteContext } from "@tanstack/react-router";
import type { RegisteredRouter, ValidateLinkOptions } from "@tanstack/react-router";
import type { ReactNode } from "react";
import { useTranslation } from "react-i18next";
import type { CourseLessonView, LessonRef } from "@/client/types.gen";
import { useEnterLesson } from "@/components/lesson/useEnterLesson";
import { ButtonLink } from "@/components/RouterLink";

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

  // The neighbours in course order are the server's answer, with the course's
  // own order behind it, not something derived here from a position in a list.
  const prevLesson = view.adjacent.previous ?? undefined;
  const nextLesson = view.adjacent.next ?? undefined;

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
          <StepLink
            enabled={Boolean(prevLesson)}
            label={t(($) => $.courses.lessons.show.prev)}
            linkOptions={{
              to: "/{-$locale}/languages/$slug/lessons/$lessonSlug",
              params: { slug: courseSlug, lessonSlug: prevLesson?.slug ?? "" },
            }}
          />

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
      </Box>
    </Stack>
  );
}

// Next, for a signed-in learner — or, on the last lesson, Finish, which leads
// to the course completion page. Unavailable until the lesson is passed, as in legacy, so the button
// tells the learner what is expected; `passed` already covers a lesson finished
// on an earlier visit, because revisiting is not the same as being stuck.
function ForwardButton({
  view,
  nextLesson,
  passed,
}: {
  view: CourseLessonView;
  nextLesson: LessonRef | undefined;
  passed: boolean;
}) {
  const { t } = useTranslation();
  const courseSlug = view.lesson.course.slug;
  // Next starts the following lesson and only then goes there.
  const { enter, isPending } = useEnterLesson(courseSlug);

  // The last lesson reads as completion and leads to the completion page, as in
  // legacy. Passing it finished the Enrollment that page checks for, and the
  // check dropped the cached course read so the page sees it (LessonWorkspace).
  if (!nextLesson) {
    return (
      <StepLink
        enabled={passed}
        label={t(($) => $.courses.lessons.show.finish)}
        linkOptions={
          view.landingPage
            ? { to: "/{-$locale}/languages/$slug/success", params: { slug: view.landingPage.slug } }
            : { to: "/{-$locale}/languages" }
        }
      />
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

  return (
    <StepLink
      enabled={passed}
      label={t(($) => $.courses.lessons.show.next)}
      linkOptions={{ to: "/{-$locale}/users/new", search: { redirect } }}
    />
  );
}

// A step through the course, in the controls' one style: a link while the step
// is open to the learner, and a disabled button in its place while it is not —
// a button because an anchor to nowhere would still be focusable and followable.
//
// The link options are checked against the route tree the way TanStack Router
// documents for a component that wraps a link.
interface StepLinkProps<TRouter extends RegisteredRouter = RegisteredRouter, TOptions = unknown> {
  enabled: boolean;
  label: string;
  linkOptions: ValidateLinkOptions<TRouter, TOptions>;
}

function StepLink<TRouter extends RegisteredRouter, TOptions>(
  props: StepLinkProps<TRouter, TOptions>,
): ReactNode;
function StepLink({ enabled, label, linkOptions }: StepLinkProps): ReactNode {
  return enabled ? (
    <ButtonLink variant="outline" color="green" {...linkOptions}>
      {label}
    </ButtonLink>
  ) : (
    <Button variant="outline" color="green" disabled>
      {label}
    </Button>
  );
}
