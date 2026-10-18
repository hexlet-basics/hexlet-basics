import { Box, Stack, useComputedColorScheme } from "@mantine/core";
import { useCallbackRef, useMediaQuery } from "@mantine/hooks";
import { ClientOnly } from "@tanstack/react-router";
import type * as Monaco from "monaco-editor";
import type { editor } from "monaco-editor";
import { Suspense, lazy, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { getEditorSettings } from "@/lib/editor-languages";
import { useSyntaxMarkers } from "@/components/lesson/useSyntaxMarkers";

// Monaco, loaded in the browser and only after the theory is on screen. It
// Cannot server-render, and it is by far the heaviest thing on the page, so it
// Hangs off a lazy import behind ClientOnly: the server ships the lesson a
// Learner reads, the browser fetches the editor afterwards.
//
// The loader module runs and is awaited first, so the editor never mounts before
// Monaco is pointed at this application's own copy instead of a CDN's.
const MonacoEditor = lazy(async () => {
  await import("@/lib/monaco");
  return import("@monaco-editor/react");
});

interface LessonEditorProps {
  courseSlug: string;
  // The solution file's extension, which names the language the learner writes.
  extension: string | null;
  // The buffer as it stands when the pane mounts. Monaco owns it from then on
  // And reports changes back through `onChange`: driving `value` on every
  // Keystroke made monaco replace its model mid-word and drop characters.
  initialCode: string;
  onChange: (code: string) => void;
  // The page pushes a buffer back only when it replaces one. `resetCount` is
  // Bumped by a reset, and `starterCode` is what a reset restores.
  starterCode: string;
  resetCount: number;
  // Ctrl+Enter, bound inside monaco. It has to be taken from monaco rather than
  // Listened for on the document: monaco binds Ctrl+Enter itself, to inserting a
  // Line below, and would swallow the keystroke and add a blank line to boot.
  onRun: () => void;
}

// Where the learner writes their solution.
export default function LessonEditor(props: LessonEditorProps) {
  return (
    <ClientOnly>
      <EditorPane {...props} />
    </ClientOnly>
  );
}

function EditorPane({
  courseSlug,
  extension,
  initialCode,
  onChange,
  starterCode,
  resetCount,
  onRun,
}: LessonEditorProps) {
  const { t } = useTranslation();
  const colorScheme = useComputedColorScheme("light", { getInitialValueInEffect: false });
  // Legacy asks the server whether the browser is a mobile one; the payload here
  // Carries no such flag, and a touch screen is what the behaviour below actually
  // Turns on: taking focus opens an on-screen keyboard over half the exercise.
  // Read on the first render — this whole pane is client-only, so there is no
  // Server markup to match.
  const isTouchScreen = useMediaQuery("(pointer: coarse)", false, {
    getInitialValueInEffect: false,
  });
  const [instance, setInstance] = useState<editor.IStandaloneCodeEditor>();
  const [monacoApi, setMonacoApi] = useState<typeof Monaco>();
  const settings = getEditorSettings(courseSlug, extension);

  useSyntaxMarkers(instance, monacoApi, settings.language);

  // Monaco keeps the callback it is given for the life of the command, so the
  // Command calls through a stable ref-backed callback rather than closing over
  // A stale run.
  const run = useCallbackRef(onRun);

  const bindEditor = (mounted: editor.IStandaloneCodeEditor, monaco: typeof Monaco) => {
    mounted.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.Enter, () => run());
    setInstance(mounted);
    setMonacoApi(monaco);
  };

  // A reset replaces the buffer and hands the cursor back, because the learner's
  // Next move is to type. `resetCount` starts at 0, which is the mount, where
  // Monaco already holds the buffer it was given.
  useEffect(() => {
    if (!instance || resetCount === 0) {
      return;
    }

    instance.setValue(starterCode);
    if (!isTouchScreen) {
      instance.focus();
    }
  }, [resetCount, starterCode, instance, isTouchScreen]);

  // Taking focus on a touch screen would open an on-screen keyboard over half
  // The exercise before the learner has read it.
  useEffect(() => {
    if (isTouchScreen) {
      return;
    }

    instance?.layout();
    instance?.focus();
  }, [instance, isTouchScreen]);

  const options: editor.IStandaloneEditorConstructionOptions = {
    tabSize: settings.tabSize,
    insertSpaces: settings.insertSpaces,
    fontSize: 14,
    scrollBeyondLastLine: false,
    minimap: { enabled: false },
    hover: { delay: 500 },
    renderWhitespace: "trailing",
    formatOnPaste: true,
    fixedOverflowWidgets: true,
    // Forced on rather than left to monaco's "auto" detection, which misses NVDA
    // On Windows and leaves the exercise unreadable to it.
    accessibilitySupport: "on",
    // And read the whole exercise to a screen reader, not monaco's default first
    // Ten lines.
    accessibilityPageSize: 500,
    ariaLabel: t(($) => $.courses.lessons.show.editor_aria_label),
  };

  return (
    <Stack h="100%" gap={0}>
      <Box flex={1} mih={0}>
        <Suspense>
          <MonacoEditor
            theme={colorScheme === "dark" ? "vs-dark" : "vs"}
            language={settings.language}
            options={options}
            defaultValue={initialCode}
            onChange={(value) => onChange(value ?? "")}
            onMount={bindEditor}
          />
        </Suspense>
      </Box>
    </Stack>
  );
}
