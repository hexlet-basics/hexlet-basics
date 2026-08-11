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
import { Link } from "@tanstack/react-router";
import type { PropsWithChildren } from "react";
import type { CourseCatalogItem } from "@/client/types.gen";
import { HoverLift } from "./HoverLift";

type Props = PropsWithChildren &
  CardProps & {
    item: CourseCatalogItem;
    lazy?: boolean;
  };

// Reused from legacy CourseBlock; data source rewired from Inertia props to the
// generated API type, and navigation from Inertia Link to TanStack Router Link.
export default function CourseBlock({ item, lazy, ...props }: Props) {
  return (
    <HoverLift h="100%">
      <Card
        renderRoot={(rootProps) => (
          <Link to="/languages/$slug" params={{ slug: item.slug }} {...rootProps} />
        )}
        shadow="sm"
        td="none"
        {...props}
      >
        <Card.Section>
          <AspectRatio ratio={4 / 3}>
            <Image
              fit="cover"
              loading={lazy ? "lazy" : "eager"}
              src={item.coverUrl}
              alt={item.header ?? item.name ?? item.course.slug}
              fallbackSrc="https://placehold.co/400x300?text=Course"
            />
          </AspectRatio>
        </Card.Section>

        <Stack pt="md" h="100%">
          <Title fw="bold" order={2} mb="md">
            {item.name}
          </Title>
          <Group c="dimmed" mt="auto">
            <Group gap="xs">
              <IconClock size="15" />
              <Text size="sm">{item.duration}</Text>
            </Group>
            <Group gap="xs">
              <IconUsers size="15" />
              <Text size="sm">
                <NumberFormatter thousandSeparator value={item.enrollmentsCount} />
              </Text>
            </Group>
          </Group>
        </Stack>
      </Card>
    </HoverLift>
  );
}
