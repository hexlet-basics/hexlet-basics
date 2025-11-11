import { CodeHighlight } from '@mantine/code-highlight';
import {
  Anchor,
  type CSSVariablesResolver,
  createTheme,
  type MantineFontSize,
  type MantineTheme,
  type StyleProp,
  Text,
  type TextProps,
  Title,
  type TitleProps,
} from '@mantine/core';

export type TitleHeader = keyof MantineTheme['headings']['sizes'];

function responsiveClamp(
  fontSize: StyleProp<
    MantineFontSize | `h${1 | 2 | 3 | 4 | 5 | 6}` | number | (string & {})
  >,
): string {
  return `clamp(${fontSize} * 0.6, 1rem + 2vw, ${fontSize})`;
}

export const theme = createTheme({
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
      styles: {
        pre: {
          padding: 0,
        },
      },
      defaultProps: {
        mb: 'lg',
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

export const resolver: CSSVariablesResolver = () => ({
  variables: {},

  // светлая тема: фон = gray-0
  light: {
    // '--mantine-color-body': 'var(--mantine-color-gray-0)',
    '--mantine-color-anchor': 'var(--mantine-color-text)',
    // '--app-color-surface': theme.colors.gray[0],
    '--app-cta-gradient':
      'linear-gradient(135deg, var(--mantine-color-yellow-2), var(--mantine-color-red-2))',
    // '--mantine-color-body': 'var(--mantine-color-gray-0)',
    // '--mantine-color-default-hover': 'var(--mantine-color-gray-1)',
  },

  // тёмная тема: фон = dark-7 (или что тебе нужно)
  dark: {
    // '--mantine-color-body': 'var(--mantine-color-dark-7)',
    '--app-cta-gradient':
      'linear-gradient(135deg, var(--mantine-color-yellow-9), var(--mantine-color-red-9))',
    // '--app-color-surface': theme.colors.dark[7],
    '--mantine-color-anchor': 'var(--mantine-color-text)',
  },
});

export const typographyStyles = (t: MantineTheme) => ({
  root: {
    fontSize: t.fontSizes.md,
    // lineHeight: t.lineHeights.lg,
    // h1: t.headings.sizes.h1,
    // h2: t.headings.sizes.h2,
    // h3: t.headings.sizes.h3,
    // h4: t.headings.sizes.h4,
    // h5: t.headings.sizes.h5,
    // h6: t.headings.sizes.h6,
  },
});
