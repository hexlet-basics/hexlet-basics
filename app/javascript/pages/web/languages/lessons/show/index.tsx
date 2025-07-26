import { AppShell, useMantineTheme } from '@mantine/core';
import { useMediaQuery } from '@mantine/hooks';
import { GripVertical } from 'lucide-react';
import { Panel, PanelGroup, PanelResizeHandle } from 'react-resizable-panels';
import LessonLayout from '@/pages/layouts/LessonLayout.tsx';
import LessonMainContent from './components/LessonMainContent.tsx';
import LessonNavbar from './components/LessonNavbar.tsx';

export default function Index() {
  const theme = useMantineTheme();
  const isTabletUp = useMediaQuery(`(min-width: ${theme.breakpoints.sm})`);

  if (isTabletUp) {
    return (
      <LessonLayout>
        <PanelGroup direction="horizontal">
          <Panel minSize={25} defaultSize={30}>
            <LessonNavbar />
          </Panel>
          <PanelResizeHandle className="w-3 bg-gray-200 group flex items-center justify-center relative hover:bg-gray-300 transition-colors">
            <GripVertical className="w-4 h-4 text-gray-500 group-hover:text-gray-700 bg-gray-200" />
          </PanelResizeHandle>
          <Panel minSize={25} defaultSize={70}>
            <LessonMainContent />
          </Panel>
        </PanelGroup>
      </LessonLayout>
    );
  }

  return (
    <LessonLayout>
      <AppShell.Navbar>
        <LessonNavbar />
      </AppShell.Navbar>
      <AppShell.Main h="100%">
        <LessonMainContent />
      </AppShell.Main>
    </LessonLayout>
  );
}
