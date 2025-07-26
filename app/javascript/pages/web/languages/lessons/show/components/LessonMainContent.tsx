import {
  AppShell,
  Box,
  Center,
  ScrollArea,
  Stack,
  Tabs,
  useMantineTheme,
} from '@mantine/core';
import { BookOpenText } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { useLessonStore } from '../store.tsx';
import ControlBox from './ControlBox.tsx';
import EditorTab from './EditorTab.tsx';
import HTMLPreviewBlock from './HTMLPreviewBlock.tsx';
import LessonTabContent from './LessonTabContent.tsx';
import OutputTab from './OutputTab.tsx';
import SolutionTab from './SolutionTab.tsx';
import TestsTab from './TestsTab.tsx';

export default function LessonMainContent() {
  const { t } = useTranslation();
  const changeTab = useLessonStore((state) => state.changeTab);
  const currentTab = useLessonStore((state) => state.currentTab);
  const theme = useMantineTheme();
  const headerHeight =
    theme.components?.AppShell?.defaultProps?.header?.height ?? 60;

  return (
    <Box h="100%" style={{ paddingTop: headerHeight }}>
      <Tabs
        style={{ height: '100%', display: 'flex', flexDirection: 'column' }}
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
          <Tabs.Tab value="editor">
            {t('languages.lessons.show.editor')}
          </Tabs.Tab>
          <Tabs.Tab value="output">
            {t('languages.lessons.show.output')}
          </Tabs.Tab>
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
                <HTMLPreviewBlock />
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
    </Box>
  );
}
