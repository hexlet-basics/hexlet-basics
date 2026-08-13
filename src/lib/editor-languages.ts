// What the editor needs to know about a course, ported from legacy's
// `app/javascript/lib/utils.ts`.
//
// A course slug is not a monaco language id: some courses teach a language
// monaco knows under another name (`racket` is scheme to it), and some are a
// course about one language written inside another's syntax (`css` is written
// in html here, as legacy has it). Anything absent is passed through, which
// covers every course whose slug already is the language.
const languageBySlug: Record<string, string> = {
  css: "html",
  racket: "scheme",
  clang: "c",
  dlang: "d",
  bash: "shell",
  "layout-designer": "html",
  "pre-course-java": "java",
  "pre-course-python": "python",
  "pre-course-javascript": "javascript",
  csConsole: "csharp",
};

// Indentation is a property of a language's community, not of the editor, so it
// follows the course. Anything absent gets the default.
const tabSizeBySlug: Record<string, number> = {
  "1c": 4,
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
  powershell: 2,
  dlang: 2,
};

// Where a tab means spaces. Absent means a real tab character — go and make are
// why this is not simply true everywhere.
const spacesBySlug: Record<string, boolean> = {
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
  "1c": true,
};

const defaultTabSize = 4;

export type EditorSettings = {
  language: string;
  tabSize: number;
  insertSpaces: boolean;
};

// The three travel together — a course is opened in one language, with one
// indentation — so they are looked up once rather than three times.
export function getEditorSettings(courseSlug: string): EditorSettings {
  return {
    language: languageBySlug[courseSlug] ?? courseSlug,
    tabSize: tabSizeBySlug[courseSlug] ?? defaultTabSize,
    insertSpaces: spacesBySlug[courseSlug] ?? false,
  };
}
