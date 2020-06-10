import React from 'react';
import MonacoEditor from 'react-monaco-editor';

const LessonPage = () => {
  const vdom = (
    <>
      <MonacoEditor
        height="600"
        // language={language}
        theme="vs-dark"
        // value={code}
        // options={options}
        // onChange={onChange}
        // editorDidMount={editorDidMount}
      />
    </>
  );
  return vdom;
};

export default LessonPage;
