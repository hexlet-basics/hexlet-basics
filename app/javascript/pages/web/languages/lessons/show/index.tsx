import { Link, usePage } from '@inertiajs/react';
import {
  Accordion,
  Alert,
  Anchor,
  AppShell,
  Box,
  Center,
  Divider,
  Group,
  List,
  NavLink,
  Paper,
  ScrollArea,
  Stack,
  Tabs,
  Text,
  Title,
  useMantineTheme,
} from '@mantine/core';
import { useMediaQuery } from '@mantine/hooks';
import i18next from 'i18next';
import { BookOpenText, Github, GripVertical, Info, Rocket } from 'lucide-react'; // TODO: the current github icon is deprecated, need to update from a new lib
import { Suspense, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Panel, PanelGroup, PanelResizeHandle } from 'react-resizable-panels';
import { XBreadcrumb } from '@/components/breadcrumbs.tsx';
import Chat from '@/components/Chat.tsx';
import MarkdownViewer from '@/components/MarkdownViewer.tsx';
import XssContent from '@/components/XssContent.tsx';
import { useIsMobile } from '@/hooks/useIsMobile.ts';
import { neededPreview } from '@/lib/utils.ts';
import ContactMethodRequestingBlock from '@/pages/layouts/blocks/ContactMethodRequestingBlock.tsx';
import LessonLayout from '@/pages/layouts/LessonLayout.tsx';
import * as Routes from '@/routes.js';
import type { BreadcrumbItem } from '@/types/index.ts';
import ControlBox from './components/ControlBox.tsx';
import EditorTab from './components/EditorTab.tsx';
import HTMLPreview from './components/HTMLPreview.tsx';
import OutputTab from './components/OutputTab.tsx';
import SolutionTab from './components/SolutionTab.tsx';
import TestsTab from './components/TestsTab.tsx';
import { useLessonStore } from './store.tsx';
import type { LessonSharedProps } from './types.ts';

function HtmlPreviewBlock() {
  const { course } = usePage<LessonSharedProps>().props;
  const currentTab = useLessonStore((state) => state.currentTab);
  const content = useLessonStore((state) => state.content);

  if (currentTab !== 'editor') {
    return null;
  }
  if (!neededPreview(course.slug!)) {
    return null;
  }

  return (
    <Suspense fallback={<div>Loading...</div>}>
      <HTMLPreview html={content} />
    </Suspense>
  );
}

export default function Index() {
  const theme = useMantineTheme();
  const isTabletUp = !useIsMobile();
  const headerHeight =
    theme.components?.AppShell?.defaultProps?.header?.height ?? 60;

  if (isTabletUp) {
    return (
      <LessonLayout>
        <Box
          component={PanelGroup}
          direction="horizontal"
          pt={headerHeight}
          h="100%"
        >
          <Panel minSize={20} defaultSize={45}>
            <Box h="100%">
              <LessonLeftBlock />
            </Box>
          </Panel>
          <PanelResizeHandle>
            <Center h="100%" w={10} bg={theme.colors.gray[2]}>
              <GripVertical size={10} color={theme.colors.gray[6]} />
            </Center>
          </PanelResizeHandle>
          <Panel minSize={30} defaultSize={55}>
            <Box h="100%">
              <LessonRightBlock />
            </Box>
          </Panel>
        </Box>
      </LessonLayout>
    );
  }

  return (
    <LessonLayout>
      <AppShell.Navbar>
        <LessonLeftBlock />
      </AppShell.Navbar>
      <AppShell.Main h="100%">
        <LessonRightBlock />
      </AppShell.Main>
    </LessonLayout>
  );
}

