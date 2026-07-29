import { ActionIcon, type ActionIconProps, Button, type ButtonProps } from "@mantine/core";
import { createLink, type LinkComponent } from "@tanstack/react-router";
import { forwardRef } from "react";

// Typed TanStack Router links rendered as Mantine controls. Mantine's
// polymorphic `component={Link}` loses the router's `to`/`params` inference (so a
// route param has to be smuggled through `renderRoot`), whereas `createLink` —
// the documented UI-library integration — keeps navigation fully type-checked
// against the generated route tree and renders a real anchor (preload,
// middle-click, right-click open). Use these instead of `component={Link}`
// whenever a Mantine control navigates.

type ButtonLinkProps = Omit<ButtonProps, "href">;

const ButtonLinkBase = forwardRef<HTMLAnchorElement, ButtonLinkProps>((props, ref) => (
  <Button ref={ref} component="a" {...props} />
));

const ButtonLinkCreated = createLink(ButtonLinkBase);

export const ButtonLink: LinkComponent<typeof ButtonLinkBase> = (props) => (
  <ButtonLinkCreated preload="intent" {...props} />
);

type ActionIconLinkProps = Omit<ActionIconProps, "href">;

const ActionIconLinkBase = forwardRef<HTMLAnchorElement, ActionIconLinkProps>((props, ref) => (
  <ActionIcon ref={ref} component="a" {...props} />
));

const ActionIconLinkCreated = createLink(ActionIconLinkBase);

export const ActionIconLink: LinkComponent<typeof ActionIconLinkBase> = (props) => (
  <ActionIconLinkCreated preload="intent" {...props} />
);
