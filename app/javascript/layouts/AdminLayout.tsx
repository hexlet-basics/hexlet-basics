import { Link, usePage } from "@inertiajs/react";
import {
  AppShell,
  Button,
  Divider,
  Drawer,
  Group,
  NavLink,
  Title,
} from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import {
  IconFileText,
  IconHome,
  IconLock,
  IconMessageCircle,
  IconTerminal2,
  IconUserCheck,
  IconUsers,
} from "@tabler/icons-react";
import type { PropsWithChildren, ReactNode } from "react";
import { useTranslation } from "react-i18next";
import XFlash from "@/components/XFlash.tsx";
import * as Routes from "@/routes.js";
import NavbarBlock from "./blocks/NavbarBlock.tsx";
import RootLayout from "./RootLayout.tsx";

type Props = PropsWithChildren & {
  header: string;
};

// resource — ключ из StaffMember::Role::Permission::Resource (serialize),
// должен совпадать со строками в staffPermissions.
type NavItem = {
  href: string;
  label: string;
  icon: ReactNode;
  resource: string;
};

function NavbarMenu() {
  const { t } = useTranslation();
  const { isAdmin, staffPermissions } = usePage().props;

  const canIndex = (resource: string) =>
    isAdmin || (staffPermissions?.[resource]?.can_index ?? false);

  const items: NavItem[] = [
    {
      href: Routes.admin_language_categories_path(),
      label: t(($) => $.layouts.admin.application.language_categories),
      icon: <IconFileText size={16} />,
      resource: "language_categories",
    },
    {
      href: Routes.admin_language_lessons_path(),
      label: t(($) => $.layouts.admin.application.language_lessons),
      icon: <IconFileText size={16} />,
      resource: "language_lessons",
    },
    {
      href: Routes.admin_language_lesson_reviews_path(),
      label: t(($) => $.layouts.admin.application.language_lesson_reviews),
      icon: <IconFileText size={16} />,
      resource: "language_lesson_reviews",
    },
    {
      href: Routes.admin_languages_path(),
      label: t(($) => $.layouts.admin.application.languages),
      icon: <IconTerminal2 size={16} />,
      resource: "languages",
    },
    {
      href: Routes.admin_language_landing_pages_path(),
      label: t(($) => $.layouts.admin.application.language_landing_pages),
      icon: <IconTerminal2 size={16} />,
      resource: "language_landing_pages",
    },
    {
      href: Routes.admin_reviews_path(),
      label: t(($) => $.layouts.admin.application.reviews),
      icon: <IconMessageCircle size={16} />,
      resource: "reviews",
    },
    {
      href: Routes.admin_banners_path(),
      label: t(($) => $.layouts.admin.application.banners),
      icon: <IconFileText size={16} />,
      resource: "banners",
    },
    {
      href: Routes.admin_blog_posts_path(),
      label: t(($) => $.layouts.admin.application.blog_posts),
      icon: <IconFileText size={16} />,
      resource: "blog_posts",
    },
    {
      href: Routes.admin_messages_path(),
      label: t(($) => $.layouts.admin.application.messages),
      icon: <IconFileText size={16} />,
      resource: "messages",
    },
    {
      href: Routes.admin_language_lesson_members_path(),
      label: t(($) => $.layouts.admin.application.language_lesson_members),
      icon: <IconFileText size={16} />,
      resource: "language_lesson_members",
    },
    {
      href: Routes.admin_leads_path(),
      label: t(($) => $.layouts.admin.application.leads),
      icon: <IconUserCheck size={16} />,
      resource: "leads",
    },
  ];

  const visibleItems = items.filter((item) => canIndex(item.resource));

  return (
    <>
      <NavLink
        component={Link}
        href={Routes.admin_root_path()}
        label={t(($) => $.layouts.admin.application.dashboard)}
        leftSection={<IconHome size={16} />}
      />
      {visibleItems.map((item) => (
        <NavLink
          key={item.resource}
          component={Link}
          href={item.href}
          label={item.label}
          leftSection={item.icon}
        />
      ))}
      {isAdmin && (
        <>
          <Divider />
          <NavLink
            component={Link}
            href={Routes.admin_management_users_path()}
            label={t(($) => $.layouts.admin.application.users)}
            leftSection={<IconUsers size={16} />}
          />
          <NavLink
            component={Link}
            href={Routes.admin_management_roles_path()}
            label={t(($) => $.admin.management.roles.index.header)}
            leftSection={<IconLock size={16} />}
          />
          <NavLink
            component={Link}
            href={Routes.admin_management_staff_members_path()}
            label={t(($) => $.admin.management.staff_members.index.header)}
            leftSection={<IconUsers size={16} />}
          />
        </>
      )}
    </>
  );
}

export default function AdminLayout({ children, header }: Props) {
  const { t } = useTranslation();

  const [mobileOpened] = useDisclosure();

  const [navbarOpened, { toggle: toggleNavbar }] = useDisclosure();
  const [opened, { open, close }] = useDisclosure();

  return (
    <RootLayout>
      <AppShell
        header={{ height: 60 }}
        navbar={{
          width: 250,
          breakpoint: "xs",
          collapsed: { mobile: !mobileOpened },
        }}
        padding="lg"
      >
        <AppShell.Header>
          <NavbarBlock onToggle={toggleNavbar} opened={navbarOpened} />
          <Drawer
            opened={opened}
            onClose={close}
            title={t(($) => $.layouts.admin.application.authentication)}
          >
            <NavbarMenu />
          </Drawer>
          <Group justify="end">
            <Button m="xs" variant="default" onClick={open} hiddenFrom="sm">
              {t(($) => $.layouts.admin.application.menu)}
            </Button>
          </Group>
        </AppShell.Header>

        <AppShell.Navbar>
          <NavbarMenu />
        </AppShell.Navbar>

        <AppShell.Main>
          <XFlash />
          {header && (
            <Title order={1} mb="xl">
              {header}
            </Title>
          )}
          {children}
        </AppShell.Main>
      </AppShell>
    </RootLayout>
  );
}
