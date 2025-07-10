import { Link } from '@inertiajs/react';
import { Anchor, Card, Group, Image, Stack, Text, Title } from '@mantine/core';
import { Clock, Users } from 'lucide-react';
import type { PropsWithChildren } from 'react';
import * as Routes from '@/routes.js';
import type {
  LanguageLandingPageForLists,
  LanguageMember,
} from '@/types/serializers';
import AppAnchor from './AppAnchor';

type Props = PropsWithChildren & {
  // course: Language;
  landingPage: LanguageLandingPageForLists;
  courseMember?: LanguageMember;
  continueButton?: boolean;
  lazy?: boolean;
};

export default function CourseBlock({
  landingPage,
  courseMember,
  lazy,
}: Props) {
  return (
    <Card pos="relative" shadow="sm">
      <Card.Section>
        <Image
          loading={lazy ? 'lazy' : 'eager'}
          src={landingPage.language.cover_list_variant}
          alt={landingPage.header}
        />
      </Card.Section>

      <Stack pt="md" h="100%">
        <Title fw="bold" order={2} mb="md">
          {courseMember && <i className="me-3 bi bi-trophy" />}
          {landingPage.name}
        </Title>
        <Group c="dimmed" mt="auto">
          <Group gap="xs">
            <Clock size="15" />
            <Text size="sm">{landingPage.duration}</Text>
          </Group>
          <Group gap="xs">
            <Users size="15" />
            <Text size="sm">{landingPage.members_count}</Text>
          </Group>
        </Group>
      </Stack>

      <AppAnchor
        pos="absolute"
        inset={0}
        href={Routes.language_url(landingPage.slug)}
      />
    </Card>
    //     {
    //   continueButton && (
    //     <Button variant="light" fullWidth>
    //       {tHelpers("continue")}
    //     </Button>
    //   )
    // }
  );
}
