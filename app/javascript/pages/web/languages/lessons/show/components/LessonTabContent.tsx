import { usePage } from '@inertiajs/react';
import {
  Accordion,
  Alert,
  Anchor,
  Box,
  Center,
  Divider,
  Group,
  Paper,
  Text,
  Title,
} from '@mantine/core';
import { Github, Rocket } from 'lucide-react'; // TODO: the current github icon is deprecated, need to update from a new lib
import { useTranslation } from 'react-i18next';
import { XBreadcrumb } from '@/components/breadcrumbs.tsx';
import MarkdownViewer from '@/components/MarkdownViewer.tsx';
import ContactMethodRequestingBlock from '@/pages/layouts/blocks/ContactMethodRequestingBlock.tsx';
import type { LessonSharedProps } from '@/pages/web/languages/lessons/show/types.ts';
import * as Routes from '@/routes';
import type { BreadcrumbItem } from '@/types';

export default function LessonTabContent() {
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
