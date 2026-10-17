import {
  ActionIcon,
  Avatar,
  Anchor,
  Box,
  Burger,
  Button,
  type ButtonProps,
  Divider,
  Drawer,
  Group,
  HoverCard,
  Image,
  Menu,
  NavLink,
  SimpleGrid,
  Stack,
  Text,
  ThemeIcon,
  UnstyledButton,
  useMantineColorScheme,
} from "@mantine/core";
import { ButtonLink } from "@/components/RouterLink";
import { useDisclosure } from "@mantine/hooks";
import {
  IconBlocks,
  IconChevronDown,
  IconChevronRight,
  IconGitBranch,
  IconHeartHandshake,
  IconLogout2,
  IconMoon,
  IconPlugConnected,
  IconRocket,
  IconSchool,
  IconSun,
  IconTarget,
  IconShieldCog,
  IconUserCog,
} from "@tabler/icons-react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Link, useLocation, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import logoImg from "@/assets/logo.svg";
import {
  deleteSessionMutation,
  getCurrentUserOptions,
  getCurrentUserQueryKey,
  listCoursesOptions,
} from "@/client/@tanstack/react-query.gen";
import { switchLocale } from "@/client/sdk.gen";
import type { CourseCatalogItem, Locale, User } from "@/client/types.gen";

// Header, ported from legacy NavbarBlock. Data comes from the hey-api generated
// Query hooks — the current user from `GET /me` (getCurrentUserOptions) and the
// courses menu from `listCourses` — never hand-written fetches. Navigation uses
// the typed TanStack Router `Link`; the optional `{-$locale}` prefix is
// preserved automatically, so links carry no explicit locale. Links point only
// at routes that exist today; auth-gated destinations (dashboard) are added as
// those pages are ported.
// The top-bar controls, as Mantine's header examples do them: a subtle button
// whose hover is a background tint, not an underline, and which keeps a focus
// ring. `gray` keeps the label in the body text colour in both schemes.
const navButtonProps = { variant: "subtle", color: "gray", c: "inherit", fw: 500 } as const;

function NavButton(props: ButtonProps) {
  return <Button {...navButtonProps} {...props} />;
}

const NavButtonLink: typeof ButtonLink = (props) => <ButtonLink {...navButtonProps} {...props} />;

export default function Header() {
  const { i18n } = useTranslation();
  const [opened, { toggle }] = useDisclosure();
  const { data } = useQuery(listCoursesOptions());
  const courses = data ?? [];

  return (
    <>
      <Group h="100%" px="md">
        <Anchor variant="nav" component={Link} to="/{-$locale}" params>
          <Image src={logoImg} w={30} h={30} fit="contain" alt="Logo" />
        </Anchor>

        <Group visibleFrom="sm">
          <CourseMenu courses={courses} />
          {i18n.language === "ru" && <SolutionsMenu />}
        </Group>

        <Group ms="auto" visibleFrom="sm">
          <AuthLinks />
          <LocaleSwitcher />
          <ThemeSwitcher />
        </Group>

        <Burger
          aria-label="Toggle Main Menu Visibility"
          opened={opened}
          onClick={toggle}
          ms="auto"
          hiddenFrom="sm"
          size="sm"
        />
      </Group>

      <Drawer opened={opened} onClose={toggle} hiddenFrom="sm" mt="md">
        <Stack gap="xs" px="md" align="start">
          <CourseMenu courses={courses} />
          {i18n.language === "ru" && <SolutionsMenu />}
          <AuthLinks />
          <LocaleSwitcher />
          <ThemeSwitcher />
        </Stack>
      </Drawer>
    </>
  );
}

// The signed-in user, resolved from the session cookie via `GET /me`. Returns
// null while anonymous (or until the auth backend lands, when `/me` resolves to
// no user).
function useCurrentUser(): User | null {
  const { data } = useQuery(getCurrentUserOptions());
  return data?.user ?? null;
}

function CourseMenu({ courses }: { courses: CourseCatalogItem[] }) {
  const { t } = useTranslation();

  return (
    <Menu
      width={320}
      shadow="md"
      position="bottom-start"
      trigger="hover"
      openDelay={100}
      closeDelay={150}
      withinPortal
    >
      <Menu.Target>
        <NavButton rightSection={<IconChevronDown size={16} />}>
          {t(($) => $.layouts.shared.nav.courses)}
        </NavButton>
      </Menu.Target>
      <Menu.Dropdown>
        <NavLink
          component={Link}
          to="/{-$locale}/languages"
          params
          fw="bold"
          label={t(($) => $.layouts.shared.all_courses)}
          rightSection={<IconChevronRight size={14} />}
        />
        <Divider mb="xs" />

        <SimpleGrid cols={2} spacing="sm" p="xs">
          {courses.map((item) => (
            <Menu.Item
              key={item.id}
              renderRoot={(props) => (
                <Link
                  to="/{-$locale}/languages/$slug"
                  params={(prev) => ({ ...prev, slug: item.slug })}
                  {...props}
                />
              )}
              leftSection={
                item.coverUrl ? (
                  <Image
                    w={18}
                    radius="sm"
                    fit="contain"
                    loading="lazy"
                    src={item.coverUrl}
                    alt={item.name ?? item.slug}
                  />
                ) : undefined
              }
            >
              <Text fz="sm">{item.name}</Text>
            </Menu.Item>
          ))}
        </SimpleGrid>
      </Menu.Dropdown>
    </Menu>
  );
}

