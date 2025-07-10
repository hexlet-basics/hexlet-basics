import { Box, type BoxProps } from '@mantine/core';
import { forwardRef } from 'react';

interface XssContentProps extends BoxProps {
  children?: string | null;
}

const XssContent = forwardRef<HTMLDivElement, XssContentProps>(
  ({ children, ...props }, ref) => (
    <Box
      ref={ref}
      // biome-ignore lint/security/noDangerouslySetInnerHtml: trusted HTML from server
      dangerouslySetInnerHTML={{ __html: children ?? '' }}
      {...props}
    />
  ),
);

XssContent.displayName = 'XssContent';

export default XssContent;
