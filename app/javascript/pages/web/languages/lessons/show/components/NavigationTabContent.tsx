import { usePage } from '@inertiajs/react';
import { Box, List } from '@mantine/core';
import AppAnchor from '@/components/AppAnchor.tsx';
import type { LessonSharedProps } from '@/pages/web/languages/lessons/show/types.ts';
import * as Routes from '@/routes';

export default function NavigationTabContent() {
  const { lessons, landingPage } = usePage<LessonSharedProps>().props;

  return (
    <Box p="lg">
      <List type="ordered">
        {lessons.map((l) => (
          <List.Item key={l.id}>
            <AppAnchor
              href={Routes.language_lesson_path(
                landingPage.language.slug!,
                l.slug!,
              )}
            >
              {l.name}
            </AppAnchor>
          </List.Item>
        ))}
      </List>
    </Box>
  );
}
