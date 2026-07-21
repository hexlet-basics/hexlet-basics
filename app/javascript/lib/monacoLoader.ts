// NOTE: используется для предотвращения дополнительной загрузки с cdn.jsdelivr.net
// NOTE: https://www.npmjs.com/package/@monaco-editor/react#loader-config

import { loader } from "@monaco-editor/react";
// import { shikiToMonaco } from "@shikijs/monaco";

import * as monaco from "monaco-editor";

// NOTE: monaco-editor 0.56 changed its package `exports` map to `"./*": "./esm/vs/*.js"`,
// so the old internal `monaco-editor/esm/vs/...` specifiers now resolve to a doubled,
// non-existent path (esm/vs/esm/vs/...) and Vite 8 / Rolldown fails to resolve them.
// Import via the public exported subpaths instead. The `?worker` query still works.
// oxlint-disable-next-line default
import editorWorker from "monaco-editor/editor/editor.worker.js?worker";
// oxlint-disable-next-line default
import cssWorker from "monaco-editor/language/css/css.worker.js?worker";
// oxlint-disable-next-line default
import htmlWorker from "monaco-editor/language/html/html.worker.js?worker";
// oxlint-disable-next-line default
import jsonWorker from "monaco-editor/language/json/json.worker.js?worker";
// oxlint-disable-next-line default
import tsWorker from "monaco-editor/language/typescript/ts.worker.js?worker";

self.MonacoEnvironment = {
  getWorker(_, label) {
    if (label === "json") {
      return new jsonWorker();
    }
    if (label === "css" || label === "scss" || label === "less") {
      return new cssWorker();
    }
    if (label === "html" || label === "handlebars" || label === "razor") {
      return new htmlWorker();
    }
    if (label === "typescript" || label === "javascript") {
      return new tsWorker();
    }
    return new editorWorker();
  },
};

// const langs = Object.keys(languages);

// const highlighter = createHighlighter({
//   themes: [lightTheme, darkTheme],
//   langs,
// });
//
// for (const name of langs) {
//   monaco.languages.register({ id: name });
// }
//
// shikiToMonaco(highlighter, monaco);

loader.config({ monaco });
