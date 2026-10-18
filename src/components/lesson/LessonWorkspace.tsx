import { Box, Center, ScrollArea, Stack, Tabs } from "@mantine/core";
import { useHotkeys, useLocalStorage } from "@mantine/hooks";
import { notifications } from "@mantine/notifications";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { ClientOnly } from "@tanstack/react-router";
import { IconBook } from "@tabler/icons-react";
import { Suspense, lazy, useState } from "react";
import type { Dispatch, ReactNode, SetStateAction } from "react";
import { useTranslation } from "react-i18next";
import {
  checkLessonMutation,
  getCourseLessonQueryKey,
  getCourseQueryKey,
} from "@/client/@tanstack/react-query.gen";
import type { CourseLessonView, LessonCheckingResponse } from "@/client/types.gen";
import LessonControls from "@/components/lesson/LessonControls";
import LessonEditor from "@/components/lesson/LessonEditor";
import LessonOutput from "@/components/lesson/LessonOutput";
import LessonSolution from "@/components/lesson/LessonSolution";
import LessonTests from "@/components/lesson/LessonTests";
import LessonTheory from "@/components/lesson/LessonTheory";
import { hasPreview } from "@/lib/editor-languages";
import { lessonCodeKey } from "@/lib/lesson-code";

// Only the markup courses show a preview, so every other course is spared the
// Chunk rather than shipped a component it never renders.
const LessonPreview = lazy(() => import("@/components/lesson/LessonPreview"));

