import {
  Center,
  Container,
  MantineProvider,
  Stack,
  Text,
  Title,
} from '@mantine/core';
import { ModalsProvider } from '@mantine/modals';
import * as Sentry from '@sentry/react';
import type { PropsWithChildren } from 'react';
import { useTranslation } from 'react-i18next';
import { resolver, theme } from '@/lib/mantine';

function FallbackComponent() {
  const { t: tLayouts } = useTranslation('layouts');
  return (
    <Container>
      <Center h="100vh">
        <Stack align="center">
          <Title order={1}>{tLayouts(($) => $.web.root.fallback.header)}</Title>
          <Text>{tLayouts(($) => $.web.root.fallback.description)}</Text>
        </Stack>
      </Center>
    </Container>
  );
}

type Props = PropsWithChildren & {
  // locale: Locale;
  // suffix: string | null
};

function Root(props: Props) {
  return (
    <MantineProvider cssVariablesResolver={resolver} theme={theme}>
      <ModalsProvider>
        <Sentry.ErrorBoundary fallback={FallbackComponent} showDialog>
          {props.children}
        </Sentry.ErrorBoundary>
      </ModalsProvider>
    </MantineProvider>
  );
}

export default Sentry.withProfiler(Root);
