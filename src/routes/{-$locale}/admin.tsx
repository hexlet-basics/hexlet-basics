import { AppShell, Burger, Group, NavLink, Title } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import {
  createFileRoute,
  Link,
  Outlet,
  useLocation,
} from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { requireAdmin } from "@/lib/auth";

// Admin back-office layout. The guard runs in beforeLoad so a non-admin never
// even loads a child route's data (SSR-safe: `user` is resolved in the root).
export const Route = createFileRoute("/{-$locale}/admin")({
  beforeLoad: ({ context, location }) => {
    requireAdmin(context.user, location.href);
  },
  component: AdminLayout,
});

// Nav entries. Each new CRUD resource adds one line here; the screens themselves
// are driven by the shared engine (CrudList/CrudForm).
const NAV = [
  {
    to: "/{-$locale}/admin/language_categories",
    // Path fragment used to mark the entry active regardless of locale prefix.
    match: "/admin/language_categories",
    key: "courseCategories",
  },
  {
    to: "/{-$locale}/admin/banners",
    match: "/admin/banners",
    key: "banners",
  },
] as const;

function AdminLayout() {
  const { t } = useTranslation();
  const [opened, { toggle }] = useDisclosure();
  const { pathname } = useLocation();

  return (
    <AppShell
      header={{ height: 56 }}
      navbar={{
        width: 260,
        breakpoint: "sm",
        collapsed: { mobile: !opened },
      }}
      padding="md"
    >
      <AppShell.Header>
        <Group h="100%" px="md">
          <Burger opened={opened} onClick={toggle} hiddenFrom="sm" size="sm" />
          <Title order={4}>{t(($) => $.admin.title)}</Title>
        </Group>
      </AppShell.Header>

      <AppShell.Navbar p="md">
        {NAV.map((entry) => (
          <NavLink
            key={entry.to}
            component={Link}
            to={entry.to}
            label={t(($) => $.admin.resources[entry.key])}
            active={pathname.includes(entry.match)}
          />
        ))}
      </AppShell.Navbar>

      <AppShell.Main>
        <Outlet />
      </AppShell.Main>
    </AppShell>
  );
}
