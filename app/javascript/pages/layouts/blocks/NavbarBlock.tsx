import logoImg from "@/images/logo.png";
import defaultAvatarImg from "@/images/user-avatar.png";
import { localesByCode } from "@/lib/utils";
import i18next from "i18next";
import {
  Burger,
  Group,
  Image,
  Anchor, Menu,
  UnstyledButton,
  Text,
  Stack,
  Center,
  Divider,
  Avatar,
  Drawer
} from "@mantine/core";

import * as Routes from "@/routes.js";
import type { SharedProps } from "@/types";
import { Link, usePage } from "@inertiajs/react";
import { useTranslation } from "react-i18next";
import { ChevronDown, User } from "lucide-react";
import AppAnchor from "@/components/AppAnchor";

export type NavbarBlockProps = {
  opened: boolean;
  onToggle: () => void;
};

export default function NavbarBlock({ opened, onToggle }: NavbarBlockProps) {
  const { landingPagesForLists } = usePage<SharedProps>().props;

  return (
    <>
      <Group h="100%" px="md">
        <Group gap={15}>
          <Anchor component={Link} href={Routes.root_path()} me="lg">
            <Image src={logoImg} width={30} height={30} alt="Logo" />
          </Anchor>

          <MyLink />
          <CourseMenu landingPages={landingPagesForLists} />
          <Group visibleFrom="sm">
            <CaseLink />
          </Group>
          <Group visibleFrom="sm">
          <BookLink />
          </Group>
        </Group>

        <Group ms="auto" visibleFrom="sm">
          <AuthLinks avatar={defaultAvatarImg} />
          <LocaleSwitcher />
        </Group>

        <Burger opened={opened} onClick={onToggle} ms="auto" hiddenFrom="sm" size="sm" />
      </Group>

      <Drawer opened={opened} onClose={onToggle} hiddenFrom="sm" mt="md">
        <MobileMenu landingPages={landingPagesForLists} avatar={defaultAvatarImg} />
      </Drawer>
    </>
  );
}

function MyLink() {
  const { auth } = usePage<SharedProps>().props;
  const { t } = useTranslation("layouts");

  if (auth.user.guest) return null;

  return (
    <AppAnchor href={Routes.my_path()}>
      {t("shared.nav.my")}
    </AppAnchor>
  );
}

function CourseMenu({ landingPages }: { landingPages: SharedProps["landingPagesForLists"] }) {
  const { t } = useTranslation("layouts");

  return (
    <Menu shadow="md">
      <Menu.Target>
        <UnstyledButton>
          <Center>
            {t("shared.nav.courses")} <ChevronDown size={14} />
          </Center>
        </UnstyledButton>
      </Menu.Target>
      <Menu.Dropdown>
        {landingPages.map((lp) => (
          <Menu.Item key={lp.id} component={Link} href={Routes.language_path(lp.slug)}>
            <Group wrap="nowrap">
              <Image w="auto" fit="contain" src={lp.language.cover_thumb_variant} alt={lp.header} />
              <Text>{lp.header}</Text>
            </Group>
          </Menu.Item>
        ))}
      </Menu.Dropdown>
    </Menu>
  );
}

function CaseLink() {
  const { t } = useTranslation("layouts");
  return (
    <Menu shadow="md" width={200}>
      <Menu.Target>
        <UnstyledButton>
          <Center>
            {t("shared.nav.cases")} <ChevronDown size={14} />
          </Center>
        </UnstyledButton>
      </Menu.Target>
      <Menu.Dropdown>
        <Menu.Item component={Link} href={Routes.for_teachers_cases_path()}>
          {t("shared.nav.for_teachers")}
        </Menu.Item>
      </Menu.Dropdown>
    </Menu>
  );
}

function BookLink() {
  const { t } = useTranslation("layouts");
  return i18next.language === "ru" ? (
    <Anchor component={Link} href={Routes.book_path()}>
      {t("shared.nav.book")}
    </Anchor>
  ) : null;
}

function AuthLinks({ avatar }: { avatar: string }) {
  const { auth } = usePage<SharedProps>().props;
  const { t } = useTranslation("layouts");

  if (auth.user.guest) {
    return (
      <>
        <AppAnchor href={Routes.new_session_path()}>
          {t("shared.nav.sign_in")}
        </AppAnchor>
        <AppAnchor href={Routes.new_user_path()}>
          {t("shared.nav.registration")}
        </AppAnchor>
      </>
    );
  }

  return (
    <Menu shadow="md" width={200}>
      <Menu.Target>
        <UnstyledButton>
          <Center>
            <User size={18} /> <ChevronDown size={14} />
          </Center>
        </UnstyledButton>
      </Menu.Target>
      <Menu.Dropdown>
        <Menu.Item component={Link} href={Routes.edit_account_profile_path()}>
          <Group>
            <Avatar src={avatar} radius="xl" size="lg" />
            <Stack gap={0}>
              <Text fw={500}>{auth.user.name}</Text>
              <Text size="xs" c="dimmed">
                {auth.user.email}
              </Text>
            </Stack>
          </Group>
        </Menu.Item>
        <Divider />
        <Menu.Item component={Link} href={Routes.edit_account_profile_path()}>
          {t("shared.nav.profile")}
        </Menu.Item>
        {auth.user.admin && (
          <Menu.Item component={Link} href={Routes.admin_root_path()}>
            {t("shared.nav.admin")}
          </Menu.Item>
        )}
        <Menu.Item component={Link} method="delete" href={Routes.session_path()}>
          {t("shared.nav.sign_out")}
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
            {localesByCode[i18next.language || "ru"].icon} <ChevronDown size={14} />
          </Center>
        </UnstyledButton>
      </Menu.Target>
      <Menu.Dropdown>
        {Object.entries(localesByCode).map(([k, v]) => (
          <Menu.Item key={k} component="a" href={Routes.switch_locale_path({ new_locale: k })}>
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

function MobileMenu({ landingPages, avatar }: { landingPages: SharedProps["landingPagesForLists"]; avatar: string }) {
  return (
    <Stack gap="xs" px="md" align="start">
      <MyLink />
      <CourseMenu landingPages={landingPages} />
      <CaseLink />
      <BookLink />
      <AuthLinks avatar={avatar} />
      <LocaleSwitcher />
    </Stack>
  );
}
