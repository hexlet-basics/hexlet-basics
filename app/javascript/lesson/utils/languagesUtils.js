// @ts-check

export const languages = {
  'layout-designer': 'layout-designer',
  'pre-course-java': 'pre-course-java',
  'pre-course-python': 'pre-course-python',
  'pre-course-javascript': 'pre-course-javascript',
  css: 'css',
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
  clang: 'clang',
  dlang: 'dlang',
  lua: 'lua',
  prolog: 'prolog',
  haskell: 'haskell',
  cpp: 'cpp',
  bash: 'bash',
  fortran: 'fortran',
  kotlin: 'kotlin',
  swift: 'swift',
};

export const neededPreview = (language) => {
  switch (language) {
    case languages.css:
    case languages.html:
    case languages['layout-designer']:
      return true;
    default:
      return false;
  }
};
