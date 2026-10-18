import { AppShell, Burger, Group, NavLink, Title } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { useQuery } from "@tanstack/react-query";
import { Link, Outlet, createFileRoute, useLocation } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { adminGetMyCapabilitiesOptions } from "@/client/@tanstack/react-query.gen";
import { requireAdmin } from "@/lib/auth";

// Admin back-office layout. The guard runs in beforeLoad so a non-admin never
// even loads a child route's data (SSR-safe: `user` is resolved in the root).
export const Route = createFileRoute("/{-$locale}/admin")({
  beforeLoad: ({ context, location }) => {
    requireAdmin(context.user, location.href);
  },
  component: AdminLayout,
});

// Nav entries. An entry with a `capability` is shown only to a staff member whose
// role holds that capability Permission (ADR-0020): the administrator flag does
// not imply it, so the menu asks the API what the caller holds. Each new CRUD
// resource adds one line here; the screens themselves
// are driven by the shared engine (CrudList/CrudForm). Ordered content-first,
// then community/marketing, then management — mirrors how the legacy back
// office grouped its menu.
const NAV = [
  {
    key: "courses",
    // Path fragment used to mark the entry active regardless of locale prefix.
    match: "/admin/courses",
    to: "/{-$locale}/admin/courses",
  },
  {
    key: "courseCategories",
    match: "/admin/course_categories",
    to: "/{-$locale}/admin/course_categories",
  },
  {
    key: "landingPages",
    match: "/admin/course_landing_pages",
    to: "/{-$locale}/admin/course_landing_pages",
  },
  {
    key: "lessonReviews",
    match: "/admin/course_lesson_reviews",
    to: "/{-$locale}/admin/course_lesson_reviews",
  },
  {
    key: "lessonProgress",
    match: "/admin/lesson_progress",
    to: "/{-$locale}/admin/lesson_progress",
  },
  {
    key: "blogPosts",
    match: "/admin/blog_posts",
    to: "/{-$locale}/admin/blog_posts",
  },
  {
    key: "banners",
    match: "/admin/banners",
    to: "/{-$locale}/admin/banners",
  },
  {
    key: "reviews",
    match: "/admin/reviews",
    to: "/{-$locale}/admin/reviews",
  },
  {
    key: "messages",
    match: "/admin/messages",
    to: "/{-$locale}/admin/messages",
  },
  {
    key: "leads",
    match: "/admin/leads",
    to: "/{-$locale}/admin/leads",
  },
  {
    key: "users",
    match: "/admin/users",
    to: "/{-$locale}/admin/users",
  },
  {
    key: "roles",
    match: "/admin/management/roles",
    to: "/{-$locale}/admin/management/roles",
  },
  {
    key: "staffMembers",
    match: "/admin/management/staff_members",
    to: "/{-$locale}/admin/management/staff_members",
  },
  {
    key: "managementUsers",
    match: "/admin/management/users",
    to: "/{-$locale}/admin/management/users",
  },
  {
    capability: "audit.read",
    key: "audit",
    match: "/admin/audit",
    to: "/{-$locale}/admin/audit",
  },
] as const;

// The capability an entry needs, when it needs one.
const capabilityOf = (entry: (typeof NAV)[number]): string | undefined =>
  "capability" in entry ? entry.capability : undefined;

function AdminLayout() {
  const { t } = useTranslation();
  const [opened, { toggle }] = useDisclosure();
  const { pathname } = useLocation();
  // Until the answer arrives nothing behind a capability is shown, so the menu
  // never offers an entry it may have to take back.
  const { data: capabilities } = useQuery(adminGetMyCapabilitiesOptions());
  const held = capabilities?.held ?? [];
  const entries = NAV.filter((entry) => {
    const needed = capabilityOf(entry);
    return needed === undefined || held.includes(needed);
  });

  return (
    <AppShell
      header={{ height: 56 }}
      navbar={{
        breakpoint: "sm",
        collapsed: { mobile: !opened },
        width: 260,
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
        {entries.map((entry) => (
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
