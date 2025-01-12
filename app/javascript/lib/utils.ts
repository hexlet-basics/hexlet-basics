import type { Grid, LanguageLesson } from "@/types/serializers";
import { router } from "@inertiajs/react";
import _ from "lodash";
import type {
  DataTableFilterMeta,
  DataTableStateEvent,
} from "primereact/datatable";
import { ChangeEvent, SyntheticEvent } from "react";

export function deviconClass(langName: string): string {
  const mapping: Record<string, string> = {
    css: "css3",
    html: "html5",
    cpp: "cplusplus",
    clang: "c",
    racket: "devicon",
    prolog: "devicon",
    fortran: "devicon",
  };
  const normalizedLangName = mapping[langName] ?? langName;
  return `devicon-${normalizedLangName}-plain`;
}

// export function assetPath(filepath: string) {
//   return `${import.meta.env.BASE_URL}images/${filepath}`;
// }

const editorMapping: Record<string, string> = {
  css: "html",
  racket: "scheme",
  clang: "c",
  dlang: "d",
  bash: "shell",
  "layout-designer": "html",
  "pre-course-java": "java",
  "pre-course-python": "python",
  "pre-course-javascript": "javascript",
};

const langToTabSizeMapping: Record<string, number> = {
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

export const langToSpacesMapping: Record<string, boolean> = {
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

// TODO: move to db
export const getEditorLanguage = (language: string): string =>
  editorMapping[language] ?? language;

export const getTabSize = (language: string): number =>
  langToTabSizeMapping[language] ?? defaultTabSize;

export const shouldReplaceTabsWithSpaces = (language: string): boolean =>
  langToSpacesMapping[language] ?? false;

export const languages = {
  "layout-designer": "layout-designer",
  // "pre-course-java": "pre-course-java",
  // "pre-course-python": "pre-course-python",
  // "pre-course-javascript": "pre-course-javascript",
  css: "css",
  html: "html",
  javascript: "javascript",
  php: "php",
  python: "python",
  java: "java",
  scheme: "scheme",
  ruby: "ruby",
  go: "go",
  elixir: "elixir",
  clojure: "clojure",
  clang: "clang",
  dlang: "dlang",
  lua: "lua",
  prolog: "prolog",
  haskell: "haskell",
  cpp: "cpp",
  bash: "bash",
  fortran: "fortran",
  kotlin: "kotlin",
  swift: "swift",
};

// TODO: move to db
export const neededPreview = (language: string) => {
  switch (language) {
    case languages.css:
    case languages.html:
    case languages["layout-designer"]:
      return true;
    default:
      return false;
  }
};

// const languagesToHelpByTutorUrls = {
//   java: "https://pre.hexlet.io/java?promo_name=base-java&promo_position=article-body&promo_type=link&promo_creative=test-failed-link",
//   python:
//     "https://pre.hexlet.io/python?promo_name=base-python&promo_position=article-body&promo_type=link&promo_creative=test-failed-link",
//   javascript:
//     "https://pre.hexlet.io/frontend?promo_name=base-frontend&promo_position=article-body&promo_type=link&promo_creative=test-failed-link",
//   php: "https://ru.hexlet.io/programs/php?promo_name=prof-php&promo_position=article-body&promo_type=link&promo_creative=test-failed-link",
//   css: "https://ru.hexlet.io/programs/frontend?promo_name=prof-frontend&promo_position=article-body&promo_type=link&promo_creative=test-failed-link",
//   html: "https://ru.hexlet.io/programs/frontend?promo_name=prof-frontend&promo_position=article-body&promo_type=link&promo_creative=test-failed-link",
// };
//
// export const getHelpByTutorUrl = (language) => {
//   const defaultUrl =
//     "https://premium.hexlet.io/?promo_name=premium&promo_position=body&promo_type=link&promo_creative=test-failed-link";
//
//   return languagesToHelpByTutorUrls[language] || defaultUrl;
// };

export const getKeyForStoringLessonCode = (lesson: LanguageLesson): string => {
  return `lesson-${lesson.id}`;
};

export function url(
  options: { withQuery?: boolean; onlyPath?: boolean } = {
    withQuery: false,
    onlyPath: false,
  },
): string {
  const parts = [];
  if (!options.onlyPath) {
    parts.push(window.location.origin);
  }
  parts.push(window.location.pathname);
  if (options.withQuery) {
    parts.push(window.location.href);
  }

  return parts.join("");
}

export function fieldsToFilters(
  fields: Record<string, string | number | string[] | undefined | null>,
): DataTableFilterMeta | undefined {
  if (!fields) {
    return;
  }
  const pairs = Object.entries(fields).map(([fieldWithMatcher, value]) => {
    return [fieldWithMatcher.split("_")[0], { value, matchMode: "contains" }];
  });
  const result = Object.fromEntries(pairs);
  return result;
}
