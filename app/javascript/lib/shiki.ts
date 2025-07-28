import { createShikiAdapter } from '@mantine/code-highlight';
import { createHighlighterCore, type LanguageRegistration } from 'shiki/core';
import { createOnigurumaEngine } from 'shiki/engine/oniguruma';

// shikiLanguages.ts
const languageModules: Record<
  string,
  () => Promise<{ default: LanguageRegistration[] }>
> = {
  // Веб
  html: () => import('@shikijs/langs/html'),
  css: () => import('@shikijs/langs/css'),
  scss: () => import('@shikijs/langs/scss'),

  // JavaScript/TypeScript
  javascript: () => import('@shikijs/langs/js'),
  js: () => import('@shikijs/langs/js'),
  typescript: () => import('@shikijs/langs/ts'),
  ts: () => import('@shikijs/langs/ts'),
  jsx: () => import('@shikijs/langs/jsx'),
  tsx: () => import('@shikijs/langs/tsx'),

  // Серверные
  php: () => import('@shikijs/langs/php'),
  ruby: () => import('@shikijs/langs/ruby'),
  python: () => import('@shikijs/langs/python'),
  go: () => import('@shikijs/langs/go'),
  java: () => import('@shikijs/langs/java'),
  kotlin: () => import('@shikijs/langs/kotlin'),
  cs: () => import('@shikijs/langs/csharp'),
  csharp: () => import('@shikijs/langs/csharp'),

  // Системные
  bash: () => import('@shikijs/langs/bash'),
  shell: () => import('@shikijs/langs/shell'),
  sh: () => import('@shikijs/langs/sh'),
  json: () => import('@shikijs/langs/json'),

  // Функциональные
  haskell: () => import('@shikijs/langs/haskell'),
  elixir: () => import('@shikijs/langs/elixir'),
  clojure: () => import('@shikijs/langs/clojure'),

  // Прочее
  perl: () => import('@shikijs/langs/perl'),
  c: () => import('@shikijs/langs/c'),
  cpp: () => import('@shikijs/langs/cpp'),
  scheme: () => import('@shikijs/langs/scheme'),
  racket: () => import('@shikijs/langs/racket'),
  prolog: () => import('@shikijs/langs/prolog'),
  '1c': () => import('@shikijs/langs/1c'),
  md: () => import('@shikijs/langs/markdown'),
};

async function loadShiki() {
  const highlighterPromise = createHighlighterCore({
    themes: [() => import('@shikijs/themes/github-light')],
    langs: Object.values(languageModules),
    engine: createOnigurumaEngine(() => import('shiki/wasm')),
  });

  return highlighterPromise;
}

const shikiAdapter = createShikiAdapter(loadShiki);
export default shikiAdapter;
