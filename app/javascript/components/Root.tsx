import {
  CodeHighlightAdapterProvider,
  createShikiAdapter,
} from '@mantine/code-highlight';
import {
  Center,
  Container,
  MantineProvider,
  type MantineProviderProps,
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

const theme: MantineProviderProps['theme'] = {
  components: {
    Anchor: {
      defaultProps: { c: 'dark' },
    },
  },
};

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
