import { Box } from "@mantine/core";
import { useEffect, useState } from "react";
import type { FrameComponentProps } from "react-frame-component";
import XssContent from "@/components/XssContent.tsx";
// import Frame from "react-frame-component";

export default function HTMLPreview({ html }: { html: string }) {
  const [Frame, setFrame] =
    useState<React.ForwardRefExoticComponent<FrameComponentProps> | null>(null);

  useEffect(() => {
    import("react-frame-component").then((mod) => setFrame(() => mod.default));
  }, []);

  if (!Frame) return null;

  return (
    <Box className="border-t border-gray-400">
      <Frame className="border-none w-full">
        <XssContent>{html}</XssContent>
      </Frame>
    </Box>
  );
}
