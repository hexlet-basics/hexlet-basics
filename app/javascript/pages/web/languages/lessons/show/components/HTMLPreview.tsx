import XssContent from "@/components/XssContent.tsx";
import { useEffect, useState } from "react";
import type { FrameComponentProps } from "react-frame-component";
// import Frame from "react-frame-component";

export default function HTMLPreview({ html }: { html: string }) {
  const [Frame, setFrame] =
    useState<React.ForwardRefExoticComponent<FrameComponentProps> | null>(null);

  useEffect(() => {
    import("react-frame-component").then((mod) => setFrame(() => mod.default));
  }, []);

  if (!Frame) return null;

  return (
    <div className="pt-3 px-2 pb-2 h-50 border-top">
      <Frame className="border-0 h-100 w-100">
        <XssContent>{html}</XssContent>
      </Frame>
    </div>
  );
}
