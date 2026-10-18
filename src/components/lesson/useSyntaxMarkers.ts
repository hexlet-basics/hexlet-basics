import type * as Monaco from "monaco-editor";
import type { editor } from "monaco-editor";
import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { canCheckSyntax, findSyntaxProblems } from "@/lib/syntax-check";

const owner = "tree-sitter";

// How long the buffer has to sit still before it is parsed again: a parse per
// keystroke would flag every half-typed line the learner is in the middle of.
const settleMs = 400;

// Marks the syntax errors tree-sitter finds in the learner's code, for the languages
// monaco has no service of (see `lib/syntax-check`). Markers are advice only; nothing
// here stops the learner from running the check.
export function useSyntaxMarkers(
  instance: editor.IStandaloneCodeEditor | undefined,
  monaco: typeof Monaco | undefined,
  language: string,
): void {
  const { t } = useTranslation();

  useEffect(() => {
    const model = instance?.getModel();
    if (!instance || !monaco || !model || !canCheckSyntax(language)) {
      return () => undefined;
    }

    let timer: ReturnType<typeof setTimeout> | 0 = 0;
    // Parsing is async (the grammar loads on first use), so a slow answer for an
    // older buffer must not overwrite the markers of a newer one.
    let latest = 0;
    let disposed = false;

    const check = async (): Promise<void> => {
      latest += 1;
      const version = latest;
      const problems = await findSyntaxProblems(language, model.getValue());
      if (disposed || version !== latest) {
        return;
      }

      monaco.editor.setModelMarkers(
        model,
        owner,
        problems.map((problem) => ({
          endColumn: Math.min(problem.endColumn, model.getLineMaxColumn(problem.endLineNumber)),
          endLineNumber: problem.endLineNumber,
          message:
            problem.kind === "missing"
              ? t(($) => $.courses.lessons.show.syntax_missing, { token: problem.token })
              : t(($) => $.courses.lessons.show.syntax_error),
          severity: monaco.MarkerSeverity.Error,
          startColumn: problem.startColumn,
          startLineNumber: problem.startLineNumber,
        })),
      );
    };

    void check();
    const subscription = model.onDidChangeContent(() => {
      clearTimeout(timer);
      timer = setTimeout(() => void check(), settleMs);
    });

    return () => {
      disposed = true;
      clearTimeout(timer);
      subscription.dispose();
      monaco.editor.setModelMarkers(model, owner, []);
    };
  }, [instance, monaco, language, t]);
}
