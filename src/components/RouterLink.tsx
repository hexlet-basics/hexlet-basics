import {
  ActionIcon,
  type ActionIconProps,
  Button,
  type ButtonProps,
  NavLink as MantineNavLink,
  type NavLinkProps,
  Text,
  type TextProps,
} from "@mantine/core";
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

type TextLinkProps = Omit<TextProps, "href">;

const TextLinkBase = forwardRef<HTMLAnchorElement, TextLinkProps>((props, ref) => (
  <Text ref={ref} component="a" {...props} />
));

const TextLinkCreated = createLink(TextLinkBase);

export const TextLink: LinkComponent<typeof TextLinkBase> = (props) => (
  <TextLinkCreated preload="intent" {...props} />
);

type NavLinkLinkProps = Omit<NavLinkProps, "href">;

const NavLinkBase = forwardRef<HTMLAnchorElement, NavLinkLinkProps>((props, ref) => (
  <MantineNavLink ref={ref} component="a" {...props} />
));

const NavLinkCreated = createLink(NavLinkBase);

// Mantine NavLink as a typed router link — the lesson navigation list.
export const NavLink: LinkComponent<typeof NavLinkBase> = (props) => (
  <NavLinkCreated preload="intent" {...props} />
);
