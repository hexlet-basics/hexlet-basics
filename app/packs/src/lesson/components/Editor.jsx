/* eslint-disable no-bitwise */
// @ts-check

import React, { useContext, useEffect, useRef } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import MonacoEditor from '@monaco-editor/react';
import { useLocalStorage } from '@rehooks/local-storage';
import { actions } from '../slices/index.js';
import { getLanguageForEditor, getTabSize, shouldReplaceTabsWithSpaces } from '../utils/editorUtils.js';

import EntityContext from '../EntityContext.js';

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

const Editor = () => {
  const { language, lessonVersion } = useContext(EntityContext);
  const { content, focusesCount } = useSelector((state) => state.editorSlice);
  const dispatch = useDispatch();
  const editorRef = useRef(null);

  const localStorageKey = `lesson-version-${lessonVersion.id}`;
  const [localStorageContent, setContent] = useLocalStorage(localStorageKey);

  useEffect(() => {
    editorRef.current?.focus();
  }, [focusesCount]);

  const handleRunCheck = () => {
    dispatch(actions.runCheck({ lessonVersion }));
  };

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
    setContent(newContent);
    dispatch(actions.changeContent({ content: newContent }));
  };

  return (
    <MonacoEditor
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

export default Editor;
