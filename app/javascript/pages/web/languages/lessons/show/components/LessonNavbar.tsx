import {
  AppShell,
  Box,
  ScrollArea,
  Tabs,
  useMantineTheme,
} from '@mantine/core';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import AssistantTabContent from '@/pages/web/languages/lessons/show/components/AssistantTabContent.tsx';
import LessonTabContent from '@/pages/web/languages/lessons/show/components/LessonTabContent.tsx';
import NavigationTabContent from '@/pages/web/languages/lessons/show/components/NavigationTabContent.tsx';

export default function LessonNavbar() {
  const { t } = useTranslation();
  const [focusesCount, setFocusCount] = useState(0);
  const theme = useMantineTheme();
  const headerHeight =
    theme.components?.AppShell?.defaultProps?.header?.height ?? 60;

  const handleSelect = (selectedKey: string | null) => {
    if (selectedKey === 'assistant') {
      setFocusCount((count) => count + 1);
    }
  };

  return (
    <Box h="100%" pt={headerHeight}>
      <Tabs
        defaultValue="lesson"
        onChange={handleSelect}
        style={{ height: '100%', display: 'flex', flexDirection: 'column' }}
      >
        <Tabs.List grow>
          <Tabs.Tab value="lesson">
            {t('languages.lessons.show.lesson')}
          </Tabs.Tab>
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
    </Box>
  );
}
