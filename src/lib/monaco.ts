import { loader } from "@monaco-editor/react";
import * as monaco from "monaco-editor";
// oxlint-disable default -- monaco's worker entrypoints have no named export.
import editorWorker from "monaco-editor/editor/editor.worker.js?worker";
import cssWorker from "monaco-editor/language/css/css.worker.js?worker";
import htmlWorker from "monaco-editor/language/html/html.worker.js?worker";
import jsonWorker from "monaco-editor/language/json/json.worker.js?worker";
import tsWorker from "monaco-editor/language/typescript/ts.worker.js?worker";
// oxlint-enable default

// Monaco, served by this application. `@monaco-editor/react` otherwise fetches
// the editor from cdn.jsdelivr.net at runtime, which puts the one thing a
// learner cannot work without behind a third party.
//
// Importing this module has the side effects that make that true, so it must be
// imported — and awaited — before the editor component mounts. It touches
// `self`, so it is client-only: nothing in the server bundle may reach it.
//
// The `?worker` specifiers use monaco's public subpaths rather than the internal
// `monaco-editor/esm/vs/...` ones: 0.56 changed its `exports` map to
// `"./*": "./esm/vs/*.js"`, so the old paths now resolve to a doubled,
// non-existent one and Vite fails to resolve them.

// The language services monaco asks for by name. Everything else — every course
// whose language has no service of its own — gets the plain editor worker.
const workers: Record<string, new () => Worker> = {
  json: jsonWorker,
  css: cssWorker,
  scss: cssWorker,
  less: cssWorker,
  html: htmlWorker,
  handlebars: htmlWorker,
  razor: htmlWorker,
  typescript: tsWorker,
  javascript: tsWorker,
};

self.MonacoEnvironment = {
  getWorker(_, label) {
    const Worker = workers[label] ?? editorWorker;
    return new Worker();
  },
};

loader.config({ monaco });
