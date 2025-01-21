import type { HTMLAttributes, PropsWithChildren } from "react";

export default function XssContent(props: PropsWithChildren & HTMLAttributes<HTMLDivElement>) {
  return (
    <div className={props.className}
      // biome-ignore lint/security/noDangerouslySetInnerHtml: <explanation>
      dangerouslySetInnerHTML={{ __html: props.children! }}
    />
  );
}
