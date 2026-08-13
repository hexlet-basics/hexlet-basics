import { Stack, Tabs } from "@mantine/core";
import { useLocalStorage } from "@mantine/hooks";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import type { CourseLessonView } from "@/client/types.gen";
import LessonControls from "@/components/lesson/LessonControls";
import LessonEditor from "@/components/lesson/LessonEditor";

// The right-hand half of the player: what the learner does, as opposed to what
// they read. It owns the buffer, because everything here works on it — the
// editor writes it, reset restores it, and the check that arrives with its own
// ticket submits it.
//
// The tabs keep every pane mounted: switching to the output and back must not
// remount the editor, which would cost the learner what they had typed.
export default function LessonWorkspace({ view }: { view: CourseLessonView }) {
  const { t } = useTranslation();
  const { lesson } = view;
  const courseSlug = lesson.course.slug;
  const starterCode = lesson.preparedCode ?? "";

  // A reload, a second tab or a return tomorrow finds the work still there.
  // Keyed per lesson and nothing else — as in legacy, the key carries no lesson
  // version, so a buffer survives the author changing the starter code.
  const [code, setCode] = useLocalStorage({
    key: `lesson-code-${courseSlug}-${lesson.slug}`,
    defaultValue: starterCode,
    // Read straight away rather than in an effect: the value is only ever
    // rendered inside the editor, which is client-only, so there is no
    // server-rendered markup for it to disagree with — and reading it late would
    // let monaco paint the starter code over work the learner already has.
    getInitialValueInEffect: false,
  });

  // Bumped to tell the editor its buffer has been replaced; see LessonEditor.
  const [resetCount, setResetCount] = useState(0);

  const reset = () => {
    setCode(starterCode);
    setResetCount((count) => count + 1);
  };

  return (
    <Stack h="100%" gap={0}>
      <Tabs
        defaultValue="editor"
        display="flex"
        style={{ flexDirection: "column", flexGrow: 1, minHeight: 0 }}
        keepMountedMode="display-none"
      >
        <Tabs.List grow>
          <Tabs.Tab value="editor">{t(($) => $.courses.lessons.show.editor)}</Tabs.Tab>
          <Tabs.Tab value="output">{t(($) => $.courses.lessons.show.output)}</Tabs.Tab>
          <Tabs.Tab value="tests">{t(($) => $.courses.lessons.show.tests)}</Tabs.Tab>
          <Tabs.Tab value="solution">{t(($) => $.courses.lessons.show.solution)}</Tabs.Tab>
        </Tabs.List>

        <Tabs.Panel value="editor" h="100%" mih={0}>
          <LessonEditor
            // Monaco is handed its buffer once, so moving to another lesson has
            // to give it a new editor rather than a new prop.
            key={lesson.slug}
            courseSlug={courseSlug}
            initialCode={code}
            onChange={setCode}
            starterCode={starterCode}
            resetCount={resetCount}
          />
        </Tabs.Panel>
      </Tabs>

      <LessonControls onReset={reset} />
    </Stack>
  );
}
