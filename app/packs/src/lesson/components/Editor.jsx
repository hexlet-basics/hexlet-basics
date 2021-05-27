// @ts-check

import React, {
  useContext, useEffect, useRef,
} from 'react';
import { useSelector, useDispatch } from 'react-redux';
import CodeMirror from '@uiw/react-codemirror';
import { useLocalStorage } from '@rehooks/local-storage';
import { actions } from '../slices/index.js';
import { getLanguageForEditor, getTabSize } from '../utils/editorUtils.js';

import EntityContext from '../EntityContext.js';

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

  const localStorageKey = `lesson-version-${lessonVersion.id}`;
  const [, setContent] = useLocalStorage(localStorageKey);

  const editorRef = useRef();
  // NOTE https://github.com/uiwjs/react-codemirror/blob/280f4586f3a01f7a416a04cf54ab9bec551f0462/website/App.js#L190
  const getInstance = (instance) => {
    if (instance?.editor) {
      editorRef.current = instance.editor;
    }
  };

  useEffect(() => {
    editorRef.current?.focus?.();
  }, [editorRef, focusesCount]);

  const onContentChange = (editor) => {
    const value = editor.getValue();
    setContent(value);
    dispatch(actions.changeContent({ content: value }));
  };

  const replaceTab = (cm) => {
    const space = ' ';
    const spaces = space.repeat(cm.getOption('indentUnit'));
    cm.replaceSelection(spaces);
  };

  const options = {
    autofocus: true,
    ...commonOptions,
    mode: getLanguageForEditor(language),
    indentUnit: getTabSize(language),
    extraKeys: {
      Tab: replaceTab,
      'Ctrl-Enter': () => {
        dispatch(actions.runCheck({ lessonVersion, editor: content }));
      },
    },
  };

  return (
    <CodeMirror
      value={content}
      options={options}
      onChange={onContentChange}
      ref={getInstance}
    />
  );
};

export default Editor;
