import { createHighlighterCore, type LanguageRegistration } from 'shiki/core';
import { createOnigurumaEngine } from 'shiki/engine/oniguruma';

const languageLoaders: Record<string, Promise<void>> = {};

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
  cs: () => import('@shikijs/langs/csharp'),

  // Системные
  bash: () => import('@shikijs/langs/bash'),
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
};

const highlighterPromise = createHighlighterCore({
  themes: [import('@shikijs/themes/github-light')],
  langs: [],
  engine: createOnigurumaEngine(import('shiki/wasm')),
}).then((highlighter) => {
  highlighter.getTheme('github-light').bg = 'var(--mantine-color-gray-0)';
  return highlighter;
});

export async function getHighlighter(lang?: string) {
  const highlighter = await highlighterPromise;

  if (lang && !highlighter.getLoadedLanguages().includes(lang)) {
    const loader = languageModules[lang];
    if (!loader) {
      console.warn(
        `[Shiki] Unsupported language: ${lang}, fallback to plaintext.`,
      );
      return highlighter;
    }

    if (!languageLoaders[lang]) {
      languageLoaders[lang] = loader().then((langImport) =>
        highlighter.loadLanguage(langImport),
      );
    }

    await languageLoaders[lang];
  }

  return highlighter;
}
