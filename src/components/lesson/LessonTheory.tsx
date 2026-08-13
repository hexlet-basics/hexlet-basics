import { Accordion, Anchor, Box, Center, Divider, Text, Title } from "@mantine/core";
import { IconBrandGithub } from "@tabler/icons-react";
import { useTranslation } from "react-i18next";
import type { CourseLessonView } from "@/client/types.gen";
import Breadcrumbs, { CurrentCrumb } from "@/components/Breadcrumbs";
import MarkdownViewer from "@/components/MarkdownViewer";
import { TextLink } from "@/components/RouterLink";

// Everything a learner reads: where they are, the theory, what they are asked to
// do, the tips, the questions everyone asks, and the way to fix what they just
// read. Ported from legacy's LessonTabContent.
//
// Breadcrumbs run home → course → lesson. Legacy has a category above the
// course, but the public category read is unimplemented and blocked on a missing
// schema, so that level waits for the category work rather than being faked.
export default function LessonTheory({ view }: { view: CourseLessonView }) {
  const { t } = useTranslation();
  const { lesson, landingPage } = view;
  const courseSlug = lesson.course.slug;

  // The marketing name of the course, not the course's own `name` — the string a
  // learner recognises, and what legacy titles this page with.
  const courseName = landingPage?.name ?? lesson.course.name ?? courseSlug;

  const questions = t(($) => $.courses.lessons.show.common_questions, {
    returnObjects: true,
  });

  return (
    <Box p="lg">
      <Breadcrumbs homeLabel={t(($) => $.courses.lessons.show.to_home_title)}>
        <TextLink to="/{-$locale}/languages/$slug" params={{ slug: courseSlug }} size="sm">
          {landingPage?.header ?? courseName}
        </TextLink>
        <CurrentCrumb>{lesson.name}</CurrentCrumb>
      </Breadcrumbs>

      <Title my="sm">{`${courseName}: ${lesson.name}`}</Title>

      <MarkdownViewer allowHtml>{lesson.theory ?? ""}</MarkdownViewer>

      <Title order={2} my="md">
        {t(($) => $.courses.lessons.show.instructions)}
      </Title>
      <MarkdownViewer allowHtml>{lesson.instructions ?? ""}</MarkdownViewer>

      {lesson.tips.length > 0 && (
        <>
          <Title order={2} my="md">
            {t(($) => $.courses.lessons.show.tips)}
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

      <Accordion mb="xs" defaultValue={questions[0]?.question}>
        {questions.map((entry) => (
          <Accordion.Item key={entry.question} value={entry.question}>
            <Accordion.Control>
              <Text>{entry.question}</Text>
            </Accordion.Control>
            <Accordion.Panel>
              <MarkdownViewer allowHtml>{entry.answer}</MarkdownViewer>
            </Accordion.Panel>
          </Accordion.Item>
        ))}
      </Accordion>

      {lesson.sourceCodeUrl && (
        <Center>
          <Text fz="sm" me="sm" component="span">
            {t(($) => $.courses.lessons.show.issues)}
          </Text>
          <Anchor
            href={lesson.sourceCodeUrl}
            target="_blank"
            rel="noopener noreferrer"
            aria-label={t(($) => $.courses.lessons.show.source_code)}
          >
            <IconBrandGithub size={12} aria-hidden="true" />
          </Anchor>
        </Center>
      )}
    </Box>
  );
}
