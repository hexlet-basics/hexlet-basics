import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import MonacoEditor from 'react-monaco-editor';
import { actions, editorSliceName } from '../slices/index.js';
import { getLanguage, getTabSize } from '../utils/editorUtils.js';

const Editor = () => {
  const { language, content } = useSelector((state) => state[editorSliceName]);
  const dispatch = useDispatch();

  const options = {
    fontSize: 14,
    automaticLayout: true,
    minimap: {
      enabled: false,
    },
    tabSize: getTabSize(language),
  };

  const onContentChange = (newContent) => {
    dispatch(actions.changeContent({ content: newContent }));
  };

  const onMount = (editor) => {
    editor.focus();
  };

  const vdom = (
    <MonacoEditor
      language={getLanguage(language)}
      value={content}
      options={options}
      onChange={onContentChange}
      editorDidMount={onMount}
    />
  );
  return vdom;
};

export default Editor;
