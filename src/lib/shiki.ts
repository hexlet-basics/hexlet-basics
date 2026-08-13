import { createShikiAdapter } from "@mantine/code-highlight";
import { createHighlighterCore, type LanguageRegistration } from "shiki/core";
import { createOnigurumaEngine } from "shiki/engine/oniguruma";

// Syntax highlighting for lesson theory, ported from legacy. Grammars are
// imported lazily, one module per language, so a page that shows no code pays
// for no grammar.
const languageModules: Record<string, () => Promise<{ default: LanguageRegistration[] }>> = {
  // Web
  html: () => import("@shikijs/langs/html"),
  css: () => import("@shikijs/langs/css"),
  scss: () => import("@shikijs/langs/scss"),

  // JavaScript / TypeScript
  javascript: () => import("@shikijs/langs/js"),
  js: () => import("@shikijs/langs/js"),
  typescript: () => import("@shikijs/langs/ts"),
  ts: () => import("@shikijs/langs/ts"),
  jsx: () => import("@shikijs/langs/jsx"),
  tsx: () => import("@shikijs/langs/tsx"),

  // Server-side
  php: () => import("@shikijs/langs/php"),
  ruby: () => import("@shikijs/langs/ruby"),
  python: () => import("@shikijs/langs/python"),
  go: () => import("@shikijs/langs/go"),
  java: () => import("@shikijs/langs/java"),
  kotlin: () => import("@shikijs/langs/kotlin"),
  cs: () => import("@shikijs/langs/csharp"),
  csharp: () => import("@shikijs/langs/csharp"),

  // Systems / shell
  bash: () => import("@shikijs/langs/bash"),
  shell: () => import("@shikijs/langs/shell"),
  sh: () => import("@shikijs/langs/sh"),
  json: () => import("@shikijs/langs/json"),

  // Functional
  haskell: () => import("@shikijs/langs/haskell"),
  elixir: () => import("@shikijs/langs/elixir"),
  clojure: () => import("@shikijs/langs/clojure"),

  // The rest
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

// The languages actually loaded into the highlighter. Anything else renders as
// plain text rather than throwing "Language `x` not found" (legacy issue #597).
export const supportedLanguages = new Set(Object.keys(languageModules));

// shiki's built-in "language" that needs no grammar.
export const plainTextLanguage = "text";

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

async function loadShiki() {
  return createHighlighterCore({
    themes: [() => import("@shikijs/themes/github-light")],
    langs: Object.values(languageModules),
    engine: createOnigurumaEngine(() => import("shiki/wasm")),
  });
}

const shikiAdapter = createShikiAdapter(loadShiki);
export default shikiAdapter;