function AuthLinks() {
  const user = useCurrentUser();
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const { mutate: logout } = useMutation({
    ...deleteSessionMutation(),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: getCurrentUserQueryKey() });
      void navigate({ to: "/{-$locale}" });
    },
  });

  if (!user) {
    return (
      <>
        <NavButtonLink to="/{-$locale}/session/new" params>
          {t(($) => $.layouts.shared.nav.sign_in)}
        </NavButtonLink>
        <NavButtonLink to="/{-$locale}/users/new" params>
          {t(($) => $.layouts.shared.nav.registration)}
        </NavButtonLink>
      </>
    );
  }

  return (
    <Menu shadow="md" width={260} position="bottom-end">
      <Menu.Target>
        <UnstyledButton aria-label={user.name ?? user.email ?? undefined}>
          <Group gap={6} wrap="nowrap">
            <Avatar name={user.name ?? user.email ?? ""} color="initials" size="sm" radius="xl" />
            <IconChevronDown size={14} />
          </Group>
        </UnstyledButton>
      </Menu.Target>
      <Menu.Dropdown>
        {/* The account block is the way to the learner's dashboard. */}
        <Menu.Item renderRoot={(props) => <Link to="/{-$locale}/my" params {...props} />}>
          <Group gap="sm" wrap="nowrap">
            <Avatar name={user.name ?? user.email ?? ""} color="initials" radius="xl" />
            <Box miw={0}>
              {user.name && (
                <Text fw={500} size="sm" truncate>
                  {user.name}
                </Text>
              )}
              {user.email && (
                <Text size="xs" c="dimmed" truncate>
                  {user.email}
                </Text>
              )}
            </Box>
          </Group>
        </Menu.Item>
        <Menu.Divider />
        <Menu.Item
          leftSection={<IconUserCog size={16} />}
          renderRoot={(props) => <Link to="/{-$locale}/account/profile/edit" params {...props} />}
        >
          {t(($) => $.layouts.shared.nav.profile)}
        </Menu.Item>
        <Menu.Item
          leftSection={<IconPlugConnected size={16} />}
          renderRoot={(props) => <Link to="/{-$locale}/account/connected_apps" params {...props} />}
        >
          {t(($) => $.layouts.shared.nav.connected_apps)}
        </Menu.Item>
        {user.canAccessAdmin && (
          <Menu.Item
            leftSection={<IconShieldCog size={16} />}
            renderRoot={(props) => <Link to="/{-$locale}/admin" params {...props} />}
          >
            {t(($) => $.layouts.shared.nav.admin)}
          </Menu.Item>
        )}
        <Menu.Divider />
        <Menu.Item leftSection={<IconLogout2 size={16} />} onClick={() => logout({})}>
          {t(($) => $.layouts.shared.nav.sign_out)}
        </Menu.Item>
      </Menu.Dropdown>
    </Menu>
  );
}

// Locale switcher (legacy LocalesController#switch). The choice is stored
// first — on the signed-in user, and in the cookie the site root reads — then
// the browser loads the same page, query included, under the new prefix (`/`,
// `/ru`, `/es`), as legacy's redirect did. The load happens even if storing
// failed: legacy redirected whatever happened. The href stays on the anchor so
// the menu is still a set of links to the other locales.
const LOCALES = [
  { code: "en", label: "English" },
  { code: "ru", label: "Русский" },
  { code: "es", label: "Español" },
] as const;

