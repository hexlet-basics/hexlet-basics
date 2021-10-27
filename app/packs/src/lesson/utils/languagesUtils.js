// @ts-check

export const languages = {
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
};

export const neededPreview = (language) => {
  switch (language) {
    case languages.css:
    case languages.html:
      return true;
    default:
      return false;
  }
};
