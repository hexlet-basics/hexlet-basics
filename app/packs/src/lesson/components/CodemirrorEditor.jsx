/* eslint-disable no-bitwise */
// @ts-check

import React, { useEffect, useRef } from 'react';
import { UnControlled as CodeMirror } from 'react-codemirror2';

import { getLanguageForEditor, getTabSize, shouldReplaceTabsWithSpaces } from '../utils/editorUtils.js';

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
import 'codemirror/mode/go/go.js';
import 'codemirror-mode-elixir';

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

const CodemirrorEditor = ({
  content,
  focusesCount,
  language,
  localStorageContent,
  handleContentChange,
  handleRunCheck,
}) => {
  const editorRef = useRef(null);

  useEffect(() => {
    editorRef.current?.focus();
  }, [focusesCount]);

  const onMount = (editor) => {
    editorRef.current = editor;
    editor.focus();
    editor.refresh();
    if (localStorageContent) {
      editor.getDoc().setValue(localStorageContent);
    }
  };

  const onContentChange = (_editor, _data, newContent) => {
    handleContentChange(newContent);
  };

  const replaceTab = (cm) => {
    const space = ' ';
    const spaces = space.repeat(cm.getOption('indentUnit'));
    cm.replaceSelection(spaces);
  };

  const extraKeys = {
    'Ctrl-Enter': handleRunCheck,
    'Cmd-Enter': handleRunCheck,
  };

  const indentWithSpaces = shouldReplaceTabsWithSpaces(language);
  const indentWithTabs = !indentWithSpaces;

  if (indentWithSpaces) {
    extraKeys.Tab = replaceTab;
  }

  const options = {
    autofocus: true,
    ...commonOptions,
    mode: getLanguageForEditor(language),
    indentUnit: getTabSize(language),
    indentWithTabs,
    extraKeys,
  };

  return (
    <CodeMirror
      value={content}
      options={options}
      onChange={onContentChange}
      detach
      editorDidMount={onMount}
      className="w-100 h-100"
    />
  );
};

export default CodemirrorEditor;
