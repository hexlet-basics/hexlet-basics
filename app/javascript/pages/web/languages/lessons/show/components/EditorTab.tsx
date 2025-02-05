import {
  getEditorLanguage,
  // getKeyForStoringLessonCode,
  getTabSize,
  shouldReplaceTabsWithSpaces,
} from "@/lib/utils.ts";
import { usePage } from "@inertiajs/react";
import MonacoEditor from "@monaco-editor/react";
// import useLocalStorageState from "use-local-storage-state";

import { useEffect, useState } from "react";
import slice, { runCheck } from "../slices/RootSlice.ts";
import type { LessonSharedProps } from "../types.ts";

import type { editor } from "monaco-editor";
import { useAppDispatch, useAppSelector } from "../slices/index.ts";

type Props = {
  code: string;
  setCode: React.Dispatch<React.SetStateAction<string>>;
};

export default function EditorTab({ code, setCode }: Props) {
  const { course, lesson, mobileBrowser } = usePage<LessonSharedProps>().props;

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
  const resetsCount = useAppSelector((state) => state.resetsCount);
  const content = useAppSelector((state) => state.content);
  const dispatch = useAppDispatch();

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
    if (mobileBrowser) return;

    editorInstance?.focus();
  }, [focusesCount, editorInstance, mobileBrowser]);

  // biome-ignore lint/correctness/useExhaustiveDependencies: <explanation>
  useEffect(() => {
    editorInstance?.setValue(content);
  }, [resetsCount]);

  const handleRunCheck = () => {
    dispatch(runCheck(lesson));
  };

  const handleEditorChange = (value: string | undefined) => {
    const newContent = value || "";
    setCode(newContent);
    dispatch(slice.actions.changeContent(newContent));
  };

  return (
    <MonacoEditor
      options={editorOptions}
      onMount={handleEditorDidMount}
      defaultValue={code}
      onChange={handleEditorChange}
      language={getEditorLanguage(course.slug!)}
      // defaultLanguage={course.slug!}
      // className="w-100 h-100"
    />
  );
}
