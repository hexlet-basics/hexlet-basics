// @ts-check

import get from 'lodash/get';

const languageMapping = {
  racket: 'scheme',
  clang: 'c',
};

const editorMapping = {
  racket: 'scheme',
  clang: 'c',
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
  clang: 2,
  csharp: 4,
  typescript: 2,
  haskell: 2,
};

export const langToSpacesMapping = {
  javascript: true,
  json: true,
  jsx: true,
  ruby: true,
  yaml: true,
  java: true,
  erlang: true,
  python: true,
  php: true,
  pug: true,
  html: true,
  css: true,
  elixir: true,
  racket: true,
  clojure: true,
  clang: true,
  csharp: true,
  typescript: true,
  haskell: true,
};

const defaultTabSize = 4;

export const getLanguage = (language) => get(languageMapping, language, language);

export const getLanguageForEditor = (language) => get(editorMapping, language, language);

export const getTabSize = (language) => get(langToTabSizeMapping, language, defaultTabSize);

export const shouldReplaceTabsWithSpaces = (language) => get(langToSpacesMapping, language, false);
