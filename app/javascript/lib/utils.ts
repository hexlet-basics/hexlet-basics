import debug from 'debug';

const log = debug('app');
export { log };

const editorMapping: Record<string, string> = {
  css: 'html',
  racket: 'scheme',
  clang: 'c',
  dlang: 'd',
  bash: 'shell',
  'layout-designer': 'html',
  'pre-course-java': 'java',
  'pre-course-python': 'python',
  'pre-course-javascript': 'javascript',
};

const langToTabSizeMapping: Record<string, number> = {
  '1c': 4,
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
  '1c': true,
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
  css: 'css',
  '1c': '1c',
  html: 'html',
  javascript: 'javascript',
  php: 'php',
  python: 'python',
  java: 'java',
  scheme: 'scheme',
  ruby: 'ruby',
  go: 'go',
  elixir: 'elixir',
  clojure: 'clojure',
  c: 'c',
  d: 'd',
  racket: 'racket',
  lua: 'lua',
  prolog: 'prolog',
  haskell: 'haskell',
  cpp: 'cpp',
  bash: 'bash',
  // fortran: "fortran",
  kotlin: 'kotlin',
  swift: 'swift',
};

// TODO: move to db
export const neededPreview = (language: string) => {
  switch (language) {
    case 'css':
    case 'html':
    case 'layout-designer':
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
//   php: "https://ru.hexlet.io/programs/php?promo_name=prof-php&promo_position=article-body&promo_type=link&promo_creative=test-failed-link&utm_source=code-basics&utm_medium=referral&utm_campaign=programs&utm_content=lesson",
//   css: "https://ru.hexlet.io/programs/frontend?promo_name=prof-frontend&promo_position=article-body&promo_type=link&promo_creative=test-failed-link&utm_source=code-basics&utm_medium=referral&utm_campaign=programs&utm_content=lesson",
//   html: "https://ru.hexlet.io/programs/frontend?promo_name=prof-frontend&promo_position=article-body&promo_type=link&promo_creative=test-failed-link&utm_source=code-basics&utm_medium=referral&utm_campaign=programs&utm_content=lesson",
// };
//
// export const getHelpByTutorUrl = (language) => {
//   const defaultUrl =
//     "https://premium.hexlet.io/?promo_name=premium&promo_position=body&promo_type=link&promo_creative=test-failed-link";
//
//   return languagesToHelpByTutorUrls[language] || defaultUrl;
// };

export function isCurrentUrl(checkingUrl: string) {
  return getCurrentUrl() === checkingUrl;
}

export function getCurrentUrl(
  options: { withQuery?: boolean; onlyPath?: boolean } = {
    withQuery: false,
    onlyPath: false,
  },
): string | undefined {
  const location = fromWindow('location');
  if (!location) {
    return;
  }

  const parts: string[] = [];
  if (!options.onlyPath) {
    parts.push(location.origin);
  }
  parts.push(location.pathname);
  if (options.withQuery) {
    parts.push(location.href);
  }

  return parts.join('');
}

export const localesByCode = {
  ru: {
    icon: '🇷🇺',
    name: 'Русский',
  },
  en: {
    icon: '🇺🇲',
    name: 'English',
  },
};

export const locales = [
  { name: 'Russian', code: 'ru' },
  { name: 'English', code: 'en' },
];

export function enumToOptions(en: Record<string, string>) {
  return Object.entries(en).map(([key, value]) => ({ id: key, name: value }));
}

export function hasObjectKey<T extends object>(
  obj: T,
  key: PropertyKey,
): key is keyof T {
  return key in obj;
}

export function fromWindow<K extends keyof Window>(
  key: K,
): Window[K] | undefined {
  if (typeof window !== 'undefined' && window[key] !== undefined) {
    return window[key];
  }
  return;
}