function LocaleSwitcher() {
  const { i18n } = useTranslation();
  const { pathname, searchStr } = useLocation();

  // Strip any existing locale prefix; `en` is served unprefixed.
  const base = pathname.replace(/^\/(ru|es)(?=\/|$)/, "") || "/";
  const hrefFor = (code: Locale) =>
    (code === "en" ? base : `/${code}${base === "/" ? "" : base}`) + searchStr;
  const switchTo = (code: Locale) => {
    void switchLocale({ query: { locale: code } }).finally(() =>
      window.location.assign(hrefFor(code)),
    );
  };

  const current = LOCALES.find((l) => l.code === i18n.language) ?? LOCALES[0];

  return (
    <Menu shadow="md" width={160}>
      <Menu.Target>
        <NavButton rightSection={<IconChevronDown size={14} />}>{current.label}</NavButton>
      </Menu.Target>
      <Menu.Dropdown>
        {LOCALES.map((locale) => (
          <Menu.Item
            key={locale.code}
            component="a"
            href={hrefFor(locale.code)}
            onClick={(event) => {
              event.preventDefault();
              switchTo(locale.code);
            }}
          >
            {locale.label}
          </Menu.Item>
        ))}
      </Menu.Dropdown>
    </Menu>
  );
}

// Both icons are rendered and the colour scheme hides one by CSS
// (Mantine's `mantine-light-hidden`/`mantine-dark-hidden`): the server cannot know the browser's scheme, so
// choosing the icon in JS would differ between the server and the first client
// render and fail hydration for every dark-theme visitor.
export function ThemeSwitcher() {
  const { t } = useTranslation();
  const { toggleColorScheme } = useMantineColorScheme();
  const label = t(($) => $.layouts.shared.nav.switch_theme);

  return (
    <ActionIcon aria-label={label} onClick={toggleColorScheme} size="sm" variant="transparent">
      <IconSun stroke={1.2} size={14} className="mantine-light-hidden" />
      <IconMoon stroke={1.2} size={14} className="mantine-dark-hidden" />
    </ActionIcon>
  );
}

// Marketing "solutions" menu, ru-only, ported from legacy. Every target but
// the teachers' case is an external Hexlet URL, kept in the locale files' links;
// that one is a local route.
// Legacy opened each in a new tab.
function SolutionsMenu() {
  const { t } = useTranslation();

  const items = [
    {
      icon: IconTarget,
      title: t(($) => $.layouts.shared.nav.courses_with_employment),
      description: t(($) => $.layouts.shared.nav.courses_with_employment_description),
      href: t(($) => $.links.hexlet_courses_for_beginners),
    },
    {
      icon: IconRocket,
      title: t(($) => $.layouts.shared.nav.career),
      description: t(($) => $.layouts.shared.nav.career_description),
      href: t(($) => $.links.hexlet_career),
    },
    {
      icon: IconGitBranch,
      title: t(($) => $.layouts.shared.nav.upskilling),
      description: t(($) => $.layouts.shared.nav.upskilling_description),
      href: t(($) => $.links.hexlet_courses_for_programmers),
    },
    {
      icon: IconHeartHandshake,
      title: t(($) => $.layouts.shared.nav.business),
      description: t(($) => $.layouts.shared.nav.business_description),
      href: t(($) => $.links.hexlet_b2b),
    },
    {
      icon: IconBlocks,
      title: t(($) => $.layouts.shared.nav.for_teachers),
      description: t(($) => $.layouts.shared.nav.for_teachers_description),
      to: "/{-$locale}/cases/for_teachers" as const,
    },
    {
      icon: IconSchool,
      title: t(($) => $.layouts.shared.nav.hexly),
      description: t(($) => $.layouts.shared.nav.hexly_description),
      href: t(($) => $.links.hexly),
    },
  ];

  const links = items.map((item) => (
    <UnstyledButton key={item.title} pos="relative">
      <Group wrap="nowrap" align="flex-start">
        <ThemeIcon size={34} variant="default" radius="md">
          <item.icon size={22} />
        </ThemeIcon>
        <Box>
          {"to" in item ? (
            <Anchor variant="nav" component={Link} to={item.to} target="_blank" rel="noreferrer">
              <Text fz="sm" fw="bold">
                {item.title}
              </Text>
            </Anchor>
          ) : (
            <Anchor variant="nav" href={item.href} target="_blank" rel="noreferrer">
              <Text fz="sm" fw="bold">
                {item.title}
              </Text>
            </Anchor>
          )}
          <Text fz="xs" c="dimmed">
            {item.description}
          </Text>
        </Box>
      </Group>
    </UnstyledButton>
  ));

  return (
    <HoverCard width={600} radius="md" shadow="md" withinPortal>
      <HoverCard.Target>
        <NavButton rightSection={<IconChevronDown size={16} />}>
          {t(($) => $.layouts.shared.nav.cases)}
        </NavButton>
      </HoverCard.Target>
      <HoverCard.Dropdown>
        <Group justify="space-between" px="md" mb="sm">
          <Text fw="bold">{t(($) => $.layouts.shared.nav.for_whom)}</Text>
        </Group>
        <Divider mb="lg" />
        <SimpleGrid cols={2} px="md" mb="md">
          {links}
        </SimpleGrid>
      </HoverCard.Dropdown>
    </HoverCard>
  );
}
