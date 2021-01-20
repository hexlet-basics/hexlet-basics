import React, { useContext, useEffect, useState } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { UnControlled as CodeMirrorEditor } from 'react-codemirror2';
import { useLocalStorage } from '@rehooks/local-storage';
import { actions } from '../slices/index.js';
import { getLanguage, getTabSize } from '../utils/editorUtils.js';
import EntityContext from '../EntityContext.js';

import 'codemirror/mode/htmlmixed/htmlmixed.js';
import 'codemirror/mode/javascript/javascript.js';
import 'codemirror/mode/css/css.js';
import 'codemirror/mode/yaml/yaml.js';
import 'codemirror/mode/shell/shell.js';
import 'codemirror/mode/jsx/jsx.js';
import 'codemirror/mode/markdown/markdown.js';
import 'codemirror/mode/ruby/ruby.js';
import 'codemirror/mode/erlang/erlang.js';
import 'codemirror/mode/python/python.js';
import 'codemirror/mode/scheme/scheme.js';
import 'codemirror/mode/php/php.js';
import 'codemirror/mode/sass/sass.js';
import 'codemirror/mode/pug/pug.js';
import 'codemirror/mode/clike/clike.js';
import 'codemirror/mode/go/go';

import 'codemirror/lib/codemirror.css';
import 'codemirror/addon/scroll/simplescrollbars.css';
import 'codemirror/addon/dialog/dialog.css';

import 'codemirror/keymap/sublime.js';

import 'codemirror/addon/edit/closebrackets.js';
import 'codemirror/addon/edit/matchtags.js';
import 'codemirror/addon/edit/matchbrackets.js';
import 'codemirror/addon/edit/closetag.js';
import 'codemirror/addon/fold/xml-fold.js';
import 'codemirror/addon/comment/comment.js';
import 'codemirror/addon/scroll/simplescrollbars.js';
import 'codemirror/addon/search/searchcursor.js';
import 'codemirror/addon/search/search.js';
import 'codemirror/addon/search/jump-to-line.js';
import 'codemirror/addon/dialog/dialog.js';

const commonOptions = {
  autoCloseBrackets: true,
  autoCloseTags: true,
  autofocus: true,
  keyMap: 'sublime',
  matchBrackets: true,
  matchTags: true,
  scrollbarStyle: 'overlay',
  lineNumbers: true,
};

const Editor = () => {
  const { language, lessonVersion } = useContext(EntityContext);
  const { content, focusesCount } = useSelector((state) => state.editorSlice);
  const dispatch = useDispatch();
  const [editor, setEditor] = useState(null);

  const localStorageKey = `lesson-version-${lessonVersion.id}`;
  const [localStorageContent, setContent] = useLocalStorage(localStorageKey);

  const onContentChange = (_editor, _data, newContent) => {
    setContent(newContent);
    dispatch(actions.changeContent({ content: newContent }));
  };

  useEffect(() => {
    editor?.focus();
  }, [editor, focusesCount]);

  const onMount = (self) => {
    setEditor(self);
    self.focus();
    self.refresh();
    // self.setOption('extraKeys', {
    //   'Ctrl-Enter': () => {
    //     if (lessonMember.id) {
    //       dispatch(actions.runCheck({ lessonVersion, editor: { content: self.getValue() } }));
    //     }
    //   },
    // });
    if (localStorageContent) {
      self.getDoc().setValue(localStorageContent);
    }
  };
  const options = {
    ...commonOptions,
    mode: getLanguage(language),
    indentUnit: getTabSize(language),
  };

  return (
    <CodeMirrorEditor
      value={content}
      options={options}
      detach
      onChange={onContentChange}
      editorDidMount={onMount}
      className="w-100 h-100"
    />
  );
};

export default Editor;
