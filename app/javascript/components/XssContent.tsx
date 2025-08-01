import { Box, type BoxProps, Typography } from '@mantine/core';
import { forwardRef } from 'react';
import { typographyStyles } from '@/lib/mantine';

interface XssContentProps extends BoxProps {
  children?: string | null;
}

const XssContent = forwardRef<HTMLDivElement, XssContentProps>(
  ({ children, ...props }, ref) => (
    <Typography styles={typographyStyles}>
      <Box
        ref={ref}
        // biome-ignore lint/security/noDangerouslySetInnerHtml: trusted HTML from server
        dangerouslySetInnerHTML={{ __html: children ?? '' }}
        {...props}
      />
    </Typography>
  ),
);

XssContent.displayName = 'XssContent';

export default XssContent;
