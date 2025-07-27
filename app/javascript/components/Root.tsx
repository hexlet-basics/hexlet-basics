import {
  CodeHighlight,
  CodeHighlightAdapterProvider,
  createShikiAdapter,
} from '@mantine/code-highlight';
import {
  Anchor,
  Center,
  Code,
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
    CodeHighlight: CodeHighlight.extend({
      defaultProps: {
        mb: 'md',
        // withBorder: true,
        // withExpandButton: false,
        // withCopyButton: false,
        // fz: 'sm',
        // bg: 'gray.0',
        // p: 'sm'
      },
    }),
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
        <Sentry.ErrorBoundary fallback={FallbackComponent} showDialog>
          {props.children}
        </Sentry.ErrorBoundary>
      </ModalsProvider>
    </MantineProvider>
  );
}

export default Sentry.withProfiler(Root);
