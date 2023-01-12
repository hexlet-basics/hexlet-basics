// @ts-check

const languageMapping = {
  racket: 'scheme',
  clang: 'c',
};

const editorMapping = {
  css: 'html',
  racket: 'scheme',
  clang: 'c',
  bash: 'shell',
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
  prolog: 1,
  cpp: 2,
  bash: 2,
  fortran: 2,
  kotlin: 4,
  swift: 2,
  rust: 4,
  perl: 4,
  ocaml: 2,
  crystal: 2,
  dart: 2,
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
  lua: true,
  cpp: true,
  bash: true,
  fortran: true,
  kotlin: true,
  swift: true,
  rust: true,
  perl: true,
};

const defaultTabSize = 4;

export const getLanguage = (language) => languageMapping[language] ?? language;

export const getLanguageForEditor = (language) => editorMapping[language] ?? language;

export const getTabSize = (language) => langToTabSizeMapping[language] ?? defaultTabSize;

export const shouldReplaceTabsWithSpaces = (language) => langToSpacesMapping[language] ?? false;
