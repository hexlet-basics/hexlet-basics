import { Link, usePage } from '@inertiajs/react';
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
  SimpleGrid,
  Space,
  Stack,
  Text,
  ThemeIcon,
  UnstyledButton,
} from '@mantine/core';
import i18next from 'i18next';
import {
  Blocks,
  ChevronDown,
  GitGraph,
  GraduationCap,
  Handshake,
  LogOut,
  Rocket,
  ShieldUser,
  Target,
  User,
  UserCog,
} from 'lucide-react';
import { useTranslation } from 'react-i18next';
import AppAnchor from '@/components/AppAnchor';
import logoImg from '@/images/logo.png';
import defaultAvatarImg from '@/images/user-avatar.png';
import { localesByCode } from '@/lib/utils';
import * as Routes from '@/routes.js';
import type { SharedProps } from '@/types';

export type NavbarBlockProps = {
  opened: boolean;
  onToggle: () => void;
};

export default function NavbarBlock({ opened, onToggle }: NavbarBlockProps) {
  const { landingPagesForLists } = usePage<SharedProps>().props;

  return (
    <>
      <Group h="100%" px="md">
        <a href={Routes.root_path()} className="me-lg">
          <Image src={logoImg} w={30} h={30} fit="contain" alt="Logo" />
        </a>

        <MyLink />
        <CourseMenu landingPages={landingPagesForLists} />
        {i18next.language === 'ru' && (
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
  const { auth } = usePage<SharedProps>().props;
  const { t } = useTranslation('layouts');

  if (auth.user.guest) return null;

  return <AppAnchor href={Routes.my_path()}>{t('shared.nav.my')}</AppAnchor>;
}

function CourseMenu({
  landingPages,
}: {
  landingPages: SharedProps['landingPagesForLists'];
}) {
  const { t } = useTranslation('layouts');

  return (
    <HoverCard shadow="md" keepMounted>
      <HoverCard.Target>
        <UnstyledButton>
          <Center inline>
            <Text me={5}>{t('shared.nav.courses')}</Text>
            <ChevronDown size={16} />
          </Center>
        </UnstyledButton>
      </HoverCard.Target>
      <HoverCard.Dropdown>
        <SimpleGrid cols={2} spacing="sm" p="xs">
          {landingPages.map((lp) => (
            <Group key={lp.id} pos="relative">
              <Image
                w="auto"
                radius="sm"
                fit="contain"
                loading="lazy"
                src={lp.language.cover_thumb_variant}
                alt={lp.header}
              />
              <Text fz="sm">{lp.header}</Text>
              <AppAnchor
                href={Routes.language_path(lp.slug)}
                inset={0}
                pos="absolute"
              />
            </Group>
          ))}
        </SimpleGrid>
      </HoverCard.Dropdown>
    </HoverCard>
  );
}

function BookLink() {
  const { t } = useTranslation('layouts');
  return i18next.language === 'ru' ? (
    <AppAnchor pseudo href={Routes.book_path()}>
      {t('shared.nav.book')}
    </AppAnchor>
  ) : null;
}

function AuthLinks({ avatar }: { avatar: string }) {
  const { auth } = usePage<SharedProps>().props;
  const { t } = useTranslation('layouts');

  if (auth.user.guest) {
    return (
      <>
        <AppAnchor pseudo href={Routes.new_session_path()}>
          {t('shared.nav.sign_in')}
        </AppAnchor>
        <AppAnchor pseudo href={Routes.new_user_path()}>
          {t('shared.nav.registration')}
        </AppAnchor>
      </>
    );
  }

  return (
    <Menu shadow="md" width={250}>
      <Menu.Target>
        <UnstyledButton>
          <Center>
            <User size={18} />
            {auth.user.name && (
              <>
                <Space me="xs" />
                <Text>{auth.user.name}</Text>
              </>
            )}
            <ChevronDown size={14} />
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
          leftSection={<UserCog size={14} />}
          component={Link}
          href={Routes.edit_account_profile_path()}
        >
          {t('shared.nav.profile')}
        </Menu.Item>
        {auth.user.admin && (
          <Menu.Item
            leftSection={<ShieldUser size={14} />}
            component={Link}
            href={Routes.admin_root_path()}
          >
            {t('shared.nav.admin')}
          </Menu.Item>
        )}
        <Menu.Item
          leftSection={<LogOut size={14} />}
          component={Link}
          method="delete"
          href={Routes.session_path()}
        >
          {t('shared.nav.sign_out')}
        </Menu.Item>
      </Menu.Dropdown>
    </Menu>
  );
}

function LocaleSwitcher() {
  return (
    <Menu shadow="md" width={200}>
      <Menu.Target>
        <UnstyledButton>
          <Center>
            {localesByCode[i18next.language || 'ru'].icon}{' '}
            <ChevronDown size={14} />
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
  landingPages: SharedProps['landingPagesForLists'];
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
  const { t: tLayouts } = useTranslation('layouts');
  const { t: tCommon } = useTranslation('common');

  const solutionMenuData = [
    {
      icon: Target,
      title: tLayouts('shared.nav.courses_with_employement'),
      description: tLayouts('shared.nav.courses_with_employement_description'),
      href: 'https://ru.hexlet.io/courses_for_beginners?utm_source=code-basics&utm_medium=referral',
    },
    {
      icon: Rocket,
      title: tLayouts('shared.nav.career'),
      description: tLayouts('shared.nav.career_description'),
      href: 'https://career.hexlet.io?utm_source=code-basics&utm_medium=referral',
    },
    {
      icon: GitGraph,
      title: tLayouts('shared.nav.upskilling'),
      description: tLayouts('shared.nav.upskilling_description'),
      href: 'https://ru.hexlet.io/courses_for_programmers?utm_source=code-basics&utm_medium=referral',
    },
    {
      icon: Handshake,
      title: tLayouts('shared.nav.business'),
      description: tLayouts('shared.nav.business_description'),
      href: 'https://b2b.hexlet.io?utm_source=code-basics&utm_medium=referral',
    },
    {
      icon: Blocks,
      title: tLayouts('shared.nav.for_teachers'),
      description: tLayouts('shared.nav.for_teachers_description'),
      href: Routes.for_teachers_cases_path(),
    },
    {
      icon: GraduationCap,
      title: tLayouts('shared.nav.hexly'),
      description: tLayouts('shared.nav.hexly_description'),
      href: 'https://hexly.ru?utm_source=code-basics&utm_medium=referral',
    },
  ];

  const links = solutionMenuData.map((item) => (
    <UnstyledButton key={item.title} pos="relative">
      <Group wrap="nowrap" align="flex-start">
        <ThemeIcon size={34} variant="default" radius="md">
          <item.icon size={22} />
        </ThemeIcon>
        <Box>
          <Text fz="sm" fw={500}>
            {item.title}
          </Text>
          <Text fz="xs" c="dimmed">
            {item.description}
          </Text>
        </Box>
      </Group>
      <AppAnchor external href={item.href} inset={0} pos="absolute" />
    </UnstyledButton>
  ));

  return (
    <HoverCard width={600} radius="md" shadow="md" withinPortal>
      <HoverCard.Target>
        <UnstyledButton>
          <Center inline>
            <Text mr={5}>{tLayouts('shared.nav.cases')}</Text>
            <ChevronDown size={16} />
          </Center>
        </UnstyledButton>
      </HoverCard.Target>

      <HoverCard.Dropdown>
        <Group justify="space-between" px="md" mb="sm">
          <Text fw={500}>{tLayouts('shared.nav.for_whom')}</Text>
          <Anchor
            target="_blank"
            href={`${tCommon('organization.site')}?utm_source=code-basics&utm_medium=referral`}
            fz="xs"
          >
            {tCommon('organization.site')}
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
