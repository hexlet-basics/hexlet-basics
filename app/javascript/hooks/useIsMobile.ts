import { useMantineTheme } from '@mantine/core';
import { useMediaQuery } from '@mantine/hooks';

export function useIsMobile() {
  const theme = useMantineTheme();
  // mobile = всё, что меньше sm
  return !useMediaQuery(`(min-width: ${theme.breakpoints.sm})`);
}
