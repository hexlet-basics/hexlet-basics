import type { PageProps } from "@inertiajs/core";
import { Link, usePage } from "@inertiajs/react";
import {
  Anchor,
  Box,
  Burger,
  Center,
  Divider,
  Drawer,
  Group,
  HoverCard,
  Image,
  Menu,
  Popover,
  SimpleGrid,
  Space,
  Stack,
  Text,
  ThemeIcon,
  UnstyledButton,
} from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import {
  IconBlocks,
  IconChevronDown,
  IconGitBranch,
  IconHeartHandshake,
  IconLogout2,
  IconRocket,
  IconSchool,
  IconTarget,
  IconUser,
  IconUserCog,
  IconUserShield,
} from "@tabler/icons-react";
import i18next from "i18next";
import { useTranslation } from "react-i18next";
import AppAnchor from "@/components/Elements/AppAnchor";
import { useIsMobile } from "@/hooks/useIsMobile";
import logoImg from "@/images/logo.svg";
import defaultAvatarImg from "@/images/user-avatar.webp";
import { hasObjectKey, localesByCode } from "@/lib/utils";
import * as Routes from "@/routes.js";

export type NavbarBlockProps = {
  opened: boolean;
  onToggle: () => void;
};

export default function NavbarBlock({ opened, onToggle }: NavbarBlockProps) {
  const { landingPagesForLists } = usePage().props;

  return (
    <>
      <Group h="100%" px="md">
        <AppAnchor href={Routes.root_path()} className="me-lg">
          <Image src={logoImg} w={30} h={30} fit="contain" alt="Logo" />
        </AppAnchor>

        <MyLink />
        <CourseMenu landingPages={landingPagesForLists} />
        {i18next.language === "ru" && (
          <Group visibleFrom="sm">
            <SolutionsMenu />
          </Group>
        )}
        <Group visibleFrom="sm">
          <BookLink />
        </Group>

        <Group ms="auto" visibleFrom="sm">
          <AuthLinks avatar={defaultAvatarImg} />
          <LocaleSwitcher />
        </Group>

        <Burger
          aria-label="Toggle Main Menu Visibility"
          opened={opened}
          onClick={onToggle}
          ms="auto"
          hiddenFrom="sm"
          size="sm"
        />
      </Group>

      <Drawer opened={opened} onClose={onToggle} hiddenFrom="sm" mt="md">
        <MobileMenu
          landingPages={landingPagesForLists}
          avatar={defaultAvatarImg}
        />
      </Drawer>
    </>
  );
}

function MyLink() {
  const { auth } = usePage().props;
  const { t } = useTranslation();

  if (auth.user.guest) return null;

  return (
    <AppAnchor href={Routes.my_path()}>
      {t(($) => $.layouts.shared.nav.my)}
    </AppAnchor>
  );
}

function CourseMenu({
  landingPages,
}: {
  landingPages: PageProps["landingPagesForLists"];
}) {
  const { t } = useTranslation();
  const isMobile = useIsMobile();
  const [opened, { open, close, toggle }] = useDisclosure(false);

  return (
    <Popover
      width={320}
      position="bottom-start"
      withArrow
      shadow="md"
      opened={opened}
      onDismiss={close}
      withinPortal
    >
      <Popover.Target>
        <UnstyledButton
          onClick={isMobile ? toggle : undefined}
          onMouseEnter={!isMobile ? open : undefined}
          onMouseLeave={!isMobile ? close : undefined}
        >
          <Center inline>
            <Text me={5}>{t(($) => $.layouts.shared.nav.courses)}</Text>
            <IconChevronDown size={16} />
          </Center>
        </UnstyledButton>
      </Popover.Target>
      <Popover.Dropdown
        onMouseEnter={!isMobile ? open : undefined}
        onMouseLeave={!isMobile ? close : undefined}
      >
        <SimpleGrid cols={2} spacing="sm" p="xs">
          {landingPages.map((lp) => (
            <Group key={lp.id} wrap="nowrap" pos="relative">
              <Image
                w="auto"
                radius="sm"
                fit="contain"
                loading="lazy"
                src={lp.language.cover_thumb_variant}
                alt={lp.header}
              />
              <Text fz="sm">{lp.name}</Text>
              <AppAnchor
                href={Routes.language_path(lp.slug)}
                inset={0}
                pos="absolute"
              />
            </Group>
          ))}
        </SimpleGrid>
      </Popover.Dropdown>
    </Popover>
  );
}

function BookLink() {
  const { t } = useTranslation();
  return i18next.language === "ru" ? (
    <AppAnchor href={Routes.book_path()}>
      {t(($) => $.layouts.shared.nav.book)}
    </AppAnchor>
  ) : null;
}

