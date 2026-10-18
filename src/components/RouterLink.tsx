import { ActionIcon, Anchor, Button, Card, NavLink as MantineNavLink, Text } from "@mantine/core";
import type {
  AnchorProps,
  ActionIconProps,
  ButtonProps,
  CardProps,
  NavLinkProps,
  TextProps,
} from "@mantine/core";
import { createLink } from "@tanstack/react-router";
import type { LinkComponent } from "@tanstack/react-router";
import { forwardRef } from "react";
import type { PropsWithChildren } from "react";

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

// An Anchor rather than TextLink's bare Text: Mantine styles Anchor with the
// theme's link colour and underline rules, which `Text component="a"` skips (it
// shows the browser's default blue). Used where a link should look like one.
type AnchorLinkProps = Omit<AnchorProps, "href">;

const AnchorLinkBase = forwardRef<HTMLAnchorElement, AnchorLinkProps>((props, ref) => (
  <Anchor ref={ref} {...props} />
));

const AnchorLinkCreated = createLink(AnchorLinkBase);

export const AnchorLink: LinkComponent<typeof AnchorLinkBase> = (props) => (
  <AnchorLinkCreated preload="intent" {...props} />
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

type CardLinkProps = PropsWithChildren<Omit<CardProps, "href">>;

const CardLinkBase = forwardRef<HTMLAnchorElement, CardLinkProps>((props, ref) => (
  <Card ref={ref} component="a" {...props} />
));

const CardLinkCreated = createLink(CardLinkBase);

// Mantine Card as a typed router link — the whole card is the click target.
export const CardLink: LinkComponent<typeof CardLinkBase> = (props) => (
  <CardLinkCreated preload="intent" {...props} />
);
