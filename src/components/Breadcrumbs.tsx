import { Breadcrumbs as MantineBreadcrumbs, Center, Text } from "@mantine/core";
import { IconHome } from "@tabler/icons-react";
import type { ReactNode } from "react";
import { TextLink } from "@/components/RouterLink";

// Breadcrumb trail, ported from legacy's XBreadcrumb. The home crumb is always
// first and is an icon; the last crumb is the current page and is not a link.
//
// Callers pass their crumbs as rendered TextLinks, so destinations stay typed
// against the route tree (the createLink convention) instead of being smuggled
// through as strings. The home crumb's label comes from the caller too — the
// pages that have breadcrumbs each keep their own copy in their own namespace.
export default function Breadcrumbs({
  homeLabel,
  children,
}: {
  homeLabel: string;
  children: ReactNode;
}) {
  return (
    <MantineBreadcrumbs styles={{ breadcrumb: { whiteSpace: "normal" } }}>
      <TextLink to="/{-$locale}" size="sm" aria-label={homeLabel}>
        <Center c="dimmed">
          <IconHome size={15} aria-hidden="true" />
        </Center>
      </TextLink>
      {children}
    </MantineBreadcrumbs>
  );
}

// The trailing crumb: the page you are on, so not a link.
export function CurrentCrumb({ children }: { children: ReactNode }) {
  return (
    <Text c="dimmed" size="sm">
      {children}
    </Text>
  );
}