// The right-hand half of the player: what the learner does, as opposed to what
// They read. It owns the buffer, because everything here works on it — the
// Editor writes it, reset restores it, and the check submits it — and it owns
// The outcome of the last run, which is what the output and the solution are
// Rendered from. The outcome is held by the page, because the assistant in the
// Other pane is asked about it too.
//
// The tabs keep every pane mounted: switching to the output and back must not
// Remount the editor, which would cost the learner what they had typed.
//
// On a phone this pane also carries the theory as a tab of its own, as legacy
// Does, so re-reading the task never means leaving the editor's side.
export default function LessonWorkspace({
  view,
  phone,
  burger,
  result,
  setResult,
}: {
  view: CourseLessonView;
  // Whether the page is laid out for a phone, one pane at a time.
  phone: boolean;
  // The pane switch at the head of the tab strip; see LessonPage.
  burger: ReactNode;
  // The outcome of the last check, held by the page for the assistant too.
  result: LessonCheckingResponse | null;
  setResult: Dispatch<SetStateAction<LessonCheckingResponse | null>>;
}) {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const { lesson } = view;
  const courseSlug = lesson.course.slug;
  const starterCode = lesson.preparedCode ?? "";

  // A reload, a second tab or a return tomorrow finds the work still there.
  const [code, setCode] = useLocalStorage({
    key: lessonCodeKey(courseSlug, lesson.slug),
    defaultValue: starterCode,
    // Read straight away rather than in an effect: the value is only ever
    // Rendered inside the editor, which is client-only, so there is no
    // Server-rendered markup for it to disagree with — and reading it late would
    // Let monaco paint the starter code over work the learner already has.
    getInitialValueInEffect: false,
  });

  // Bumped to tell the editor its buffer has been replaced; see LessonEditor.
  const [resetCount, setResetCount] = useState(0);
  const [tab, setTab] = useState<string | null>("editor");

  const reset = () => {
    setCode(starterCode);
    setResetCount((count) => count + 1);
  };

  const check = useMutation({
    ...checkLessonMutation(),
    onSuccess: async (outcome) => {
      setResult(outcome);
      // The server has just recorded progress; the marks in the lesson list and
      // Everything else the payload carries are re-read rather than guessed at
      // Here.
      await queryClient.invalidateQueries({
        queryKey: getCourseLessonQueryKey({ path: { courseSlug, slug: lesson.slug } }),
      });
      // A pass also moves the course's own progress, and on the last lesson it
      // Finishes the Enrollment the completion page checks for. The Course page
      // Is read by its landing page's slug, so that is the key.
      if (outcome.passed && view.landingPage) {
        await queryClient.invalidateQueries({
          queryKey: getCourseQueryKey({ path: { slug: view.landingPage.slug } }),
        });
      }
    },
    // A network failure must never read as a wrong answer, so it says what it
    // Was and leaves the page where it stood.
    onError: () => notifications.show({ message: t(($) => $.common.errors.network) }),
  });

  const run = () => {
    if (check.isPending) {
      return;
    }

    // The contract refuses an empty submission, and the starter code of some
    // Courses is empty: pressing run on it would come back as a 400, which the
    // Error handler below reads as a network failure.
    if (code.trim() === "") {
      notifications.show({ message: t(($) => $.courses.lessons.show.empty_solution) });
      return;
    }

    // Move to where the answer will appear before asking for it.
    setTab("output");
    setResult(null);
    check.mutate({
      path: { id: lesson.id },
      body: { code, versionId: lesson.versionId ?? 0 },
    });
  };

  // For a learner whose focus is outside the editor; inside it, monaco's own
  // Binding runs the check (see LessonEditor).
  useHotkeys([["mod+Enter", run]]);

  // Passing opens the reference solution, and a lesson already finished opens it
  // From the start: revisiting is not the same as being stuck.
  const finished = view.lessons.find((item) => item.slug === lesson.slug)?.finished;
  const solutionUnlocked = Boolean(finished) || Boolean(result?.passed);

  return (
    <Stack h="100%" gap={0}>
      <Tabs
        // A phone turned into a desktop loses the theory tab; the editor takes
        // Its place rather than leaving the pane blank.
        value={!phone && tab === "lesson" ? "editor" : tab}
        onChange={setTab}
        display="flex"
        flex={1}
        mih={0}
        style={{ flexDirection: "column" }}
        keepMountedMode="display-none"
      >
        <Tabs.List grow>
          {burger}
          {phone && (
            <Tabs.Tab value="lesson" aria-label={t(($) => $.courses.lessons.show.lesson)}>
              <Center>
                <IconBook size={14} />
              </Center>
            </Tabs.Tab>
          )}
          <Tabs.Tab value="editor">{t(($) => $.courses.lessons.show.editor)}</Tabs.Tab>
          <Tabs.Tab value="output">{t(($) => $.courses.lessons.show.output)}</Tabs.Tab>
          <Tabs.Tab value="tests">{t(($) => $.courses.lessons.show.tests)}</Tabs.Tab>
          <Tabs.Tab value="solution">{t(($) => $.courses.lessons.show.solution)}</Tabs.Tab>
        </Tabs.List>

        {/* Only on a phone: the theory then renders twice, once per pane, as it
            does in legacy — and a desktop is spared a second, hidden copy. */}
        {phone && (
          <Tabs.Panel value="lesson" h="100%" mih={0}>
            <ScrollArea h="100%">
              <LessonTheory view={view} />
            </ScrollArea>
          </Tabs.Panel>
        )}

        <Tabs.Panel value="editor" h="100%" mih={0}>
          <Stack h="100%" gap={0}>
            <Box flex={1} mih={0}>
              <LessonEditor
                // Monaco is handed its buffer once, so moving to another lesson
                // Has to give it a new editor rather than a new prop.
                key={lesson.slug}
                courseSlug={courseSlug}
                extension={lesson.exerciseExtension}
                initialCode={code}
                onChange={setCode}
                starterCode={starterCode}
                resetCount={resetCount}
                onRun={run}
              />
            </Box>

            {hasPreview(courseSlug) && (
              // Client-only for the same reason the editor is: the buffer is read
              // From storage on the first client render, so markup the server
              // Rendered from the starter code would disagree with it.
              <ClientOnly>
                <Suspense>
                  <Box flex="none">
                    <LessonPreview code={code} />
                  </Box>
                </Suspense>
              </ClientOnly>
            )}
          </Stack>
        </Tabs.Panel>

        <Tabs.Panel value="output" h="100%" mih={0}>
          {/* A plain scrolling box rather than a ScrollArea: the output has to fill
              the pane, and a ScrollArea's content wrapper is only as tall as its
              content. */}
          <Box h="100%" style={{ overflow: "auto" }}>
            <LessonOutput result={result} running={check.isPending} />
          </Box>
        </Tabs.Panel>

        <Tabs.Panel value="tests" h="100%" mih={0}>
          <ScrollArea h="100%">
            <LessonTests
              courseSlug={courseSlug}
              testExtension={lesson.testExtension}
              testCode={lesson.testCode ?? ""}
            />
          </ScrollArea>
        </Tabs.Panel>

        <Tabs.Panel value="solution" h="100%" mih={0}>
          <ScrollArea h="100%">
            <LessonSolution
              courseSlug={courseSlug}
              extension={lesson.exerciseExtension}
              referenceSolution={lesson.originalCode ?? ""}
              learnerCode={code}
              unlocked={solutionUnlocked}
            />
          </ScrollArea>
        </Tabs.Panel>
      </Tabs>

      <LessonControls
        view={view}
        passed={solutionUnlocked}
        onReset={reset}
        onRun={run}
        running={check.isPending}
      />
    </Stack>
  );
}
