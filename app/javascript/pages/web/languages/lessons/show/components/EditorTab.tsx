// import "@/lib/monacoLoader.ts";

import { usePage } from '@inertiajs/react';
import MonacoEditor from '@monaco-editor/react';
import type { editor } from 'monaco-editor';

import { useEffect, useState } from 'react';
import {
  getEditorLanguage,
  // getKeyForStoringLessonCode,
  getTabSize,
  shouldReplaceTabsWithSpaces,
} from '@/lib/utils.ts';
import { useLessonStore } from '../store.tsx';
import type { LessonSharedProps } from '../types.ts';

export default function EditorTab() {
  const { course, mobileBrowser } = usePage<LessonSharedProps>().props;

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
    renderWhitespace: 'trailing',
    formatOnPaste: true,
    // renderLineHighlight: false,
    fixedOverflowWidgets: true,
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

    editorInstance?.focus();
  }, [focusesCount, editorInstance, mobileBrowser]);

  const handleEditorChange = (value: string | undefined) => {
    const newContent = value || '';
    // setCode(newContent);
    changeContent(newContent);
  };

  return (
    <MonacoEditor
      theme="github-light"
      options={editorOptions}
      onMount={handleEditorDidMount}
      value={content}
      onChange={handleEditorChange}
      language={getEditorLanguage(course.slug!)}
      // defaultLanguage={course.slug!}
      // className="w-100 h-100"
    />
  );
}
