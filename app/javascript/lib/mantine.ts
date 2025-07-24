import type { MantineTheme } from '@mantine/core';

export const typographyStyles = (t: MantineTheme) => ({
  root: {
    fontSize: t.fontSizes.md,
    lineHeight: t.lineHeights.lg,
    h1: t.headings.sizes.h1,
    h2: t.headings.sizes.h2,
    h3: t.headings.sizes.h3,
    h4: t.headings.sizes.h4,
    h5: t.headings.sizes.h5,
    h6: t.headings.sizes.h6,
  },
});
