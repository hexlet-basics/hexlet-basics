import { Link, usePage } from '@inertiajs/react';
import { Group, NavLink } from '@mantine/core';

export type CrudHorizontalMenuItem = {
  label: string | React.ReactNode;
  href: string;
  external?: boolean;
};

type Props = {
  items: CrudHorizontalMenuItem[];
};

export function CrudHorizontalMenu({ items }: Props) {
  const { url } = usePage();

  const getCurrentArea = (href: string) => (url === href ? 'page' : undefined);

  return (
    <Group mb="md" gap={0}>
      {items.map(({ label, href, external }) =>
        external ? (
          <NavLink
            key={href}
            w="auto"
            href={href}
            target="_blank"
            component="a"
            aria-current={getCurrentArea(href)}
            label={label}
          />
        ) : (
          <NavLink
            key={href}
            w="auto"
            href={href}
            component={Link}
            aria-current={getCurrentArea(href)}
            label={label}
          />
        ),
      )}
    </Group>
  );
}