function AuthLinks({ avatar }: { avatar: string }) {
  const { auth } = usePage().props;
  const { t } = useTranslation();

  if (auth.user.guest) {
    return (
      <>
        <AppAnchor href={Routes.new_session_path()}>
          {t(($) => $.layouts.shared.nav.sign_in)}
        </AppAnchor>
        <AppAnchor href={Routes.new_user_path()}>
          {t(($) => $.layouts.shared.nav.registration)}
        </AppAnchor>
      </>
    );
  }

  return (
    <Menu shadow="md" width={250}>
      <Menu.Target>
        <UnstyledButton>
          <Center>
            <IconUser size={18} />
            {auth.user.name && (
              <>
                <Space me="xs" />
                <Text>{auth.user.name}</Text>
              </>
            )}
            <IconChevronDown size={14} />
          </Center>
        </UnstyledButton>
      </Menu.Target>
      <Menu.Dropdown>
        <Menu.Label>
          <Stack gap={0}>
            <Text size="xs" c="dimmed">
              {auth.user.email}
            </Text>
          </Stack>
        </Menu.Label>
        <Menu.Item
          leftSection={<IconUserCog size={14} />}
          component={Link}
          href={Routes.edit_account_profile_path()}
        >
          {t(($) => $.layouts.shared.nav.profile)}
        </Menu.Item>
        {auth.user.admin && (
          <Menu.Item
            leftSection={<IconUserShield size={14} />}
            component={Link}
            href={Routes.admin_root_path()}
          >
            {t(($) => $.layouts.shared.nav.admin)}
          </Menu.Item>
        )}
        <Menu.Item
          leftSection={<IconLogout2 size={14} />}
          component={Link}
          method="delete"
          href={Routes.session_path()}
        >
          {t(($) => $.layouts.shared.nav.sign_out)}
        </Menu.Item>
      </Menu.Dropdown>
    </Menu>
  );
}

function LocaleSwitcher() {
  const localeKey = hasObjectKey(localesByCode, i18next.language)
    ? i18next.language
    : "ru";

  return (
    <Menu shadow="md" width={200}>
      <Menu.Target>
        <UnstyledButton>
          <Center>
            {localesByCode[localeKey].icon} <IconChevronDown size={14} />
          </Center>
        </UnstyledButton>
      </Menu.Target>
      <Menu.Dropdown>
        {Object.entries(localesByCode).map(([k, v]) => (
          <Menu.Item
            key={k}
            component="a"
            href={Routes.switch_locale_path({ new_locale: k })}
          >
            <Group gap={5}>
              {v.icon}
              <Text>{v.name}</Text>
            </Group>
          </Menu.Item>
        ))}
      </Menu.Dropdown>
    </Menu>
  );
}

function MobileMenu({
  landingPages,
  avatar,
}: {
  landingPages: PageProps["landingPagesForLists"];
  avatar: string;
}) {
  return (
    <Stack gap="xs" px="md" align="start">
      <MyLink />
      <CourseMenu landingPages={landingPages} />
      <SolutionsMenu />
      <BookLink />
      <AuthLinks avatar={avatar} />
      <LocaleSwitcher />
    </Stack>
  );
}

function SolutionsMenu() {
  const { t } = useTranslation();

  const solutionMenuData = [
    {
      icon: IconTarget,
      title: t(($) => $.layouts.shared.nav.courses_with_employement),
      description: t(
        ($) => $.layouts.shared.nav.courses_with_employement_description,
      ),
      href: "https://ru.hexlet.io/courses_for_beginners?utm_source=code-basics&utm_medium=referral",
    },
    {
      icon: IconRocket,
      title: t(($) => $.layouts.shared.nav.career),
      description: t(($) => $.layouts.shared.nav.career_description),
      href: "https://career.hexlet.io?utm_source=code-basics&utm_medium=referral",
    },
    {
      icon: IconGitBranch,
      title: t(($) => $.layouts.shared.nav.upskilling),
      description: t(($) => $.layouts.shared.nav.upskilling_description),
      href: "https://ru.hexlet.io/courses_for_programmers?utm_source=code-basics&utm_medium=referral",
    },
    {
      icon: IconHeartHandshake,
      title: t(($) => $.layouts.shared.nav.business),
      description: t(($) => $.layouts.shared.nav.business_description),
      href: "https://b2b.hexlet.io?utm_source=code-basics&utm_medium=referral",
    },
    {
      icon: IconBlocks,
      title: t(($) => $.layouts.shared.nav.for_teachers),
      description: t(($) => $.layouts.shared.nav.for_teachers_description),
      href: Routes.for_teachers_cases_path(),
    },
    {
      icon: IconSchool,
      title: t(($) => $.layouts.shared.nav.hexly),
      description: t(($) => $.layouts.shared.nav.hexly_description),
      href: "https://hexly.ru?utm_source=code-basics&utm_medium=referral",
    },
  ];

  const links = solutionMenuData.map((item) => (
    <UnstyledButton key={item.title} pos="relative">
      <Group wrap="nowrap" align="flex-start">
        <ThemeIcon size={34} variant="default" radius="md">
          <item.icon size={22} />
        </ThemeIcon>
        <Box>
          <AppAnchor
            external
            href={item.href}
            className="after:absolute after:inset-0"
          >
            <Text fz="sm" fw={500}>
              {item.title}
            </Text>
          </AppAnchor>
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
        <UnstyledButton>
          <Center inline>
            <Text mr={5}>{t(($) => $.layouts.shared.nav.cases)}</Text>
            <IconChevronDown size={16} />
          </Center>
        </UnstyledButton>
      </HoverCard.Target>
      <HoverCard.Dropdown>
        <Group justify="space-between" px="md" mb="sm">
          <Text fw={500}>{t(($) => $.layouts.shared.nav.for_whom)}</Text>
          <Anchor
            target="_blank"
            href={`${t(($) => $.common.organization.site)}?utm_source=code-basics&utm_medium=referral`}
            fz="xs"
          >
            {t(($) => $.common.organization.site)}
          </Anchor>
        </Group>

        <Divider mb="lg" />

        <SimpleGrid cols={2} spacing="md" px="md" mb="md">
          {links}
        </SimpleGrid>
      </HoverCard.Dropdown>
    </HoverCard>
  );
}
