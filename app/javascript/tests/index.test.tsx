import { renderToString } from 'react-dom/server';
import { expect, test } from 'vitest';
import configure from '@/lib/configure';
import Show from '@/pages/web/blog_posts/show';
import type { BlogPost, User } from '@/types';

import '@/init.ts';
import { MantineProvider } from '@mantine/core';
import Markdown from 'react-markdown';
import MarkdownViewer from '@/components/MarkdownViewer';

configure('ru', 'ru');

test('renders to string without errors', () => {
  const content = 'JOPA';
  const vdom = (
    <MantineProvider>
      <MarkdownViewer>{content}</MarkdownViewer>
    </MantineProvider>
  );
  const html = renderToString(vdom);
  expect(html).toContain(content);
  // expect(html).toContain('Hello from SSR');
  //   expect(html).toMatchInlineSnapshot(
  //     '"<div data-reactroot=\\"\\">Hello from SSR</div>"'
  //   );
});
