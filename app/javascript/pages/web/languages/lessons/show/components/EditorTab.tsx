import { usePage } from "@inertiajs/react";
import { Box, useComputedColorScheme } from "@mantine/core";
import MonacoEditor from "@monaco-editor/react";
import type { editor } from "monaco-editor";

import { useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import useAsyncModule from "@/hooks/useAsyncModule.ts";
import {
  getEditorLanguage,
  // getKeyForStoringLessonCode,
  getTabSize,
  shouldReplaceTabsWithSpaces,
} from "@/lib/utils.ts";
import type { LessonSharedProps } from "@/types";
import { useLessonStore } from "../store.tsx";

export default function EditorTab() {
  const { course, mobileBrowser } = usePage<LessonSharedProps>().props;
  const { t } = useTranslation();
  const colorScheme = useComputedColorScheme("light", {
    getInitialValueInEffect: false,
  });
  const loadMonaco = useCallback(() => import("@/lib/monacoLoader"), []);
  const { isLoading: isMonacoLoading, error: monacoError } =
    useAsyncModule(loadMonaco);

  const editorOptions: editor.IStandaloneEditorConstructionOptions = {
    tabSize: getTabSize(course.slug!),
    insertSpaces: shouldReplaceTabsWithSpaces(course.slug!),
    fontSize: 14,
    scrollBeyondLastLine: false,
    minimap: {
      enabled: false,
    },
    hover: {
      delay: 500,
    },
    renderWhitespace: "trailing",
    formatOnPaste: true,
    // renderLineHighlight: false,
    fixedOverflowWidgets: true,
    // Force accessibility mode on so screen readers (e.g. NVDA) can read the
    // code — Monaco's "auto" detection misses NVDA on Windows (#722).
    accessibilitySupport: "on",
    // Expose the whole exercise to the screen reader, not just the default 10
    // lines, so longer solutions are fully readable.
    accessibilityPageSize: 500,
    ariaLabel: t(($) => $.languages.lessons.show.editor_aria_label),
  };

  const focusesCount = useLessonStore((state) => state.focusesCount);
  // const resetsCount = useLessonStore((state) => state.resetsCount);
  const content = useLessonStore((state) => state.content);
  const changeContent = useLessonStore((state) => state.changeContent);

  // const defaultCode = lesson.prepared_code || "";

  const [editorInstance, setEditorInstance] =
    useState<editor.IStandaloneCodeEditor>();

  const handleEditorDidMount = (editor: editor.IStandaloneCodeEditor) => {
    setEditorInstance(editor);

    // const extraKeys = [
    //   {
    //     key: monaco.KeyMod.CtrlCmd | monaco.KeyCode.Enter,
    //     action: handleRunCheck,
    //   },
    // ];
    //
    // for (const v of extraKeys) {
    //   editor.addCommand(v.key, v.action);
    // }
  };

  // biome-ignore lint/correctness/useExhaustiveDependencies: -
  useEffect(() => {
    if (mobileBrowser) return;

    editorInstance?.layout();
    editorInstance?.focus();
  }, [focusesCount, editorInstance, mobileBrowser]);

  const handleEditorChange = (value: string | undefined) => {
    const newContent = value || "";
    // setCode(newContent);
    changeContent(newContent);
  };

  const editorTheme = colorScheme === "dark" ? "vs-dark" : "vs";

  if (isMonacoLoading || monacoError) {
    return <Box h="100%" />;
  }

  return (
    <MonacoEditor
      theme={editorTheme}
      options={editorOptions}
      onMount={handleEditorDidMount}
      value={content}
      onChange={handleEditorChange}
      language={getEditorLanguage(course.slug!)}
    />
  );
}