function LessonLeftBlock() {
  const { t } = useTranslation();
  const [focusesCount, setFocusCount] = useState(0);
  const handleSelect = (selectedKey: string | null) => {
    if (selectedKey === 'assistant') {
      setFocusCount((count) => count + 1);
    }
  };

  return (
    <Tabs
      defaultValue="lesson"
      onChange={handleSelect}
      h="100%"
      display="flex"
      style={{ flexDirection: 'column' }}
    >
      <Tabs.List grow>
        <Tabs.Tab value="lesson">{t('languages.lessons.show.lesson')}</Tabs.Tab>
        <Tabs.Tab value="assistant">
          {t('languages.lessons.show.discuss')}
        </Tabs.Tab>
        <Tabs.Tab value="navigation">
          {t('languages.lessons.show.navigation')}
        </Tabs.Tab>
      </Tabs.List>

      <AppShell.Section grow mih={0}>
        <Tabs.Panel value="lesson" h="100%">
          <ScrollArea h="100%">
            <LessonTabContent />
          </ScrollArea>
        </Tabs.Panel>

        <Tabs.Panel value="assistant" h="100%">
          <ScrollArea h="100%">
            <AssistantTabContent focusesCount={focusesCount} />
          </ScrollArea>
        </Tabs.Panel>

        <Tabs.Panel value="navigation" h="100%">
          <ScrollArea h="100%">
            <NavigationTabContent />
          </ScrollArea>
        </Tabs.Panel>
      </AppShell.Section>
    </Tabs>
  );
}

function LessonRightBlock() {
  const { t } = useTranslation();
  const changeTab = useLessonStore((state) => state.changeTab);
  const currentTab = useLessonStore((state) => state.currentTab);

  return (
    <Tabs
      h="100%"
      display="flex"
      style={{ flexDirection: 'column' }}
      value={currentTab}
      onChange={(key) => changeTab(key as typeof currentTab)}
      // keepMounted={false}
    >
      <Tabs.List grow>
        <Tabs.Tab value="lesson" hiddenFrom="sm">
          <Center>
            <BookOpenText size={14} />
          </Center>
        </Tabs.Tab>
        <Tabs.Tab value="editor">{t('languages.lessons.show.editor')}</Tabs.Tab>
        <Tabs.Tab value="output">{t('languages.lessons.show.output')}</Tabs.Tab>
        <Tabs.Tab value="tests">{t('languages.lessons.show.tests')}</Tabs.Tab>
        <Tabs.Tab value="solution">
          {t('languages.lessons.show.solution')}
        </Tabs.Tab>
      </Tabs.List>

      <AppShell.Section grow mih={0}>
        <Tabs.Panel value="lesson" h="100%" hiddenFrom="sm">
          <Stack h="100%" gap={0}>
            <ScrollArea h="100%">
              <LessonTabContent />
            </ScrollArea>
            <Box style={{ flexShrink: 0 }}>
              <ControlBox />
            </Box>
          </Stack>
        </Tabs.Panel>

        <Tabs.Panel value="editor" h="100%">
          <Stack h="100%" gap={0} mih={0} style={{ flexGrow: 1 }}>
            <Stack mih={0} style={{ flexGrow: 1 }}>
              <EditorTab />
            </Stack>

            <Box style={{ flexShrink: 0 }}>
              <HtmlPreviewBlock />
              <ControlBox />
            </Box>
          </Stack>
        </Tabs.Panel>

        <Tabs.Panel value="output" h="100%">
          <Stack h="100%" gap={0}>
            <ScrollArea h="100%">
              <OutputTab />
            </ScrollArea>
            <Box style={{ flexShrink: 0 }}>
              <ControlBox />
            </Box>
          </Stack>
        </Tabs.Panel>

        <Tabs.Panel value="tests" h="100%">
          <Stack h="100%" gap={0}>
            <ScrollArea h="100%">
              <Box p="md">
                <TestsTab />
              </Box>
            </ScrollArea>
            <ControlBox />
          </Stack>
        </Tabs.Panel>

        <Tabs.Panel value="solution" h="100%">
          <Stack h="100%" gap={0}>
            <ScrollArea h="100%">
              <Box p="md">
                <SolutionTab />
              </Box>
            </ScrollArea>
            <ControlBox />
          </Stack>
        </Tabs.Panel>
      </AppShell.Section>
    </Tabs>
  );
}

