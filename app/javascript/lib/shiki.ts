import { createHighlighterCore } from 'shiki/core';
import { createJavaScriptRegexEngine } from 'shiki/engine/javascript';

const shikiPromise = createHighlighterCore({
  langs: [
    import('@shikijs/langs/tsx'),
    import('@shikijs/langs/scss'),
    import('@shikijs/langs/css'),
    import('@shikijs/langs/html'),
    import('@shikijs/langs/bash'),
    import('@shikijs/langs/json'),
    import('@shikijs/langs/elixir'),
    import('@shikijs/langs/clojure'),
    import('@shikijs/langs/racket'),
    import('@shikijs/langs/1c'),
    // import('@shikijs/langs/c'),
    // import('@shikijs/langs/cpp'),
    import('@shikijs/langs/haskell'),
    import('@shikijs/langs/js'),
    import('@shikijs/langs/ts'),
    import('@shikijs/langs/jsx'),
    import('@shikijs/langs/php'),
    import('@shikijs/langs/java'),
    import('@shikijs/langs/ruby'),
    import('@shikijs/langs/csharp'),
    import('@shikijs/langs/go'),
    import('@shikijs/langs/python'),
  ],
  themes: [
    import('@shikijs/themes/github-dark'),
    import('@shikijs/themes/github-light'),
  ],
  engine: createJavaScriptRegexEngine(),
}).then((shiki) => {
  shiki.getTheme('github-light').bg = 'var(--mantine-color-gray-0)';
  shiki.getTheme('github-dark').bg = 'var(--mantine-color-gray-7)';

  return shiki
});

export async function loadShiki() {
  const shiki = await shikiPromise
  return shiki;
}
