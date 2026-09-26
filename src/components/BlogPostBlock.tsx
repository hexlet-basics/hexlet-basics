import { AspectRatio, Card, Group, Image, Text, Title } from "@mantine/core";
import { IconClock, IconThumbUp } from "@tabler/icons-react";
import { useTranslation } from "react-i18next";
import type { BlogPost } from "@/client/types.gen";
import { CardLink } from "@/components/RouterLink";
import { timeAgo } from "@/lib/time-ago";
import { HoverLift } from "./HoverLift";

// Blog post card, ported from legacy BlogPostBlock. Reading time is the legacy
// hardcoded five minutes: the serializer's readingTime was never shown.
export default function BlogPostBlock({ post, lazy }: { post: BlogPost; lazy?: boolean }) {
  const { t, i18n } = useTranslation();

  return (
    <HoverLift h="100%">
      <CardLink
        to="/{-$locale}/blog_posts/$slug"
        params={{ slug: post.slug ?? "" }}
        shadow="sm"
        radius="md"
        h="100%"
        td="none"
      >
        {post.coverListVariant && (
          <Card.Section>
            <AspectRatio ratio={2 / 1}>
              <Image
                fit="cover"
                loading={lazy ? "lazy" : "eager"}
                src={post.coverListVariant}
                alt={`Cover for ${post.name}`}
              />
            </AspectRatio>
          </Card.Section>
        )}
        <Title order={2} fz="h5" fw="bold" my="md">
          {post.name}
        </Title>
        <Text mb="xs">{post.description}</Text>
        <Group gap="xs" c="dimmed" mt="auto">
          <Text me="auto">{timeAgo(post.createdAt, i18n.language)}</Text>
          <Group gap={5}>
            <IconThumbUp size={14} />
            <Text>{post.likesCount}</Text>
          </Group>
          <Group gap={5}>
            <IconClock size={14} />
            <Text>~{t(($) => $.common.time.minutes, { count: 5 })}</Text>
          </Group>
        </Group>
      </CardLink>
    </HoverLift>
  );
}
