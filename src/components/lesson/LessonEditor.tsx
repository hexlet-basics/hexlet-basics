import { Alert, Box, Stack, useComputedColorScheme } from "@mantine/core";
import { useLocalStorage, useMediaQuery } from "@mantine/hooks";
import { ClientOnly } from "@tanstack/react-router";
import type { editor } from "monaco-editor";
import { lazy, Suspense, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { getEditorSettings } from "@/lib/editor-languages";

// Monaco, loaded in the browser and only after the theory is on screen. It
// cannot server-render, and it is by far the heaviest thing on the page, so it
// hangs off a lazy import behind ClientOnly: the server ships the lesson a
// learner reads, the browser fetches the editor afterwards.
//
// The loader module runs and is awaited first, so the editor never mounts before
// monaco is pointed at this application's own copy instead of a CDN's.
const MonacoEditor = lazy(async () => {
  await import("@/lib/monaco");
  return await import("@monaco-editor/react");
});

type LessonEditorProps = {
  courseSlug: string;
  code: string;
  onChange: (code: string) => void;
  // Bumped by the page when it wants the cursor in the editor: after a reset,
  // where the learner's next move is to type. The editor also takes focus when
  // it first mounts, as legacy does.
  focusSignal: number;
};

// Where the learner writes their solution.
export default function LessonEditor(props: LessonEditorProps) {
  return (
    <ClientOnly>
      <EditorPane {...props} />
    </ClientOnly>
  );
}

function EditorPane({ courseSlug, code, onChange, focusSignal }: LessonEditorProps) {
  const { t } = useTranslation();
  const colorScheme = useComputedColorScheme("light", { getInitialValueInEffect: false });
  // Legacy asks the server whether the browser is a mobile one; the payload here
  // carries no such flag, and a touch screen is what the two behaviours below
  // actually turn on: there is no Tab key to advertise, and taking focus opens
  // an on-screen keyboard over half the exercise. Read on the first render —
  // this whole pane is client-only, so there is no server markup to match, and
  // reading it late would paint the hint and then take it away again.
  const isTouchScreen = useMediaQuery("(pointer: coarse)", false, {
    getInitialValueInEffect: false,
  });
  // Autocomplete is not a feature people find on their own, so it is advertised
  // once — and never again, which is what makes it worth reading the first time.
  const [hintDismissed, setHintDismissed] = useLocalStorage({
    key: "lesson-editor-autocomplete-hint-dismissed",
    defaultValue: false,
    getInitialValueInEffect: false,
  });

  const [instance, setInstance] = useState<editor.IStandaloneCodeEditor>();
  const settings = getEditorSettings(courseSlug);

  useEffect(() => {
    if (isTouchScreen) return;

    instance?.layout();
    instance?.focus();
  }, [focusSignal, instance, isTouchScreen]);

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
    // on Windows and leaves the exercise unreadable to it.
    accessibilitySupport: "on",
    // And read the whole exercise to a screen reader, not monaco's default first
    // ten lines.
    accessibilityPageSize: 500,
    ariaLabel: t(($) => $.courses.lessons.show.editor_aria_label),
  };

  return (
    <Stack h="100%" gap={0}>
      {!isTouchScreen && !hintDismissed && (
        <Alert
          variant="light"
          py="xs"
          radius={0}
          withCloseButton
          closeButtonLabel={t(($) => $.courses.lessons.show.hint_close)}
          onClose={() => setHintDismissed(true)}
        >
          {t(($) => $.courses.lessons.show.autocomplete_hint)}
        </Alert>
      )}

      <Box style={{ flexGrow: 1, minHeight: 0 }}>
        <Suspense>
          <MonacoEditor
            theme={colorScheme === "dark" ? "vs-dark" : "vs"}
            language={settings.language}
            options={options}
            value={code}
            onChange={(value) => onChange(value ?? "")}
            onMount={setInstance}
          />
        </Suspense>
      </Box>
    </Stack>
  );
}
