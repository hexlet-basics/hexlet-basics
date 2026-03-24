import { Box, Divider, Stack } from "@mantine/core";
import { useEffect, useState } from "react";
import type { FrameComponentProps } from "react-frame-component";
import XssContent from "@/components/XssContent.tsx";
// import Frame from "react-frame-component";

export default function HTMLPreview({ html }: { html: string }) {
  const [Frame, setFrame] =
    useState<React.ComponentType<FrameComponentProps> | null>(null);

  useEffect(() => {
    import("react-frame-component").then((mod) => setFrame(() => mod.default));
  }, []);

  if (!Frame) return null;

  return (
    <Stack gap={0}>
      <Divider color="gray.4" />
      <Box>
        <Frame frameBorder="0" width="100%">
          <XssContent>{html}</XssContent>
        </Frame>
      </Box>
    </Stack>
  );
}
