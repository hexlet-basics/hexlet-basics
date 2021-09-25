/* eslint-disable no-bitwise */
// @ts-check

import React, { useEffect, useRef } from 'react';
import Editor from '@monaco-editor/react';
import { getLanguageForEditor, getTabSize, shouldReplaceTabsWithSpaces } from '../utils/editorUtils.js';

const commonOptions = {
  fontSize: 14,
  scrollBeyondLastLine: false,
  minimap: {
    enabled: false,
  },
  hover: {
    delay: 1000,
  },
  renderWhitespace: 'trailing',
  formatOnPaste: true,
  renderLineHighlight: false,
};

const MonacoEditor = ({
  language,
  content,
  focusesCount,
  localStorageContent,
  handleContentChange,
  handleRunCheck,
}) => {
  const editorRef = useRef(null);

  useEffect(() => {
    editorRef.current?.focus();
  }, [focusesCount]);

  const editorOptions = {
    tabSize: getTabSize(language),
    insertSpaces: shouldReplaceTabsWithSpaces(language),
  };

  const onMount = (editor, monaco) => {
    editorRef.current = editor;
    const model = editor.getModel();
    model.updateOptions(editorOptions);
    model.pushEOL(0);

    editorRef.current.focus();

    const extraKeys = [
      {
        key: monaco.KeyMod.CtrlCmd | monaco.KeyCode.Enter,
        action: handleRunCheck,
      },
    ];

    extraKeys.forEach(({ key, action }) => {
      editorRef.current.addCommand(key, action);
    });
  };

  const onContentChange = (newContent) => {
    handleContentChange(newContent);
  };

  return (
    <Editor
      defaultValue={localStorageContent || ''}
      value={content}
      options={commonOptions}
      defaultLanguage={getLanguageForEditor(language)}
      onChange={onContentChange}
      onMount={onMount}
      className="w-100 h-100"
    />
  );
};

export default MonacoEditor;
