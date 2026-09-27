import { CodeHighlightAdapterProvider } from "@mantine/code-highlight";
import {
  Burger,
  Center,
  Loader,
  ScrollArea,
  Splitter,
  Tabs,
  Text,
  useMantineTheme,
} from "@mantine/core";
import {
  type SplitterPaneSize,
  useDisclosure,
  useLocalStorage,
  useMediaQuery,
} from "@mantine/hooks";
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

// The desktop split before the learner drags the divider: theory, then workspace.
const DEFAULT_PANE_SIZES: [SplitterPaneSize, SplitterPaneSize] = ["40%", "60%"];

// The lesson player's shell: theory, the assistant and navigation on the left,
// the workspace on the right. Both panes stay mounted for the life of the page —
// the editor in the right pane must never remount and lose a learner's buffer.
// The assistant sits beside the editor rather than in its tabs, so the code
// stays in view while the learner asks about it.
//
// On a phone the two panes cannot sit side by side, so one fills the screen at a
// time and a burger in each pane's tab strip swaps which — legacy's mechanism.
// Only the pane widths change; nothing unmounts, so a swap costs the learner
// nothing they have typed.
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
    defaultValue: DEFAULT_PANE_SIZES,
  });

  // Read in an effect with a desktop default, for the same reason as the sizes
  // above: the server renders the desktop split, and a phone moves to its
  // one-pane layout once it has hydrated.
  const theme = useMantineTheme();
  const isDesktop = useMediaQuery(`(min-width: ${theme.breakpoints.sm})`, true, {
    getInitialValueInEffect: true,
  });
  // Whether a phone shows the theory pane rather than the workspace. The
  // workspace comes first, as in legacy: it is where the learner acts, and its
  // own theory tab keeps the reading one tap away.
  const [theoryOpened, { toggle: togglePanes }] = useDisclosure(false);

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

  const phoneSizes: SplitterPaneSize[] = theoryOpened ? ["100%", "0%"] : ["0%", "100%"];
  const sizes = isDesktop ? paneSizes : phoneSizes;
  // Storage holds whatever array was written to it, so a pane with no entry
  // falls back to the default split instead of an undefined size.
  const [theorySize = DEFAULT_PANE_SIZES[0], workspaceSize = DEFAULT_PANE_SIZES[1]] = sizes;

  const burger = (
    <PaneBurger
      opened={theoryOpened}
      onToggle={togglePanes}
      label={t(($) => $.courses.lessons.show.navigation)}
    />
  );

  return (
    <CodeHighlightAdapterProvider adapter={shikiAdapter}>
      {/* On a phone the burger drives the panes, not a drag: the handle loses
          its thumb and its line, and any size change it reports is dropped, so
          the controlled sizes stay where the burger put them. */}
      <Splitter
        h="100%"
        sizes={sizes}
        onSizeChange={(next) => {
          if (isDesktop) setPaneSizes(next);
        }}
        withHandle={isDesktop}
        lineSize={isDesktop ? undefined : 0}
      >
        {/* The pane a phone has folded away is `inert`: still mounted, but out
            of the tab order and the accessibility tree, as it is out of sight. */}
        <Splitter.Pane
          defaultSize={theorySize}
          min={isDesktop ? "25%" : "0%"}
          inert={!isDesktop && !theoryOpened}
        >
          <Tabs defaultValue="lesson" h="100%" display="flex" style={{ flexDirection: "column" }}>
            <Tabs.List grow>
              {burger}
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

        <Splitter.Pane defaultSize={workspaceSize} inert={!isDesktop && theoryOpened}>
          <LessonWorkspace
            view={data}
            phone={!isDesktop}
            burger={burger}
            result={result}
            setResult={setResult}
          />
        </Splitter.Pane>
      </Splitter>
    </CodeHighlightAdapterProvider>
  );
}

// The pane switch a phone gets in place of the divider, at the head of each
// pane's tab strip. Hidden from `sm` up by CSS rather than by the media query, so
// a desktop never paints it, not even before hydration.
function PaneBurger({
  opened,
  onToggle,
  label,
}: {
  opened: boolean;
  onToggle: () => void;
  label: string;
}) {
  return (
    <Burger
      opened={opened}
      onClick={onToggle}
      hiddenFrom="sm"
      size="sm"
      aria-label={label}
      mx="xs"
      style={{ flexGrow: 0, alignSelf: "center" }}
    />
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
