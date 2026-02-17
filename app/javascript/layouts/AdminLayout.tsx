import { Link } from "@inertiajs/react";
import {
  AppShell,
  Box,
  Button,
  Container,
  Divider,
  Drawer,
  Grid,
  Group,
  NavLink,
  Title,
} from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import {
  IconFileText,
  IconHelpHexagon,
  IconHome,
  IconLayoutGrid,
  IconMessageCircle,
  IconRoadSign,
  IconTerminal2,
  IconUserCheck,
  IconUsers,
} from "@tabler/icons-react";
import type { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import XFlash from "@/components/XFlash.tsx";
import * as Routes from "@/routes.js";
import FooterBlock from "./blocks/FooterBlock.tsx";
import NavbarBlock from "./blocks/NavbarBlock.tsx";
import RootLayout from "./RootLayout.tsx";

type Props = PropsWithChildren & {
  header: string;
};

function NavbarMenu() {
  const { t } = useTranslation();

  return (
    <>
      <NavLink
        component={Link}
        href={Routes.admin_root_path()}
        label={t(($) => $.layouts.admin.application.dashboard)}
        leftSection={<IconHome size={16} />}
      />
      <NavLink
        component={Link}
        href={Routes.admin_language_categories_path()}
        label={t(($) => $.layouts.admin.application.language_categories)}
        leftSection={<IconFileText size={16} />}
      />
      <NavLink
        component={Link}
        href={Routes.admin_language_lessons_path()}
        label={t(($) => $.layouts.admin.application.language_lessons)}
        leftSection={<IconFileText size={16} />}
      />
      <NavLink
        component={Link}
        href={Routes.admin_language_lesson_reviews_path()}
        label={t(($) => $.layouts.admin.application.language_lesson_reviews)}
        leftSection={<IconFileText size={16} />}
      />
      <NavLink
        component={Link}
        href={Routes.admin_languages_path()}
        label={t(($) => $.layouts.admin.application.languages)}
        leftSection={<IconTerminal2 size={16} />}
      />
      <NavLink
        component={Link}
        href={Routes.admin_language_landing_pages_path()}
        label={t(($) => $.layouts.admin.application.language_landing_pages)}
        leftSection={<IconTerminal2 size={16} />}
      />
      <NavLink
        component={Link}
        href={Routes.admin_reviews_path()}
        label={t(($) => $.layouts.admin.application.reviews)}
        leftSection={<IconMessageCircle size={16} />}
      />
      <NavLink
        component={Link}
        href={Routes.admin_blog_posts_path()}
        label={t(($) => $.layouts.admin.application.blog_posts)}
        leftSection={<IconFileText size={16} />}
      />
      <NavLink
        component={Link}
        href={Routes.admin_messages_path()}
        label={t(($) => $.layouts.admin.application.messages)}
        leftSection={<IconFileText size={16} />}
      />
      <NavLink
        component={Link}
        href={Routes.admin_language_lesson_members_path()}
        label={t(($) => $.layouts.admin.application.language_lesson_members)}
        leftSection={<IconFileText size={16} />}
      />
      <NavLink
        component={Link}
        href={Routes.admin_surveys_path()}
        label={t(($) => $.layouts.admin.application.surveys)}
        leftSection={<IconHelpHexagon size={16} />}
      />
      <NavLink
        component={Link}
        href={Routes.admin_survey_scenarios_path()}
        label={t(($) => $.layouts.admin.application.survey_scenarios)}
        leftSection={<IconLayoutGrid size={16} />}
      />
      <NavLink
        component={Link}
        href={Routes.admin_survey_answers_path()}
        label={t(($) => $.layouts.admin.application.survey_answers)}
        leftSection={<IconRoadSign size={16} />}
      />
      <NavLink
        component={Link}
        href={Routes.admin_leads_url()}
        label={t(($) => $.layouts.admin.application.leads)}
        leftSection={<IconUserCheck size={16} />}
      />
      <Divider />
      <NavLink
        component={Link}
        href={Routes.admin_management_users_path()}
        label={t(($) => $.layouts.admin.application.users)}
        leftSection={<IconUsers size={16} />}
      />
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
