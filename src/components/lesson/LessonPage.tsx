import { CodeHighlightAdapterProvider } from "@mantine/code-highlight";
import { Center, Loader, ScrollArea, Splitter, Tabs, Text } from "@mantine/core";
import { type SplitterPaneSize, useLocalStorage } from "@mantine/hooks";
import { useQuery } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { getCourseLessonOptions } from "@/client/@tanstack/react-query.gen";
import LessonNavigation from "@/components/lesson/LessonNavigation";
import LessonTheory from "@/components/lesson/LessonTheory";
import LessonWorkspace from "@/components/lesson/LessonWorkspace";
import shikiAdapter from "@/lib/shiki";

// The lesson player's shell: theory and navigation on the left, the workspace on
// the right. Both panes stay mounted for the life of the page — the editor in
// the right pane must never remount and lose a learner's buffer.
//
// The workspace's output and reference-solution panes arrive with their own
// tickets; what is here is the editor a learner writes their solution in.
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
  // Read in an effect (the hook's default), not during render: this page is
  // server-rendered, and reading storage on the first client render would paint
  // a split the server's HTML does not have.
  const [paneSizes, setPaneSizes] = useLocalStorage<SplitterPaneSize[]>({
    key: `lesson-panes-${courseSlug}-${lessonSlug}`,
    defaultValue: ["40%", "60%"],
  });

  if (isPending) {
    return (
      <Center h="100%">
        <Loader />
      </Center>
    );
  }

  if (isError || !data) {
    return <LessonMissing />;
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
          <LessonWorkspace view={data} />
        </Splitter.Pane>
      </Splitter>
    </CodeHighlightAdapterProvider>
  );
}

// A lesson slug that resolves to nothing — a mistyped URL, or a lesson dropped
// from the course. The route renders this too, because its loader rejects before
// the page above ever runs.
export function LessonMissing() {
  const { t } = useTranslation();

  return (
    <Center h="100%">
      <Text c="red">{t(($) => $.courses.lessons.show.lesson_not_found)}</Text>
    </Center>
  );
}