function LessonTabContent() {
  const { t } = useTranslation();
  const {
    courseCategory,
    landingPage,
    course,
    lesson,
    shouldAddContactMethod,
  } = usePage<LessonSharedProps>().props;

  const commonQuestions = t('languages.lessons.show.common_questions', {
    returnObjects: true,
  }) as Array<{ question: string; answer: string }>;

  const items: BreadcrumbItem[] = [
    {
      name: courseCategory?.name ?? '-',
      url: courseCategory
        ? Routes.language_category_url(courseCategory.slug!)
        : '#',
    },
    {
      name: landingPage.header!,
      url: Routes.language_url(landingPage.slug!),
    },
    {
      name: lesson.name!,
      url: Routes.language_lesson_url(landingPage.slug!, lesson.slug!),
    },
  ];

  return (
    <Box p="lg">
      <XBreadcrumb items={items} />
      <Title my="sm">{`${landingPage.name}: ${lesson.name}`}</Title>

      {shouldAddContactMethod && (
        <Paper withBorder shadow="sm" p="md" mt="md">
          <ContactMethodRequestingBlock />
        </Paper>
      )}

      <MarkdownViewer allowHtml>{lesson.theory || ''}</MarkdownViewer>

      <Title order={2} my="md">
        {t('languages.lessons.show.instructions')}
      </Title>
      <MarkdownViewer allowHtml>{lesson.instructions || ''}</MarkdownViewer>

      {course.hexlet_program_landing_page && (
        <Alert variant="primary" my="xl" radius="lg">
          <Group justify="center" gap={6}>
            <Rocket size={15} />
            <Anchor
              target="_blank"
              href={`${course.hexlet_program_landing_page}?utm_source=code-basics&utm_medium=referral`}
            >
              {t('languages.lessons.show.profession_description')}
            </Anchor>
          </Group>
        </Alert>
      )}

      {lesson.tips.length > 0 && (
        <>
          <Title order={2} my="md">
            {t('languages.lessons.show.tips')}
          </Title>
          <ul>
            {lesson.tips.map((tip) => (
              <li key={tip}>
                <MarkdownViewer allowHtml>{tip}</MarkdownViewer>
              </li>
            ))}
          </ul>
        </>
      )}

      <Divider my="xl" c="gray" />

      <Accordion mb="xs" defaultValue={commonQuestions[0]?.question}>
        {commonQuestions.map((v) => (
          <Accordion.Item key={v.question} value={v.question}>
            <Accordion.Control>
              <Text>{v.question}</Text>
            </Accordion.Control>
            <Accordion.Panel>
              <MarkdownViewer allowHtml>{v.answer}</MarkdownViewer>
            </Accordion.Panel>
          </Accordion.Item>
        ))}
      </Accordion>

      <Center>
        <Text fz="sm" me="sm" component="span">
          {t('languages.lessons.show.issues')}
        </Text>
        <a
          href={lesson.source_code_url!}
          target="_blank"
          rel="noreferrer noopener"
        >
          <Github size={12} />
        </a>
      </Center>
    </Box>
  );
}

function AssistantTabContent({ focusesCount }: { focusesCount: number }) {
  const { t } = useTranslation();
  const { t: tCommon } = useTranslation('common');

  const {
    previousMessages,
    canCreateAssistantMessage,
    course,
    lesson,
    lessonMember,
  } = usePage<LessonSharedProps>().props;

  const userCode = useLessonStore((state) => state.content);
  const output = useLessonStore((state) => state.output);

  return (
    <Box p="lg">
      {i18next.language === 'ru' && (
        <Alert icon={<Info />} mb="lg">
          <XssContent>
            {t('languages.lessons.show.if_stuck_html', {
              url: tCommon('community_url'),
            })}
          </XssContent>
        </Alert>
      )}
      <Chat
        focusesCount={focusesCount}
        previousMessages={previousMessages}
        enabled={canCreateAssistantMessage}
        userCode={userCode}
        output={output}
        course={course}
        lesson={lesson}
        lessonMember={lessonMember}
      />
    </Box>
  );
}

function NavigationTabContent() {
  const { lesson, lessons, landingPage } = usePage<LessonSharedProps>().props;
  return (
    <Box p="lg">
      {lessons.map((l) => (
        <NavLink
          component={Link}
          key={l.id}
          href={Routes.language_lesson_path(
            landingPage.language.slug!,
            l.slug!,
          )}
          label={l.name}
          active={l.id === lesson.id}
        />
      ))}
    </Box>
  );
}
