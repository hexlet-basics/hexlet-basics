import { createShikiAdapter } from "@mantine/code-highlight";
import { createHighlighterCore, type LanguageRegistration } from "shiki/core";
import { createOnigurumaEngine } from "shiki/engine/oniguruma";

// shikiLanguages.ts
const languageModules: Record<
  string,
  () => Promise<{ default: LanguageRegistration[] }>
> = {
  // Веб
  html: () => import("@shikijs/langs/html"),
  css: () => import("@shikijs/langs/css"),
  scss: () => import("@shikijs/langs/scss"),

  // JavaScript/TypeScript
  javascript: () => import("@shikijs/langs/js"),
  js: () => import("@shikijs/langs/js"),
  typescript: () => import("@shikijs/langs/ts"),
  ts: () => import("@shikijs/langs/ts"),
  jsx: () => import("@shikijs/langs/jsx"),
  tsx: () => import("@shikijs/langs/tsx"),

  // Серверные
  php: () => import("@shikijs/langs/php"),
  ruby: () => import("@shikijs/langs/ruby"),
  python: () => import("@shikijs/langs/python"),
  go: () => import("@shikijs/langs/go"),
  java: () => import("@shikijs/langs/java"),
  kotlin: () => import("@shikijs/langs/kotlin"),
  cs: () => import("@shikijs/langs/csharp"),
  csharp: () => import("@shikijs/langs/csharp"),

  // Системные
  bash: () => import("@shikijs/langs/bash"),
  shell: () => import("@shikijs/langs/shell"),
  sh: () => import("@shikijs/langs/sh"),
  json: () => import("@shikijs/langs/json"),

  // Функциональные
  haskell: () => import("@shikijs/langs/haskell"),
  elixir: () => import("@shikijs/langs/elixir"),
  clojure: () => import("@shikijs/langs/clojure"),

  // Прочее
  perl: () => import("@shikijs/langs/perl"),
  c: () => import("@shikijs/langs/c"),
  cpp: () => import("@shikijs/langs/cpp"),
  scheme: () => import("@shikijs/langs/scheme"),
  racket: () => import("@shikijs/langs/racket"),
  prolog: () => import("@shikijs/langs/prolog"),
  "1c": () => import("@shikijs/langs/1c"),
  md: () => import("@shikijs/langs/markdown"),
  diff: () => import("@shikijs/langs/diff"),
  swift: () => import("@shikijs/langs/swift"),
  rust: () => import("@shikijs/langs/rust"),
  lua: () => import("@shikijs/langs/lua"),
};

// Языки, которые реально загружены в highlighter. Всё остальное рендерим как
// обычный текст, чтобы не падать с "Language `x` not found" (см. issue #597).
export const supportedLanguages = new Set(Object.keys(languageModules));

// Специальные "языки" shiki, которые не требуют загрузки грамматики.
export const plainTextLanguage = "text";

// Явное указание «обычного текста» — не считаем неизвестным языком и не шумим в
// Sentry.
const plainTextAliases = new Set(["text", "plaintext", "plain", "txt"]);

// Классифицируем язык из ```-блока. `language` гарантированно загружен в
// highlighter (иначе shiki падает с "Language `x` not found", issue #597), а
// `unknown` показывает, стоит ли сообщать о нём в Sentry.
//
// Отдельно гасим «склеенные» теги вроде `javapublic`/`csvar`: это баг контента,
// где язык ```-блока прилип к первой строке кода (```java + public → javapublic),
// а не язык, который нам надо поддержать. Такой тег начинается с уже известного
// языка, поэтому молча откатываемся на текст и НЕ шумим. Действительно новые
// одиночные токены (например, язык, который мы ещё не завезли) по-прежнему
// репортим, чтобы не потерять сигнал.
// Приводим язык к тому, что точно загружен в highlighter. Нужно там, где язык
// берётся не из ```-блока, а из слага курса (getEditorLanguage) и может не иметь
// грамматики (например `csConsole`, `d`) — иначе CodeHighlight падает с
// "Language `x` not found".
export function toSupportedLanguage(language: string): string {
  return supportedLanguages.has(language) ? language : plainTextLanguage;
}

export function classifyLanguage(requested: string): {
  language: string;
  unknown: boolean;
} {
  if (supportedLanguages.has(requested)) {
    return { language: requested, unknown: false };
  }
  if (plainTextAliases.has(requested)) {
    return { language: plainTextLanguage, unknown: false };
  }
  const isGluedContentTag = [...supportedLanguages].some(
    (lang) => lang.length >= 2 && requested.startsWith(lang),
  );
  return { language: plainTextLanguage, unknown: !isGluedContentTag };
}

async function loadShiki() {
  const highlighterPromise = createHighlighterCore({
    themes: [() => import("@shikijs/themes/github-light")],
    langs: Object.values(languageModules),
    engine: createOnigurumaEngine(() => import("shiki/wasm")),
  });

  return highlighterPromise;
}

const shikiAdapter = createShikiAdapter(loadShiki);
export default shikiAdapter;
