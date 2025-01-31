import {
  getEditorLanguage,
  getKeyForStoringLessonCode,
  getTabSize,
  shouldReplaceTabsWithSpaces,
} from "@/lib/utils.ts";
import { usePage } from "@inertiajs/react";
import MonacoEditor from "@monaco-editor/react";
import useLocalStorageState from "use-local-storage-state";

import { useEffect, useState } from "react";
import slice from "../slices/RootSlice.ts";
import type { Props } from "../types.ts";

import type { editor } from "monaco-editor";
import { useAppDispatch, useAppSelector } from "../slices/index.ts";


export default function EditorTab() {
  const { course, lesson } = usePage<Props>().props;

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
  };

  const focusesCount = useAppSelector((state) => state.focusesCount);
  const content = useAppSelector((state) => state.content);
  const dispatch = useAppDispatch();

  const [_, setLocalStorageContent] = useLocalStorageState<string>(
    getKeyForStoringLessonCode(lesson),
    { defaultValue: lesson.prepared_code || "" },
  );

  const [editorInstance, setEditorInstance] =
    useState<editor.IStandaloneCodeEditor>();

  const handleEditorDidMount = (editor: editor.IStandaloneCodeEditor) => {
    setEditorInstance(editor);

    const extraKeys = [
      {
        key: monaco.KeyMod.CtrlCmd | monaco.KeyCode.Enter,
        action: handleRunCheck,
      },
    ];

    for (const v of extraKeys) {
      editor.addCommand(v.key, v.action);
    }
  };

  // biome-ignore lint/correctness/useExhaustiveDependencies: <explanation>
  useEffect(() => {
    editorInstance?.focus();
  }, [focusesCount, editorInstance]);

  const handleRunCheck = () => {
    // dispatch(slice.actions.runCheck({ lesson }));
  };

  const handleEditorChange = (value: string | undefined) => {
    const newContent = value || "";
    setLocalStorageContent(newContent);
    dispatch(slice.actions.changeContent(newContent));
  };

  return (
    <MonacoEditor
      options={editorOptions}
      onMount={handleEditorDidMount}
      defaultValue={content}
      onChange={handleEditorChange}
      language={getEditorLanguage(course.slug!)}
      // defaultLanguage={course.slug!}
      // className="w-100 h-100"
    />
  );
}
