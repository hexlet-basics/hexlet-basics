import { type Method, shouldIntercept } from "@inertiajs/core";
import { Link, router } from "@inertiajs/react";
import {
  Anchor,
  type AnchorProps,
  Button,
  type ButtonProps,
  type ElementProps,
} from "@mantine/core";
import type { MouseEventHandler, PropsWithChildren } from "react";
import useConfirmation from "@/hooks/useConfirmation";
import { getCurrentUrl } from "@/lib/utils";

type AppAnchorProps = AnchorProps &
  Partial<ButtonProps> &
  Omit<
    ElementProps<"a", keyof AnchorProps>,
    | "onClick"
    | "onError"
    | "onProgress"
    | "onSuccess"
    | "onCancel"
    | "onStart"
    | "onFinish"
  > &
  PropsWithChildren & {
    activeProps?: AnchorProps;

    /**
     * If true, opens in a new tab and disables referrer
     */
    external?: boolean;
    toInertia?: boolean;
    inertiaElement?: "link" | "button";
    withConfirmation?: boolean;

    /**
     * HTTP method for Inertia-powered links
     */
    method?: Method;

    /**
     * Always required
     */
    href: string;

    target?: string;
    onClick?: MouseEventHandler<Element>;
  };

export default function AppAnchor({
  external,
  toInertia,
  inertiaElement = "link",
  method = "get",
  withConfirmation,
  activeProps,
  href,
  onClick,
  ...props
}: AppAnchorProps) {
  const visit = () => {
    router.visit(href, { method });
  };

  const confirmDeleting = useConfirmation({
    callback: () => {
      visit();
    },
  });

  const handleInertiaLinkClick: MouseEventHandler<Element> = (event) => {
    onClick?.(event);

    if (!shouldIntercept(event)) {
      return;
    }

    event.preventDefault();
    if (withConfirmation) {
      confirmDeleting(event);
      return;
    }

    visit();
  };

  const handleInertiaButtonClick: MouseEventHandler<Element> = (event) => {
    onClick?.(event);

    if (!withConfirmation) {
      return;
    }

    event.preventDefault();
    confirmDeleting(event);
  };

  let additionalProps = {};
  if (activeProps) {
    const active = getCurrentUrl({ onlyPath: true }) === href;
    additionalProps = active ? activeProps : {};
  }

  if (external) {
    return (
      <Anchor
        href={href}
        target="_blank"
        rel="noopener noreferrer nofollow"
        {...props}
        onClick={onClick}
      />
    );
  }

  if (!toInertia) {
    return <Anchor href={href} {...props} {...additionalProps} />;
  }

  if (inertiaElement === "button") {
    return (
      <Button
        component={Link}
        href={href}
        method={method}
        onClick={handleInertiaButtonClick}
        {...props}
        {...additionalProps}
      />
    );
  }

  // default: Inertia-powered as regular link
  return (
    <Anchor
      href={href}
      onClick={handleInertiaLinkClick}
      {...props}
      {...additionalProps}
    />
  );
}
