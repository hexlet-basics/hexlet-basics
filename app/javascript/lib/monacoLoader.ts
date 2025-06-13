// NOTE: https://www.npmjs.com/package/@monaco-editor/react#loader-config

import { loader } from "@monaco-editor/react";
// import { shikiToMonaco } from "@shikijs/monaco";


import * as monaco from "monaco-editor";

import editorWorker from "monaco-editor/esm/vs/editor/editor.worker?worker";
import cssWorker from "monaco-editor/esm/vs/language/css/css.worker?worker";
import htmlWorker from "monaco-editor/esm/vs/language/html/html.worker?worker";
import jsonWorker from "monaco-editor/esm/vs/language/json/json.worker?worker";
import tsWorker from "monaco-editor/esm/vs/language/typescript/ts.worker?worker";

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
