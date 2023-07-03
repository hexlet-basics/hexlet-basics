// @ts-check

export const languages = {
  'layout-designer': 'layout-designer',
  'pre-course-java': 'pre-course-java',
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
