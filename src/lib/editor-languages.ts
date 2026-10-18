// What the editor needs to know about a course, ported from legacy's
// `app/javascript/lib/utils.ts`.
//
// A course slug is not a monaco language id: some courses teach a language
// monaco knows under another name (`racket` is scheme to it), and some are a
// course about one language written inside another's syntax (`css` is written
// in html here, as legacy has it). Anything absent is passed through, which
// covers every course whose slug already is the language.
const languageBySlug: Record<string, string> = {
  bash: "shell",
  clang: "c",
  csConsole: "csharp",
  css: "html",
  "layout-designer": "html",
  "pre-course-java": "java",
  "pre-course-javascript": "javascript",
  "pre-course-python": "python",
  racket: "scheme",
};

// What a file extension says about the language of its contents. The build's
// spec.yml names the extension of the solution and of the test file, which is
// the truth about the code on screen; the course slug is only a guess that holds
// while a course's slug is its language (the AI course is markdown with
// JavaScript tests). Values are monaco ids. Anything absent falls back to the
// slug-derived language.
const languageByExtension: Record<string, string> = {
  c: "c",
  clj: "clojure",
  cpp: "cpp",
  cs: "csharp",
  css: "css",
  ex: "elixir",
  exs: "elixir",
  go: "go",
  h: "c",
  hs: "haskell",
  html: "html",
  java: "java",
  js: "javascript",
  json: "json",
  jsx: "javascript",
  kt: "kotlin",
  lua: "lua",
  md: "markdown",
  mjs: "javascript",
  php: "php",
  pl: "perl",
  py: "python",
  rb: "ruby",
  rkt: "scheme",
  rs: "rust",
  scm: "scheme",
  sh: "shell",
  swift: "swift",
  ts: "typescript",
  tsx: "typescript",
};

// Indentation is a property of a language's community, not of the editor, so it
// follows the course. Anything absent gets the default.
const tabSizeBySlug: Record<string, number> = {
  "1c": 4,
  bash: 2,
  clang: 2,
  cpp: 2,
  crystal: 2,
  csharp: 4,
  css: 2,
  dart: 2,
  dlang: 2,
  elixir: 2,
  erlang: 2,
  fortran: 2,
  go: 4,
  haskell: 2,
  html: 2,
  java: 4,
  javascript: 2,
  kotlin: 4,
  ocaml: 2,
  perl: 4,
  powershell: 2,
  prolog: 1,
  python: 4,
  racket: 2,
  ruby: 2,
  rust: 4,
  swift: 2,
  typescript: 2,
};

// Where a tab means spaces. Absent means a real tab character — go and make are
// why this is not simply true everywhere.
const spacesBySlug: Record<string, boolean> = {
  "1c": true,
  bash: true,
  clang: true,
  clojure: true,
  cpp: true,
  csharp: true,
  css: true,
  elixir: true,
  erlang: true,
  fortran: true,
  haskell: true,
  html: true,
  java: true,
  javascript: true,
  json: true,
  jsx: true,
  kotlin: true,
  lua: true,
  perl: true,
  php: true,
  pug: true,
  python: true,
  racket: true,
  ruby: true,
  rust: true,
  swift: true,
  typescript: true,
  yaml: true,
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
// usable extension it is the course's language, as before.
export function getTestLanguage(courseSlug: string, testExtension?: string | null): string {
  return languageFor(courseSlug, testExtension);
}

// The three travel together — a course is opened in one language, with one
// indentation — so they are looked up once rather than three times.
//
// `extension` is the solution file's, from the build; it wins over the slug for
// the language, while indentation stays a property of the course.
export function getEditorSettings(courseSlug: string, extension?: string | null): EditorSettings {
  return {
    insertSpaces: spacesBySlug[courseSlug] ?? false,
    language: languageFor(courseSlug, extension),
    tabSize: tabSizeBySlug[courseSlug] ?? defaultTabSize,
  };
}

// The courses whose exercise is markup rather than a program, so the learner
// watches what they build instead of waiting for a run to show it. Legacy keeps
// the same list in code (`neededPreview`) and the payload carries no flag for
// it, so it stays a property of the course slug like everything above.
const previewedSlugs = new Set(["css", "html", "layout-designer"]);

export function hasPreview(courseSlug: string): boolean {
  return previewedSlugs.has(courseSlug);
}
