import { useEffect, useState } from 'react';
import { getHighlighter } from '@/lib/shiki';
import XssContent from './XssContent';

type LazyCodeHighlightProps = {
  code: string;
  language?: string;
  className?: string;
};

/**
 * Универсальный компонент подсветки кода
 * с ленивой загрузкой языков через Shiki
 */
export function LazyCodeHighlight({
  code,
  language = 'plaintext',
  className,
}: LazyCodeHighlightProps) {
  const [html, setHtml] = useState<string | null>(null);

  useEffect(() => {
    let mounted = true;

    (async () => {
      const highlighter = await getHighlighter(language);
      const htmlOutput = highlighter.codeToHtml(code, {
        lang: language,
        theme: 'github-light', // можно сделать пропсом
      });
      if (mounted) setHtml(htmlOutput);
    })();

    return () => {
      mounted = false;
    };
  }, [code, language]);

  if (html === null) {
    return <code>{code}</code>;
  }

  return <XssContent className={className}>{html}</XssContent>;
}
