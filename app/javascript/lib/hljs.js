import hljs from 'highlight.js/lib/core';
import javascript from 'highlight.js/lib/languages/javascript';
import python from 'highlight.js/lib/languages/python';
import shell from 'highlight.js/lib/languages/shell';
import php from 'highlight.js/lib/languages/php';
import java from 'highlight.js/lib/languages/java';
import ruby from 'highlight.js/lib/languages/ruby';
import sql from 'highlight.js/lib/languages/sql';
import scss from 'highlight.js/lib/languages/scss';
import css from 'highlight.js/lib/languages/css';
import xml from 'highlight.js/lib/languages/xml';

const langs = {
  javascript,
  ruby,
  shell,
  php,
  java,
  sql,
  python,
  scss,
  css,
  html: xml,
  xml,
};

Object.keys(langs).forEach((key) => hljs.registerLanguage(key, langs[key]));

export default hljs;
