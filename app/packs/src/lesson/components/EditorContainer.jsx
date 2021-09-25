// @ts-check
import React, { useContext } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import useLocalStorage from '@rehooks/local-storage';

import MonacoEditor from './MonacoEditor.jsx';
import CodemirrorEditor from './CodemirrorEditor.jsx';
import { EDITORS } from '../slices/editorSlice';
import EntityContext from '../EntityContext.js';
import { actions } from '../slices/index.js';

const EditorContainer = () => {
  const dispatch = useDispatch();
  const editorType = useSelector((state) => state.editorSlice.editorType);
  const { language, lessonVersion } = useContext(EntityContext);
  const content = useSelector((state) => state.editorSlice.content);
  const focusesCount = useSelector((state) => state.editorSlice.focusesCount);

  const localStorageKey = `lesson-version-${lessonVersion.id}`;
  const [localStorageContent, setContent] = useLocalStorage(localStorageKey);

  const handleContentChange = (value) => {
    setContent(value);
    dispatch(actions.changeContent({ content: value }));
  };

  const handleRunCheck = (value) => {
    dispatch(actions.runCheck({ lessonVersion, editor: { content: value } }));
  };

  const props = {
    language,
    content,
    focusesCount,
    localStorageContent,
    handleContentChange,
    handleRunCheck,
  };

  let Editor;

  switch (editorType) {
    case EDITORS.monaco: {
      Editor = MonacoEditor;
      break;
    }
    case EDITORS.codemirror: {
      Editor = CodemirrorEditor;
      break;
    }
    default:
      throw new Error(`Editor type is unknown: ${editorType}`);
  }

  // eslint-disable-next-line react/jsx-props-no-spreading
  return <Editor {...props} />;
};

export default EditorContainer;
