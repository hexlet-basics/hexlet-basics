// @ts-check

export const languages = {
  "layout-designer": "layout-designer",
  "pre-course-java": "pre-course-java",
  "pre-course-python": "pre-course-python",
  "pre-course-javascript": "pre-course-javascript",
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

export const neededPreview = (language) => {
  switch (language) {
    case languages.css:
    case languages.html:
    case languages["layout-designer"]:
      return true;
    default:
      return false;
  }
};

const languagesToHelpByTutorUrls = {
  java: "https://pre.hexlet.io/java?promo_name=base-java&promo_position=article-body&promo_type=link&promo_creative=test-failed-link",
  python:
    "https://pre.hexlet.io/python?promo_name=base-python&promo_position=article-body&promo_type=link&promo_creative=test-failed-link",
  javascript:
    "https://pre.hexlet.io/frontend?promo_name=base-frontend&promo_position=article-body&promo_type=link&promo_creative=test-failed-link",
  php: "https://ru.hexlet.io/programs/php?promo_name=prof-php&promo_position=article-body&promo_type=link&promo_creative=test-failed-link",
  css: "https://ru.hexlet.io/programs/frontend?promo_name=prof-frontend&promo_position=article-body&promo_type=link&promo_creative=test-failed-link",
  html: "https://ru.hexlet.io/programs/frontend?promo_name=prof-frontend&promo_position=article-body&promo_type=link&promo_creative=test-failed-link",
};

export const getHelpByTutorUrl = (language) => {
  const defaultUrl =
    "https://premium.hexlet.io/?promo_name=premium&promo_position=body&promo_type=link&promo_creative=test-failed-link";

  return languagesToHelpByTutorUrls[language] || defaultUrl;
};
