import { createShikiAdapter } from "@mantine/code-highlight";
import { createHighlighterCore } from "shiki/core";
import type { HighlighterCore, LanguageRegistration } from "shiki/core";
import { createOnigurumaEngine } from "shiki/engine/oniguruma";

// Syntax highlighting for lesson theory, ported from legacy. Grammars are
// imported lazily, one module per language, so a page that shows no code pays
// for no grammar.
const languageModules: Record<string, () => Promise<{ default: LanguageRegistration[] }>> = {
  css: () => import("@shikijs/langs/css"),
  // Web
  html: () => import("@shikijs/langs/html"),
  scss: () => import("@shikijs/langs/scss"),

  // JavaScript / TypeScript
  javascript: () => import("@shikijs/langs/js"),
  js: () => import("@shikijs/langs/js"),
  jsx: () => import("@shikijs/langs/jsx"),
  ts: () => import("@shikijs/langs/ts"),
  tsx: () => import("@shikijs/langs/tsx"),
  typescript: () => import("@shikijs/langs/ts"),

  cs: () => import("@shikijs/langs/csharp"),
  csharp: () => import("@shikijs/langs/csharp"),
  go: () => import("@shikijs/langs/go"),
  java: () => import("@shikijs/langs/java"),
  kotlin: () => import("@shikijs/langs/kotlin"),
  // Server-side
  php: () => import("@shikijs/langs/php"),
  python: () => import("@shikijs/langs/python"),
  ruby: () => import("@shikijs/langs/ruby"),

  // Systems / shell
  bash: () => import("@shikijs/langs/bash"),
  json: () => import("@shikijs/langs/json"),
  sh: () => import("@shikijs/langs/sh"),
  shell: () => import("@shikijs/langs/shell"),

  clojure: () => import("@shikijs/langs/clojure"),
  elixir: () => import("@shikijs/langs/elixir"),
  // Functional
  haskell: () => import("@shikijs/langs/haskell"),

  "1c": () => import("@shikijs/langs/1c"),
  c: () => import("@shikijs/langs/c"),
  cpp: () => import("@shikijs/langs/cpp"),
  diff: () => import("@shikijs/langs/diff"),
  lua: () => import("@shikijs/langs/lua"),
  markdown: () => import("@shikijs/langs/markdown"),
  md: () => import("@shikijs/langs/markdown"),
  // The rest
  perl: () => import("@shikijs/langs/perl"),
  prolog: () => import("@shikijs/langs/prolog"),
  racket: () => import("@shikijs/langs/racket"),
  rust: () => import("@shikijs/langs/rust"),
  scheme: () => import("@shikijs/langs/scheme"),
  swift: () => import("@shikijs/langs/swift"),
};

// The languages actually loaded into the highlighter. Anything else renders as
// plain text rather than throwing "Language `x` not found" (legacy issue #597).
export const supportedLanguages = new Set(Object.keys(languageModules));

// Shiki's built-in "language" that needs no grammar.
export const plainTextLanguage = "text";

// Likewise built in: shiki colours ANSI escape sequences itself, which is what a
// test runner's report is made of.
export const ansiLanguage = "ansi";

// An author asking for plain text explicitly is not an unknown language.
const plainTextAliases = new Set(["text", "plaintext", "plain", "txt"]);

// Resolve a language name to one the highlighter has a grammar for. Needed
// wherever the name comes from somewhere other than a fenced block — a course
// slug, say — and may have no grammar at all.
export function toSupportedLanguage(language: string): string {
  return supportedLanguages.has(language) ? language : plainTextLanguage;
}

// Classify the language tag of a fenced code block. `language` is guaranteed to
// be loaded; `unknown` says whether the tag is worth reporting as a gap.
//
// A "glued" tag like `javapublic` or `csvar` is a content bug — the fence
// language ran into the first line of code (```java + public) — not a language
// to support. Such a tag starts with a language we already know, so it falls
// back to plain text quietly. A genuinely new single token still reports as
// unknown, so a language we have yet to add is not lost in the noise.
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

function loadShiki(): Promise<HighlighterCore> {
  return createHighlighterCore({
    engine: createOnigurumaEngine(() => import("shiki/wasm")),
    langs: Object.values(languageModules),
    themes: [() => import("@shikijs/themes/github-light")],
  });
}

const shikiAdapter = createShikiAdapter(loadShiki);
export default shikiAdapter;
