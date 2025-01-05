import type { PropsWithChildren } from "react";

export default function XssContent(props: PropsWithChildren) {
  return (
    <div
      // biome-ignore lint/security/noDangerouslySetInnerHtml: <explanation>
      dangerouslySetInnerHTML={{ __html: props.children! }}
    />
  );
}
