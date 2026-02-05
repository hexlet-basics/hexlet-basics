import { Link } from '@inertiajs/react';
import {
  Alert,
  Anchor,
  Card,
  Container,
  Grid,
  Group,
  Image,
  Paper,
  Progress,
  SimpleGrid,
  Stack,
  Text,
  Title,
} from '@mantine/core';
import i18next from 'i18next';
import { MessageSquareDiff } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import CourseBlock from '@/components/CourseBlock';
import AppAnchor from '@/components/Elements/AppAnchor';
import XssContent from '@/components/XssContent';
import ApplicationLayout from '@/pages/layouts/ApplicationLayout';
import * as Routes from '@/routes.js';
import type {
  LanguageLandingPageForLists,
  LanguageMember,
} from '@/types/serializers';

type Props = {
  startedCourseMembers: LanguageMember[];
  finishedCourseMembers: LanguageMember[];
  landingPageResourcesByCourseId: Record<number, LanguageLandingPageForLists>;
};

function StartedCourse({
  lp,
  cm,
}: {
  lp: LanguageLandingPageForLists;
  cm: LanguageMember;
}) {
  return (
    <Card p={0} radius="md" withBorder pos="relative">
      <Group wrap="nowrap" align="stretch">
        <Image
          visibleFrom="sm"
          alt={lp.header}
          // w="100%"
          maw={150}
          w="auto"
          fit="contain"
          src={lp.language.cover_list_variant}
        />

        <Stack w="100%" p="sm" gap="xs">
          <AppAnchor
            className="after:absolute after:inset-0"
            href={Routes.language_url(lp.slug!)}
          >
            <Title order={3}>{lp.header}</Title>
          </AppAnchor>

          {cm.next_lesson_name && <Text>{cm.next_lesson_name} →</Text>}

          <Progress.Root mt="auto">
            <Progress.Section
              aria-label={`${cm.progress}%`}
              value={cm.progress}
            />
          </Progress.Root>
        </Stack>
      </Group>
    </Card>
  );
}

export default function My(props: Props) {
  const {
    startedCourseMembers,
    finishedCourseMembers,
    landingPageResourcesByCourseId,
  } = props;
  const { t: tViews } = useTranslation('web');
  const { t: tCommon } = useTranslation('common');

  return (
    <ApplicationLayout>
      <Container>
        <Title order={2} my="xl">
          {tViews(($) => $.my.show.started)}
        </Title>
        {startedCourseMembers.length === 0 && <p>{tCommon(($) => $.empty)}</p>}

        <SimpleGrid cols={{ base: 1, xs: 2 }}>
          {startedCourseMembers.map((cm) => (
            <StartedCourse
              key={cm.id}
              cm={cm}
              lp={landingPageResourcesByCourseId[cm.language_id]}
            />
          ))}
        </SimpleGrid>

        <Title my="xl" order={2}>
          {tViews(($) => $.my.show.finished)}
        </Title>

        {finishedCourseMembers.length > 0 && i18next.language === 'ru' && (
          <Alert icon={<MessageSquareDiff />} mb="xl" px="xl">
            <XssContent>{tViews(($) => $.my.show.add_review)}</XssContent>
          </Alert>
        )}

        {finishedCourseMembers.length === 0 && <p>{tCommon(($) => $.empty)}</p>}

        <SimpleGrid cols={{ base: 2, xs: 3, sm: 4 }}>
          {finishedCourseMembers.map((cm) => (
            <CourseBlock
              key={cm.id}
              courseMember={cm}
              landingPage={landingPageResourcesByCourseId[cm.language_id]}
            />
          ))}
        </SimpleGrid>
      </Container>
    </ApplicationLayout>
  );
}
