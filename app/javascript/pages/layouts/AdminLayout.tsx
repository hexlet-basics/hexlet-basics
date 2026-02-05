import { Link } from '@inertiajs/react';
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
} from '@mantine/core';
import { useDisclosure } from '@mantine/hooks';
import type { PropsWithChildren } from 'react';
import { useTranslation } from 'react-i18next';
import XFlash from '@/components/XFlash.tsx';
import * as Routes from '@/routes.js';
import FooterBlock from './blocks/FooterBlock.tsx';
import NavbarBlock from './blocks/NavbarBlock.tsx';
import RootLayout from './RootLayout.tsx';

type Props = PropsWithChildren & {
  header: string;
};

function NavbarMenu() {
  const { t: tLayouts } = useTranslation('layouts');

  return (
    <>
      <NavLink
        component={Link}
        href={Routes.admin_root_path()}
        label={tLayouts(($) => $.web.admin.application.dashboard)}
        leftSection={<i className="bi bi-house" />}
      />
      <NavLink
        component={Link}
        href={Routes.admin_language_categories_path()}
        label={tLayouts(($) => $.web.admin.application.language_categories)}
        leftSection={<i className="bi bi-file-text" />}
      />
      <NavLink
        component={Link}
        href={Routes.admin_language_lessons_path()}
        label={tLayouts(($) => $.web.admin.application.language_lessons)}
        leftSection={<i className="bi bi-file-text" />}
      />
      <NavLink
        component={Link}
        href={Routes.admin_language_lesson_reviews_path()}
        label={tLayouts(($) => $.web.admin.application.language_lesson_reviews)}
        leftSection={<i className="bi bi-file-text" />}
      />
      <NavLink
        component={Link}
        href={Routes.admin_languages_path()}
        label={tLayouts(($) => $.web.admin.application.languages)}
        leftSection={<i className="bi bi-terminal" />}
      />
      <NavLink
        component={Link}
        href={Routes.admin_language_landing_pages_path()}
        label={tLayouts(($) => $.web.admin.application.language_landing_pages)}
        leftSection={<i className="bi bi-terminal" />}
      />
      <NavLink
        component={Link}
        href={Routes.admin_reviews_path()}
        label={tLayouts(($) => $.web.admin.application.reviews)}
        leftSection={<i className="bi bi-chat-left-quote" />}
      />
      <NavLink
        component={Link}
        href={Routes.admin_blog_posts_path()}
        label={tLayouts(($) => $.web.admin.application.blog_posts)}
        leftSection={<i className="bi bi-file-text" />}
      />
      <NavLink
        component={Link}
        href={Routes.admin_messages_path()}
        label={tLayouts(($) => $.web.admin.application.messages)}
        leftSection={<i className="bi bi-file-text" />}
      />
      <NavLink
        component={Link}
        href={Routes.admin_language_lesson_members_path()}
        label={tLayouts(($) => $.web.admin.application.language_lesson_members)}
        leftSection={<i className="bi bi-file-text" />}
      />
      <NavLink
        component={Link}
        href={Routes.admin_surveys_path()}
        label={tLayouts(($) => $.web.admin.application.surveys)}
        leftSection={<i className="bi bi-patch-question" />}
      />
      <NavLink
        component={Link}
        href={Routes.admin_survey_scenarios_path()}
        label={tLayouts(($) => $.web.admin.application.survey_scenarios)}
        leftSection={<i className="bi bi-collection" />}
      />
      <NavLink
        component={Link}
        href={Routes.admin_survey_answers_path()}
        label={tLayouts(($) => $.web.admin.application.survey_answers)}
        leftSection={<i className="bi bi-signpost-2" />}
      />
      <NavLink
        component={Link}
        href={Routes.admin_leads_url()}
        label={tLayouts(($) => $.web.admin.application.leads)}
        leftSection={<i className="bi bi-person-check" />}
      />
      <Divider />
      <NavLink
        component={Link}
        href={Routes.admin_management_users_path()}
        label={tLayouts(($) => $.web.admin.application.users)}
        leftSection={<i className="bi bi-people" />}
      />
    </>
  );
}

export default function AdminLayout({ children, header }: Props) {
  const { t: tLayouts } = useTranslation('layouts');
  const [mobileOpened] = useDisclosure();

  const [navbarOpened, { toggle: toggleNavbar }] = useDisclosure();
  const [opened, { open, close }] = useDisclosure();

  return (
    <RootLayout>
      <AppShell
        header={{ height: 60 }}
        navbar={{
          width: 250,
          breakpoint: 'xs',
          collapsed: { mobile: !mobileOpened },
        }}
        padding="lg"
      >
        <AppShell.Header>
          <NavbarBlock onToggle={toggleNavbar} opened={navbarOpened} />
          <Drawer opened={opened} onClose={close} title="Authentication">
            <NavbarMenu />
          </Drawer>
          <Group justify="end">
            <Button m="xs" variant="default" onClick={open} hiddenFrom="sm">
              Menu
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
