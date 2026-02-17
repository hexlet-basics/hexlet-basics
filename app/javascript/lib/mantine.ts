import { CodeHighlight } from "@mantine/code-highlight";
import { generateColors } from "@mantine/colors-generator";
import {
  Anchor,
  type CSSVariablesResolver,
  createTheme,
  DEFAULT_THEME,
  getGradient,
  mergeMantineTheme,
  type TypographyProps,
} from "@mantine/core";

const VIEWPORT_MIN_PX = 320;
const VIEWPORT_MAX_PX = 1200;
function normalizeStandardFontSize(value: string): string {
  const match = value.match(/^calc\(([^*]+)\s*\*\s*var\(--mantine-scale\)\)$/);

  if (match) {
    return match[1].trim();
  }

  return value;
}

const STANDARD_FONT_SIZE = normalizeStandardFontSize(
  DEFAULT_THEME.fontSizes.md,
);

function fontSizeToPx(fontSize: string): number | null {
  const value = fontSize.trim();

  if (value.endsWith("rem")) {
    const parsed = Number.parseFloat(value);
    return Number.isNaN(parsed) ? null : parsed * 16;
  }

  if (value.endsWith("px")) {
    const parsed = Number.parseFloat(value);
    return Number.isNaN(parsed) ? null : parsed;
  }

  if (/^[\d.]+$/.test(value)) {
    const parsed = Number.parseFloat(value);
    return Number.isNaN(parsed) ? null : parsed;
  }

  return null;
}

function responsiveClamp(fontSize: string, minSize: string): string {
  const basePx = fontSizeToPx(fontSize);
  const minPx = fontSizeToPx(minSize);

  if (basePx == null || minPx == null) {
    return fontSize;
  }

  if (minPx >= basePx) {
    return fontSize;
  }

  const rangePx = VIEWPORT_MAX_PX - VIEWPORT_MIN_PX;
  const slope = (basePx - minPx) / rangePx;
  const intercept = minPx - slope * VIEWPORT_MIN_PX;
  const preferred = `calc(${intercept.toFixed(4)}px + ${(slope * 100).toFixed(4)}vw)`;

  return `clamp(${minPx.toFixed(4)}px, ${preferred}, ${basePx.toFixed(4)}px)`;
}

function responsiveFontSize(
  value: string,
  minSize = STANDARD_FONT_SIZE,
): string {
  return responsiveClamp(value, minSize);
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
  defaultRadius: 0,
  spacing: {
    xxl: "calc(4rem * var(--mantine-scale))",
  },
  headings: {
    fontWeight: "normal",
    sizes: {
      h1: { fontSize: responsiveFontSize("2.5rem", "1.75rem") }, // 40px
      h2: { fontSize: responsiveFontSize("2rem", "1.5rem") }, // 32px
      h3: { fontSize: responsiveFontSize("1.75rem", "1.25rem") }, // 28px
      h4: { fontSize: responsiveFontSize("1.5rem") }, // 24px
      h5: { fontSize: responsiveFontSize("1.25rem") }, // 20px
      h6: { fontSize: responsiveFontSize("1rem") }, // 16px
    },
  },
  fontFamily: "Arial, sans-serif",
  fontSizes: {
    xs: responsiveFontSize("12px"),
    sm: responsiveFontSize("14px"),
    md: responsiveFontSize("16px"),
    lg: responsiveFontSize("18px"),
    xl: responsiveFontSize("20px"),
    "display-3": responsiveFontSize("3rem", "2rem"), // 48px
    "display-2": responsiveFontSize("4rem", "2.25rem"), // 56px
    "display-1": responsiveFontSize("5rem", "2.5rem"), // 64px
    h1: responsiveFontSize("2.5rem", "1.75rem"), // 40px
    h2: responsiveFontSize("2rem", "1.5rem"), // 32px
    h3: responsiveFontSize("1.75rem", "1.25rem"), // 28px
    h4: responsiveFontSize("1.5rem"), // 24px
    h5: responsiveFontSize("1.25rem"), // 20px
    h6: responsiveFontSize("1rem"), // 16px
  },
  // lineHeights: {
  //   xs: '1.4',
  //   sm: '1.45',
  //   md: '1.5', // line-height как в Bootstrap
  //   lg: '1.6',
  //   xl: '1.65',
  // },
  components: {
    CodeHighlight: CodeHighlight.extend({
      styles: {
        pre: {
          padding: 0,
        },
      },
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
    "--mantine-color-dimmed": theme.colors.gray[7],
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
    "--mantine-color-dimmed": theme.colors.gray[7],
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
