import {
  AspectRatio,
  Card,
  type CardProps,
  Group,
  Image,
  NumberFormatter,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { IconClock, IconUsers } from "@tabler/icons-react";
import type { PropsWithChildren } from "react";
import * as Routes from "@/routes.js";
import type {
  LanguageLandingPageForLists,
  LanguageMember,
} from "@/types/serializers";
import AppAnchor from "./Elements/AppAnchor";
import { HoverLift } from "./HoverLift";

type Props = PropsWithChildren &
  CardProps & {
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
  ...props
}: Props) {
  return (
    <AppAnchor href={Routes.language_url(landingPage.slug)} td="none">
      <HoverLift h="100%">
        <Card shadow="sm" {...props}>
          <Card.Section>
            <AspectRatio ratio={4 / 3}>
              <Image
                fit="cover"
                loading={lazy ? "lazy" : "eager"}
                src={landingPage.language.cover_list_variant}
                alt={landingPage.header}
              />
            </AspectRatio>
          </Card.Section>

          <Stack pt="md" h="100%">
            <Title fw="bold" order={2} mb="md">
              {landingPage.name}
            </Title>
            <Group c="dimmed" mt="auto">
              <Group gap="xs">
                <IconClock size="15" />
                <Text size="sm">{landingPage.duration}</Text>
              </Group>
              <Group gap="xs">
                <IconUsers size="15" />
                <Text size="sm">
                  <NumberFormatter
                    thousandSeparator
                    value={landingPage.members_count}
                  />
                </Text>
              </Group>
            </Group>
          </Stack>
        </Card>
      </HoverLift>
    </AppAnchor>
    //     {
    //   continueButton && (
    //     <Button variant="light" fullWidth>
    //       {tHelpers("continue")}
    //     </Button>
    //   )
    // }
  );
}
