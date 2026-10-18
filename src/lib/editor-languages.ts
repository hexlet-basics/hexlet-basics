// What the editor needs to know about a course, ported from legacy's
// `app/javascript/lib/utils.ts`.
//
// A course slug is not a monaco language id: some courses teach a language
// Monaco knows under another name (`racket` is scheme to it), and some are a
// Course about one language written inside another's syntax (`css` is written
// In html here, as legacy has it). Anything absent is passed through, which
// Covers every course whose slug already is the language.
const languageBySlug: Record<string, string> = {
  css: "html",
  racket: "scheme",
  clang: "c",
  bash: "shell",
  "layout-designer": "html",
  "pre-course-java": "java",
  "pre-course-python": "python",
  "pre-course-javascript": "javascript",
  csConsole: "csharp",
};

// What a file extension says about the language of its contents. The build's
// Spec.yml names the extension of the solution and of the test file, which is
// The truth about the code on screen; the course slug is only a guess that holds
// While a course's slug is its language (the AI course is markdown with
// JavaScript tests). Values are monaco ids. Anything absent falls back to the
// Slug-derived language.
const languageByExtension: Record<string, string> = {
  rb: "ruby",
  py: "python",
  js: "javascript",
  mjs: "javascript",
  jsx: "javascript",
  ts: "typescript",
  tsx: "typescript",
  java: "java",
  kt: "kotlin",
  go: "go",
  php: "php",
  cs: "csharp",
  c: "c",
  h: "c",
  cpp: "cpp",
  rs: "rust",
  swift: "swift",
  hs: "haskell",
  ex: "elixir",
  exs: "elixir",
  clj: "clojure",
  pl: "perl",
  lua: "lua",
  rkt: "scheme",
  scm: "scheme",
  sh: "shell",
  md: "markdown",
  html: "html",
  css: "css",
  json: "json",
};

// Indentation is a property of a language's community, not of the editor, so it
// Follows the course. Anything absent gets the default.
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
// Why this is not simply true everywhere.
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

export interface EditorSettings {
  language: string;
  tabSize: number;
  insertSpaces: boolean;
}

function languageFor(courseSlug: string, extension?: string | null): string {
  const fromExtension = extension ? languageByExtension[extension.toLowerCase()] : undefined;
  return fromExtension ?? languageBySlug[courseSlug] ?? courseSlug;
}

// The language of the test file, which can differ from the solution's. Without a
// Usable extension it is the course's language, as before.
export function getTestLanguage(courseSlug: string, testExtension?: string | null): string {
  return languageFor(courseSlug, testExtension);
}

// The three travel together — a course is opened in one language, with one
// Indentation — so they are looked up once rather than three times.
//
// `extension` is the solution file's, from the build; it wins over the slug for
// The language, while indentation stays a property of the course.
export function getEditorSettings(courseSlug: string, extension?: string | null): EditorSettings {
  return {
    language: languageFor(courseSlug, extension),
    tabSize: tabSizeBySlug[courseSlug] ?? defaultTabSize,
    insertSpaces: spacesBySlug[courseSlug] ?? false,
  };
}

// The courses whose exercise is markup rather than a program, so the learner
// Watches what they build instead of waiting for a run to show it. Legacy keeps
// The same list in code (`neededPreview`) and the payload carries no flag for
// It, so it stays a property of the course slug like everything above.
const previewedSlugs = new Set(["css", "html", "layout-designer"]);

export function hasPreview(courseSlug: string): boolean {
  return previewedSlugs.has(courseSlug);
}
