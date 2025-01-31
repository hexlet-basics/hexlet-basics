import XssContent from "@/components/XssContent.tsx";
import Frame from "react-frame-component";

export default function HTMLPreview({ html }: { html: string }) {
  return (
    <div className="pt-3 px-2 pb-2 h-50 border-top">
      <Frame className="border-0 h-100 w-100">
        <XssContent>{html}</XssContent>
      </Frame>
    </div>
  );
}
