import _ from 'lodash';

const languageMapping = {
  racket: 'scheme',
  css: 'html',
  html: 'text/html',
};

const langToTabSizeMapping = {
  javascript: 2,
  ruby: 2,
  racket: 2,
  erlang: 2,
  elixir: 2,
  html: 2,
};

const defaultTabSize = 4;

export const getLanguage = (language) => _.get(languageMapping, language, language);

export const getTabSize = (language) => _.get(langToTabSizeMapping, language, defaultTabSize);
