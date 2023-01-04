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
import scheme from 'highlight.js/lib/languages/scheme';
import clojure from 'highlight.js/lib/languages/clojure';
import c from 'highlight.js/lib/languages/c';
import csharp from 'highlight.js/lib/languages/csharp';
import typescript from 'highlight.js/lib/languages/typescript';
import cpp from 'highlight.js/lib/languages/cpp';
import go from 'highlight.js/lib/languages/go';
import lua from 'highlight.js/lib/languages/lua';
import prolog from 'highlight.js/lib/languages/prolog';
import swift from 'highlight.js/lib/languages/swift';
import kotlin from 'highlight.js/lib/languages/kotlin';
import elixir from 'highlight.js/lib/languages/elixir';
import bash from 'highlight.js/lib/languages/bash';
import fortran from 'highlight.js/lib/languages/fortran';
import haskell from 'highlight.js/lib/languages/haskell';
import rust from 'highlight.js/lib/languages/rust';
import perl from 'highlight.js/lib/languages/perl';
import ocaml from 'highlight.js/lib/languages/ocaml';
import powershell from 'highlight.js/lib/languages/powershell';

const langs = {
  ocaml,
  perl,
  rust,
  haskell,
  fortran,
  bash,
  elixir,
  kotlin,
  swift,
  prolog,
  lua,
  javascript,
  ruby,
  shell,
  php,
  java,
  sql,
  powershell,
  python,
  scss,
  css,
  html: xml,
  xml,
  clojure,
  scheme,
  c,
  csharp,
  typescript,
  cpp,
  go,
};

Object.keys(langs).forEach((key) => hljs.registerLanguage(key, langs[key]));

export default hljs;
