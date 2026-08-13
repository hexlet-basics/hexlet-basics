import { CodeHighlightAdapterProvider } from "@mantine/code-highlight";
import { Center, Loader, ScrollArea, Splitter, Tabs, Text } from "@mantine/core";
import { type SplitterPaneSize, useLocalStorage } from "@mantine/hooks";
import { useQuery } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { getCourseLessonOptions } from "@/client/@tanstack/react-query.gen";
import LessonNavigation from "@/components/lesson/LessonNavigation";
import LessonTheory from "@/components/lesson/LessonTheory";
import shikiAdapter from "@/lib/shiki";

// The lesson player's shell: theory and navigation on the left, the workspace on
// the right. Both panes stay mounted for the life of the page — the editor that
// lands in the right pane must never remount and lose a learner's buffer.
//
// The right pane is its tab strip and nothing else for now; the editor, the
// output and the reference solution arrive with their own tickets.
export default function LessonPage({
  courseSlug,
  lessonSlug,
}: {
  courseSlug: string;
  lessonSlug: string;
}) {
  const { t } = useTranslation();
  const { data, isPending, isError } = useQuery(
    getCourseLessonOptions({ path: { courseSlug, slug: lessonSlug } }),
  );

  // Where the learner dragged the divider, remembered per lesson exactly as
  // legacy remembered it.
  const [paneSizes, setPaneSizes] = useLocalStorage<SplitterPaneSize[]>({
    key: `lesson-panes-${courseSlug}-${lessonSlug}`,
    defaultValue: ["40%", "60%"],
    getInitialValueInEffect: false,
  });

  if (isPending) {
    return (
      <Center h="100%">
        <Loader />
      </Center>
    );
  }

  if (isError || !data) {
    return (
      <Center h="100%">
        <Text c="red">{t(($) => $.courses.lessons.show.lesson_not_found)}</Text>
      </Center>
    );
  }

  return (
    <CodeHighlightAdapterProvider adapter={shikiAdapter}>
      <Splitter h="100%" sizes={paneSizes} onSizeChange={setPaneSizes} withHandle>
        <Splitter.Pane defaultSize={paneSizes[0]} min="25%">
          <Tabs defaultValue="lesson" h="100%" display="flex" style={{ flexDirection: "column" }}>
            <Tabs.List grow>
              <Tabs.Tab value="lesson">{t(($) => $.courses.lessons.show.lesson)}</Tabs.Tab>
              <Tabs.Tab value="navigation">{t(($) => $.courses.lessons.show.navigation)}</Tabs.Tab>
            </Tabs.List>

            <Tabs.Panel value="lesson" h="100%" mih={0}>
              <ScrollArea h="100%">
                <LessonTheory view={data} />
              </ScrollArea>
            </Tabs.Panel>

            <Tabs.Panel value="navigation" h="100%" mih={0}>
              <ScrollArea h="100%">
                <LessonNavigation view={data} />
              </ScrollArea>
            </Tabs.Panel>
          </Tabs>
        </Splitter.Pane>

        <Splitter.Pane defaultSize={paneSizes[1]}>
          <Tabs
            defaultValue="editor"
            h="100%"
            display="flex"
            style={{ flexDirection: "column" }}
            keepMountedMode="display-none"
          >
            <Tabs.List grow>
              <Tabs.Tab value="editor">{t(($) => $.courses.lessons.show.editor)}</Tabs.Tab>
              <Tabs.Tab value="output">{t(($) => $.courses.lessons.show.output)}</Tabs.Tab>
              <Tabs.Tab value="tests">{t(($) => $.courses.lessons.show.tests)}</Tabs.Tab>
              <Tabs.Tab value="solution">{t(($) => $.courses.lessons.show.solution)}</Tabs.Tab>
            </Tabs.List>
          </Tabs>
        </Splitter.Pane>
      </Splitter>
    </CodeHighlightAdapterProvider>
  );
}
