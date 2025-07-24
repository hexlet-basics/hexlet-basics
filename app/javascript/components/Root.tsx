import {
  CodeHighlightAdapterProvider,
  createShikiAdapter,
} from '@mantine/code-highlight';
import {
  Anchor,
  Center,
  Container,
  createTheme,
  MantineProvider,
  Stack,
  Text,
  Title,
} from '@mantine/core';
import { ModalsProvider } from '@mantine/modals';
import * as Sentry from '@sentry/react';
import type { PropsWithChildren } from 'react';
import { useTranslation } from 'react-i18next';
import { getHighlighter } from '@/lib/shiki';

const shikiAdapter = createShikiAdapter(getHighlighter);

function FallbackComponent() {
  const { t: tLayouts } = useTranslation('layouts');
  return (
    <Container>
      <Center h="100vh">
        <Stack align="center">
          <Title order={1}>{tLayouts('web.root.fallback.header')}</Title>
          <Text>{tLayouts('web.root.fallback.description')}</Text>
        </Stack>
      </Center>
    </Container>
  );
}

type Props = PropsWithChildren & {
  // locale: Locale;
  // suffix: string | null
};

const theme = createTheme({
  components: {
    Anchor: Anchor.extend({
      defaultProps: {
        c: 'dark',
      },
    }),
  },
});

function Root(props: Props) {
  return (
    <MantineProvider theme={theme}>
      <ModalsProvider>
        <CodeHighlightAdapterProvider adapter={shikiAdapter}>
          <Sentry.ErrorBoundary fallback={FallbackComponent} showDialog>
            {props.children}
          </Sentry.ErrorBoundary>
        </CodeHighlightAdapterProvider>
      </ModalsProvider>
    </MantineProvider>
  );
}

export default Sentry.withProfiler(Root);
