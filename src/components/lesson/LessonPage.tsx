import { CodeHighlightAdapterProvider } from "@mantine/code-highlight";
import { Center, Loader, ScrollArea, Splitter, Tabs, Text } from "@mantine/core";
import { type SplitterPaneSize, useLocalStorage } from "@mantine/hooks";
import { useQuery } from "@tanstack/react-query";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { getCourseLessonOptions } from "@/client/@tanstack/react-query.gen";
import type { LessonCheckingResponse } from "@/client/types.gen";
import LessonAssistant from "@/components/lesson/LessonAssistant";
import LessonNavigation from "@/components/lesson/LessonNavigation";
import LessonTheory from "@/components/lesson/LessonTheory";
import LessonWorkspace from "@/components/lesson/LessonWorkspace";
import shikiAdapter from "@/lib/shiki";

// The lesson player's shell: theory, the assistant and navigation on the left,
// the workspace on the right. Both panes stay mounted for the life of the page —
// the editor in the right pane must never remount and lose a learner's buffer.
// The assistant sits beside the editor rather than in its tabs, so the code
// stays in view while the learner asks about it.
//
// The workspace's output and reference-solution panes arrive with their own
// tickets; what is here is the editor a learner writes their solution in.
export default function LessonPage({
  courseSlug,
  lessonSlug,
  signedIn,
}: {
  courseSlug: string;
  lessonSlug: string;
  signedIn: boolean;
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

  // The outcome of the last check: the workspace renders it, the assistant is
  // asked about it.
  const [result, setResult] = useState<LessonCheckingResponse | null>(null);

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
              <Tabs.Tab value="assistant">{t(($) => $.courses.lessons.show.assistant)}</Tabs.Tab>
              <Tabs.Tab value="navigation">{t(($) => $.courses.lessons.show.navigation)}</Tabs.Tab>
            </Tabs.List>

            <Tabs.Panel value="lesson" h="100%" mih={0}>
              <ScrollArea h="100%">
                <LessonTheory view={data} />
              </ScrollArea>
            </Tabs.Panel>

            <Tabs.Panel value="assistant" h="100%" mih={0}>
              <ScrollArea h="100%">
                <LessonAssistant view={data} output={result?.output ?? ""} signedIn={signedIn} />
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
          <LessonWorkspace view={data} result={result} setResult={setResult} />
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
