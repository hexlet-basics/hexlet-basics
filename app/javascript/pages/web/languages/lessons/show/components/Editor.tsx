import { usePage } from "@inertiajs/react";
import MonacoEditor, { loader, type Monaco } from "@monaco-editor/react";
import { useLocalStorage } from "@rehooks/local-storage";
import { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { actions } from "../slices/index.js";
import type { Props } from "../types.js";
import {
  getLanguageForEditor,
  getTabSize,
  shouldReplaceTabsWithSpaces,
} from "@/lib/utils.js";

import * as monaco from "monaco-editor";
import editorWorker from "monaco-editor/esm/vs/editor/editor.worker?worker";
import cssWorker from "monaco-editor/esm/vs/language/css/css.worker?worker";
import htmlWorker from "monaco-editor/esm/vs/language/html/html.worker?worker";
import jsonWorker from "monaco-editor/esm/vs/language/json/json.worker?worker";
import tsWorker from "monaco-editor/esm/vs/language/typescript/ts.worker?worker";
import type { editor } from "monaco-editor";

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

function Editor() {
  const { course, lesson } = usePage<Props>().props;

  const editorOptions: editor.IStandaloneEditorConstructionOptions = {
    tabSize: getTabSize(course.slug),
    insertSpaces: shouldReplaceTabsWithSpaces(course.slug),
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

  const { content, focusesCount } = useSelector((state) => state.editorSlice);
  const dispatch = useDispatch();

  const localStorageKey = `lesson-${lesson.id}`;
  const [localStorageContent, setContent] = useLocalStorage<string | undefined>(
    localStorageKey,
  );

  const [editorInstance, setEditorInstance] =
    useState<editor.IStandaloneCodeEditor>();

  const handleEditorDidMount = (
    editor: editor.IStandaloneCodeEditor,
    monaco: Monaco,
  ) => {
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
    dispatch(actions.runCheck({ lesson }));
  };

  const handleContentChange = (value: string | undefined) => {
    setContent(value);
    // dispatch(actions.changeContent({ content: newContent }));
  };

  return (
    <MonacoEditor
      options={editorOptions}
      onMount={handleEditorDidMount}
      defaultValue={localStorageContent || content}
      onChange={handleContentChange}
      language={getLanguageForEditor(course.slug)}
      height="90vh"
      defaultLanguage={course.slug!}
      className="w-100 h-100"
    />
  );

  // return (
  // <MonacoEditor
  //   value={content}
  //   options={commonOptions}
  //   language={getLanguageForEditor(language)}
  //   onChange={onContentChange}
  //   editorDidMount={onMount}
  //   className="w-100 h-100"
  // />
  // );
}

export default Editor;
