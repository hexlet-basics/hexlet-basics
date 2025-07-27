import { CodeHighlight } from '@mantine/code-highlight';
import {
  Anchor,
  Center,
  Container,
  createTheme,
  type MantineFontSize,
  MantineProvider,
  type MantineTheme,
  Stack,
  type StyleProp,
  Text,
  type TextProps,
  Title,
  type TitleProps,
} from '@mantine/core';
import { ModalsProvider } from '@mantine/modals';
import * as Sentry from '@sentry/react';
import type { PropsWithChildren } from 'react';
import { useTranslation } from 'react-i18next';
import type { TitleHeader } from '@/lib/mantine';

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

function responsiveClamp(
  fontSize: StyleProp<
    MantineFontSize | `h${1 | 2 | 3 | 4 | 5 | 6}` | number | (string & {})
  >,
): string {
  return `clamp(${fontSize} * 0.6, 1rem + 2vw, ${fontSize})`;
}

const theme = createTheme({
  components: {
    Title: Title.extend({
      vars: (
        theme: MantineTheme,
        params: TitleProps & { responsive?: boolean },
      ) => {
        if (params.responsive) {
          const order = params.order || 1;
          const key = (params.size || `h${order}`) as TitleHeader;
          const baseFontSize = theme.headings.sizes[key].fontSize;
          const fz = responsiveClamp(baseFontSize);
          return {
            root: {
              '--title-fz': fz,
            },
          };
        }
        return { root: { '--title-fz': undefined } };
      },
    }),
    Text: Text.extend({
      vars: (_theme: MantineTheme, params: TextProps) => {
        if (params.responsive) {
          // `params.fz` может быть 'sm', 'md', 'lg' или конкретным rem
          const baseFontSize = '1rem'; // fallback
          const fz = params.responsive || baseFontSize;
          // if (typeof params.fz === 'string') {
          //   const fontSizeValue = theme.fontSizes[params.fz as keyof typeof theme.fontSizes];
          //   baseFontSize = fontSizeValue
          //     ? `${fontSizeValue / 16}rem`
          //     : params.fz.endsWith('rem')
          //       ? params.fz
          //       : '1rem';
          // }
          return {
            root: {
              '--text-fz': responsiveClamp(fz),
            },
          };
        }
        return { root: {} };
      },
    }),
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
