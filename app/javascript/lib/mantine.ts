import { CodeHighlight } from "@mantine/code-highlight";
import { generateColors } from "@mantine/colors-generator";
import {
  type CSSVariablesResolver,
  createTheme,
  DEFAULT_THEME,
  getGradient,
  mergeMantineTheme,
  type TypographyProps,
} from "@mantine/core";

const VIEWPORT_MIN_PX = 320;
const VIEWPORT_MAX_PX = 1200;
function fontSizeToPx(value: `${number}rem` | `${number}px`): number {
  const size = value.trim();

  if (size.endsWith("rem")) {
    return Number.parseFloat(size) * 16;
  }

  if (size.endsWith("px")) {
    return Number.parseFloat(size);
  }

  return Number.NaN;
}

function fluidHeading(
  minSize: `${number}rem` | `${number}px`,
  maxSize: `${number}rem` | `${number}px`,
): string {
  const minPx = fontSizeToPx(minSize);
  const maxPx = fontSizeToPx(maxSize);

  if (Number.isNaN(minPx) || Number.isNaN(maxPx) || minPx >= maxPx) {
    return maxSize;
  }

  const rangePx = VIEWPORT_MAX_PX - VIEWPORT_MIN_PX;
  const slope = (maxPx - minPx) / rangePx;
  const intercept = minPx - slope * VIEWPORT_MIN_PX;
  const preferred = `calc(${intercept.toFixed(4)}px + ${(slope * 100).toFixed(4)}vw)`;

  return `clamp(${minSize}, ${preferred}, ${maxSize})`;
}

const HEXLET_BASE = "#3B37E0";
const HEXLET_VIOLET = "#D16FFF"; // additional
const HEXLET_CYAN = "#00C2FF"; // additional

// семантические (не из brand-assets, но нужны системе)
const HEXLET_GREEN = "#2FB344";
const HEXLET_YELLOW = "#F59F00";
const HEXLET_RED = "#E03131";

const myTheme = createTheme({
  colors: {
    indigo: generateColors(HEXLET_BASE),
    violet: generateColors(HEXLET_VIOLET),
    cyan: generateColors(HEXLET_CYAN),
    green: generateColors(HEXLET_GREEN),
    yellow: generateColors(HEXLET_YELLOW),
    red: generateColors(HEXLET_RED),
  },
  primaryColor: "indigo",
  primaryShade: { light: 6, dark: 5 },
  // defaultRadius: 0,
  spacing: {
    xxl: "calc(4rem * var(--mantine-scale))",
  },
  headings: {
    fontWeight: "normal",
    sizes: {
      h1: { fontSize: fluidHeading("1.75rem", "2.5rem") }, // 28-40px
      h2: { fontSize: fluidHeading("1.5rem", "2rem") }, // 24-32px
      h3: { fontSize: fluidHeading("1.25rem", "1.75rem") }, // 20-28px
      h4: { fontSize: fluidHeading("1rem", "1.5rem") }, // 16-24px
      h5: { fontSize: fluidHeading("1rem", "1.25rem") }, // 16-20px
      h6: { fontSize: "1rem" }, // 16px
    },
  },
  fontFamily: "Arial, sans-serif",
  // lineHeights: {
  //   xs: '1.4',
  //   sm: '1.45',
  //   md: '1.5', // line-height как в Bootstrap
  //   lg: '1.6',
  //   xl: '1.65',
  // },
  components: {
    CodeHighlight: CodeHighlight.extend({
      defaultProps: {
        mb: "lg",
        // withBorder: true,
        // withExpandButton: false,
        // withCopyButton: false,
        // fz: 'sm',
        // bg: 'gray.0',
        // p: 'sm'
      },
    }),
    // Anchor: Anchor.extend({
    //   defaultProps: {
    //     c: 'light-dark(var(--mantine-color-dark-9), var(--mantine-color-gray-0))',
    //   },
    // }),
  },
});

export const theme = mergeMantineTheme(DEFAULT_THEME, myTheme);

export const resolver: CSSVariablesResolver = () => ({
  variables: {},

  // светлая тема: фон = gray-0
  light: {
    "--mantine-color-dimmed": theme.colors.gray[8],
    // '--mantine-color-body': 'var(--mantine-color-gray-0)',
    "--mantine-color-anchor": "var(--mantine-color-text)",
    "--app-color-surface": theme.colors.gray[1],
    // '--app-cta-gradient':
    //   'linear-gradient(135deg, var(--mantine-color-yellow-2), var(--mantine-color-red-2))',
    // "--app-cta-gradient": "linear-gradient(to right, rgba(46, 42, 223, 0.90), rgba(46, 42, 223, 1.00), rgba(46, 42, 223, 0.90))",
    "--app-cta-gradient": getGradient(
      { deg: 90, from: "blue", to: "cyan.5" },
      theme,
    ),
    // '--mantine-color-body': 'var(--mantine-color-gray-0)',
    // '--mantine-color-default-hover': 'var(--mantine-color-gray-1)',
  },

  // тёмная тема: фон = dark-7 (или что тебе нужно)
  dark: {
    "--mantine-color-dimmed": theme.colors.gray[5],
    // '--mantine-color-body': 'var(--mantine-color-dark-7)',
    // '--app-cta-gradient':
    //   'linear-gradient(135deg, var(--mantine-color-yellow-9), var(--mantine-color-red-9))',
    "--app-cta-gradient": getGradient(
      { deg: 90, from: "blue.9", to: "cyan.7" },
      theme,
    ),
    "--app-color-surface": theme.colors.dark[6],
    "--mantine-color-anchor": "var(--mantine-color-text)",
  },
});

export const typographyStyles: TypographyProps["styles"] = (theme) => ({
  root: {
    overflowWrap: "break-word",
    // fontSize: theme.fontSizes.lg,
    lineHeight: theme.lineHeights.lg,
    // h1: t.headings.sizes.h1,
    h2: theme.headings.sizes.h5,
    h3: theme.headings.sizes.h6,
    h4: theme.headings.sizes.h6,
    h5: theme.headings.sizes.h6,
    h6: theme.headings.sizes.h6,
  },
});
