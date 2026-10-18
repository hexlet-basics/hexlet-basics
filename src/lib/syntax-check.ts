// Syntax errors for the languages monaco has no language service of and nothing
// else here checks.
//
// Monaco reports errors on its own for TypeScript, JavaScript, JSON, CSS and HTML,
// and other tools cover python, ruby, go and the rest. C, C++, C#, Kotlin, make and
// shell would otherwise never hear that a bracket is unclosed until the learner has
// run the check. tree-sitter (web-tree-sitter, in the browser) parses the buffer and
// the places it could not make sense of become markers.
//
// This only finds syntax errors, never type or name errors, and it is advice: the
// learner's code is checked by the server as before, whatever the parser says.
//
// Each grammar is a separate wasm asset fetched on first use, so a course pays only
// for its own language.
import type { Language, Node, Parser } from "web-tree-sitter";
import runtimeUrl from "web-tree-sitter/web-tree-sitter.wasm?url";

// Keyed by monaco language id (`shell` is bash; `make` is the course slug, which
// editor-languages passes through unchanged).
const grammars: Record<string, () => Promise<{ default: string }>> = {
  c: () => import("tree-sitter-c/tree-sitter-c.wasm?url"),
  cpp: () => import("tree-sitter-cpp/tree-sitter-cpp.wasm?url"),
  csharp: () => import("tree-sitter-c-sharp/tree-sitter-c_sharp.wasm?url"),
  kotlin: () => import("@tree-sitter-grammars/tree-sitter-kotlin/tree-sitter-kotlin.wasm?url"),
  make: () => import("tree-sitter-make/tree-sitter-make.wasm?url"),
  shell: () => import("tree-sitter-bash/tree-sitter-bash.wasm?url"),
};

// More than a handful of markers on one half-typed line is noise, and the first
// error is the one that explains the rest.
const maxProblems = 5;

export interface SyntaxProblem {
  // A parser could not read this text, or expected a token that is not there.
  kind: "unexpected" | "missing";
  // For `missing`, the token the grammar expected, e.g. `)`.
  token: string;
  // 1-based, as monaco wants them.
  startLineNumber: number;
  startColumn: number;
  endLineNumber: number;
  endColumn: number;
}

export function canCheckSyntax(language: string): boolean {
  return language in grammars;
}

interface Runtime {
  Parser: typeof Parser;
  Language: typeof Language;
}

let runtime: Promise<Runtime> | null = null;
const parsers = new Map<string, Promise<Parser>>();

async function loadRuntime(): Promise<Runtime> {
  const module = await import("web-tree-sitter");
  await module.Parser.init({ locateFile: () => runtimeUrl });
  return module;
}

async function loadParser(language: string): Promise<Parser> {
  const loadGrammar = grammars[language];
  if (!loadGrammar) {
    throw new Error(`no grammar for ${language}`);
  }

  runtime ??= loadRuntime();
  const [{ Parser: ParserClass, Language: LanguageClass }, grammar] = await Promise.all([
    runtime,
    loadGrammar(),
  ]);
  const created = new ParserClass();
  // Fetched here and handed over as bytes: given a URL, the library first sniffs
  // `process.versions.node`, which throws wherever a `process` stub lacks `versions`.
  const response = await fetch(grammar.default);
  if (!response.ok) {
    throw new Error(`grammar ${language}: ${response.status}`);
  }
  const bytes = new Uint8Array(await response.arrayBuffer());
  created.setLanguage(await LanguageClass.load(bytes));
  return created;
}

async function forgetWhenFailed(language: string, parser: Promise<Parser>): Promise<void> {
  try {
    await parser;
  } catch {
    parsers.delete(language);
  }
}

// The parser for the language, or null when it could not be loaded.
async function parserOrNull(language: string): Promise<Parser | null> {
  try {
    return await parserFor(language);
  } catch {
    return null;
  }
}

function parserFor(language: string): Promise<Parser> {
  let parser = parsers.get(language);
  if (!parser) {
    parser = loadParser(language);
    // A failed load (offline, blocked asset) must not stick: the next edit retries.
    void forgetWhenFailed(language, parser);
    parsers.set(language, parser);
  }
  return parser;
}

// The syntax errors in `code`, or none when the language has no grammar or the
// grammar could not be loaded — the editor works the same without the hints.
export async function findSyntaxProblems(language: string, code: string): Promise<SyntaxProblem[]> {
  if (!canCheckSyntax(language)) {
    return [];
  }

  const parser = await parserOrNull(language);
  if (!parser) {
    return [];
  }

  const tree = parser.parse(code);
  if (!tree) {
    return [];
  }

  try {
    const problems: SyntaxProblem[] = [];
    collect(tree.rootNode, problems);
    return problems;
  } finally {
    tree.delete();
  }
}

// Walks only the branches that contain an error, and stops at the first error node
// on each: what is inside an ERROR node is the parser's own guesswork.
function collect(node: Node, problems: SyntaxProblem[]): void {
  if (problems.length >= maxProblems || !node.hasError) {
    return;
  }

  if (node.isMissing) {
    problems.push(toProblem("missing", node));
  } else if (node.isError) {
    problems.push(toProblem("unexpected", node));
  } else {
    for (const child of node.children) {
      collect(child, problems);
    }
  }
}

function toProblem(kind: SyntaxProblem["kind"], node: Node): SyntaxProblem {
  const start = node.startPosition;
  let end = node.endPosition;

  // An ERROR node can swallow the rest of the file, which as one red block hides
  // the code under it: mark the line it starts on.
  if (end.row !== start.row) {
    end = { column: Number.MAX_SAFE_INTEGER, row: start.row };
  }

  // A missing token has no width, and monaco draws nothing for an empty range.
  const endColumn = end.column === start.column ? end.column + 1 : end.column;

  return {
    endColumn: endColumn + 1,
    endLineNumber: end.row + 1,
    kind,
    startColumn: start.column + 1,
    startLineNumber: start.row + 1,
    token: node.type,
  };
}
