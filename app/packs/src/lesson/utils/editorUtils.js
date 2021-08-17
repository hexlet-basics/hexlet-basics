// @ts-check

import get from 'lodash/get';

const languageMapping = {
  racket: 'scheme',
};

const editorMapping = {
  racket: 'scheme',
  java: 'text/x-java',
  html: 'text/html',
};

const langToTabSizeMapping = {
  javascript: 2,
  ruby: 2,
  racket: 2,
  erlang: 2,
  elixir: 2,
  html: 2,
  css: 2,
  python: 4,
  java: 4,
  go: 4,
};

export const langToSpacesMapping = {
  javascript: true,
  json: true,
  jsx: true,
  ruby: true,
  yaml: true,
  'text/x-java': true,
  java: true,
  scheme: true,
  erlang: true,
  python: true,
  php: true,
  pug: true,
  'text/html': true,
  html: true,
  css: true,
};

const defaultTabSize = 4;

export const getLanguage = (language) => get(languageMapping, language, language);

export const getLanguageForEditor = (language) => get(editorMapping, language, language);

export const getTabSize = (language) => get(langToTabSizeMapping, language, defaultTabSize);

export const shouldReplaceTabsWithSpaces = (language) => get(langToSpacesMapping, language, false);
