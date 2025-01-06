import { useLocalStorage } from "@rehooks/local-storage";
import React, { useEffect, useRef } from "react";
// import MonacoEditor from "react-monaco-editor";
import MonacoEditor from "@monaco-editor/react";
import { useDispatch, useSelector } from "react-redux";
import { actions } from "../slices/index.js";
import {
  getLanguageForEditor,
  getTabSize,
  shouldReplaceTabsWithSpaces,
} from "../utils/editorUtils.js";
import { usePage } from "@inertiajs/react";
import type { Props } from "../types.js";

const commonOptions = {
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
  renderLineHighlight: false,
  fixedOverflowWidgets: true,
};

function Editor() {
  const { course, lesson } = usePage<Props>().props;

  const { content, focusesCount } = useSelector((state) => state.editorSlice);
  const dispatch = useDispatch();
  const editorRef = useRef(null);

  const localStorageKey = `lesson-version-${lesson.id}`;
  const [localStorageContent, setContent] = useLocalStorage(localStorageKey);

  useEffect(() => {
    editorRef.current?.focus();
  }, [focusesCount]);

  const handleRunCheck = () => {
    dispatch(actions.runCheck({ lesson }));
  };

  const editorOptions = {
    tabSize: getTabSize(course),
    insertSpaces: shouldReplaceTabsWithSpaces(course),
  };

  const onMount = (editor, monaco) => {
    editorRef.current = editor;
    const model = editor.getModel();
    model.updateOptions(editorOptions);
    model.pushEOL(0);

    // editorRef.current.focus();

    const extraKeys = [
      {
        key: monaco.KeyMod.CtrlCmd | monaco.KeyCode.Enter,
        action: handleRunCheck,
      },
    ];

    // extraKeys.forEach(({ key, action }) => {
    //   editorRef.current.addCommand(key, action);
    // });

    // NOTE: fix typescript validation error — `'name' is deprecated.(6385)`
    monaco.languages.typescript.javascriptDefaults.setDiagnosticsOptions({
      diagnosticCodesToIgnore: [6385],
    });
  };

  const onContentChange = (newContent) => {
    setContent(newContent);
    dispatch(actions.changeContent({ content: newContent }));
  };

  return (
    <MonacoEditor
      height="90vh"
      defaultLanguage="javascript"
      defaultValue="// some comment"
    />
  );

  // return (
  // <MonacoEditor
  //   defaultValue={localStorageContent || ""}
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
