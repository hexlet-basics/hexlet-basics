import { Link } from "@inertiajs/react";
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
} from "@mantine/core";
import { IconMessage2 } from "@tabler/icons-react";
import i18next from "i18next";
import { Trans, useTranslation } from "react-i18next";
import CourseBlock from "@/components/CourseBlock";
import AppAnchor from "@/components/Elements/AppAnchor";
import ApplicationLayout from "@/layouts/ApplicationLayout";
import * as Routes from "@/routes.js";
import type {
  LanguageLandingPageForLists,
  LanguageMember,
} from "@/types/serializers";

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
  const { t } = useTranslation();

  return (
    <ApplicationLayout>
      <Container>
        <Title order={2} my="xl">
          {t(($) => $.my.show.started)}
        </Title>
        {startedCourseMembers.length === 0 && <p>{t(($) => $.common.empty)}</p>}

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
          {t(($) => $.my.show.finished)}
        </Title>

        {finishedCourseMembers.length > 0 && i18next.language === "ru" && (
          <Alert icon={<IconMessage2 />} mb="xl" px="xl">
            <Trans
              t={t}
              i18nKey={($) => $.my.show.add_review}
              components={{
                a: (
                  <AppAnchor
                    external
                    href="https://taplink.cc/codebasics_reviews"
                  />
                ),
              }}
            />
          </Alert>
        )}

        {finishedCourseMembers.length === 0 && (
          <p>{t(($) => $.common.empty)}</p>
        )}

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
