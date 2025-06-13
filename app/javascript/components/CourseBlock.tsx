import * as Routes from "@/routes.js";
import type { SharedProps } from "@/types";
import type {
  LanguageLandingPageForLists,
  LanguageMember,
} from "@/types/serializers";
import { Link, usePage } from "@inertiajs/react";
import type { PropsWithChildren } from "react";
import { Card, Button, Text, Group, Stack, Image, Title, Anchor } from '@mantine/core';
import { useTranslation } from "react-i18next";
import { Clock, Users } from "lucide-react";

type Props = PropsWithChildren & {
  // course: Language;
  landingPage: LanguageLandingPageForLists;
  courseMember?: LanguageMember;
  continueButton?: boolean;
};

export default function CourseBlock({
  landingPage,
  courseMember,
  continueButton,
}: Props) {
  const { suffix } = usePage<SharedProps>().props;
  const { t: tHelpers } = useTranslation("helpers");

  return (

    <Card pos="relative" shadow="sm">
      <Card.Section>
        <Image
          loading="lazy"
          src={landingPage.language.cover_list_variant}
          alt={landingPage.header}
        />
      </Card.Section>

      <Card.Section p="md">
        <Title fw="bold" order={2}>
          {courseMember && <i className="me-3 bi bi-trophy" />}
          {landingPage.header}
        </Title>
        <Group c="dimmed" mt="xs">
          <Group gap="xs">
            <Clock size="15" />
            <Text size="sm">{landingPage.duration}</Text>
          </Group>
          <Group gap="xs">
            <Users size="15" />
            <Text size="sm">{landingPage.members_count}</Text>
          </Group>
        </Group>
      </Card.Section>

      <Anchor
        component={Link}
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
