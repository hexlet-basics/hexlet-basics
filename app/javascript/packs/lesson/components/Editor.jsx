import React, { useContext } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import MonacoEditor from 'react-monaco-editor';
import { actions } from '../slices/index.js';
import { getLanguage, getTabSize } from '../utils/editorUtils.js';
import EntityContext from '../EntityContext.js';

const Editor = () => {
  const { language } = useContext(EntityContext);
  const { content } = useSelector((state) => state.editorSlice);
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
    // TODO: add hot key for check code on ctrl+Enter
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
