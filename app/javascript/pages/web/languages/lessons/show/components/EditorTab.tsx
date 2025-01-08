import { usePage } from "@inertiajs/react";
import MonacoEditor, { loader } from "@monaco-editor/react";
import { useLocalStorage } from "@rehooks/local-storage";
import { useEffect, useState } from "react";
import slice from "../slices/RootSlice.ts";
import type { Props } from "../types.ts";
import {
  getEditorLanguage,
  getKeyForStoringLessonCode,
  getTabSize,
  shouldReplaceTabsWithSpaces,
} from "@/lib/utils.ts";

import * as monaco from "monaco-editor";
import editorWorker from "monaco-editor/esm/vs/editor/editor.worker?worker";
import cssWorker from "monaco-editor/esm/vs/language/css/css.worker?worker";
import htmlWorker from "monaco-editor/esm/vs/language/html/html.worker?worker";
import jsonWorker from "monaco-editor/esm/vs/language/json/json.worker?worker";
import tsWorker from "monaco-editor/esm/vs/language/typescript/ts.worker?worker";
import type { editor } from "monaco-editor";
import { useAppDispatch, useAppSelector } from "../slices/index.ts";

self.MonacoEnvironment = {
  getWorker(_, label) {
    if (label === "json") {
      return new jsonWorker();
    }
    if (label === "css" || label === "scss" || label === "less") {
      return new cssWorker();
    }
    if (label === "html" || label === "handlebars" || label === "razor") {
      return new htmlWorker();
    }
    if (label === "typescript" || label === "javascript") {
      return new tsWorker();
    }
    return new editorWorker();
  },
};

loader.config({ monaco });

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

  const [_, setLocalStorageContent] = useLocalStorage<string>(
    getKeyForStoringLessonCode(lesson),
    lesson.prepared_code || "",
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
    const newContent = value || ""
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
