import { Link, usePage } from "@inertiajs/react";
import {
  Anchor,
  AspectRatio,
  Card,
  Group,
  Image,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { IconClock, IconThumbUp } from "@tabler/icons-react";
import dayjs from "dayjs";
import type { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import * as Routes from "@/routes.js";
import type { BlogPost } from "@/types/serializers";
import AppAnchor from "./Elements/AppAnchor";
import { HoverLift } from "./HoverLift";

type Props = PropsWithChildren & {
  post: BlogPost;
  lazy?: boolean;
};

export default function BlogPostBlock({ post, lazy }: Props) {
  const { suffix } = usePage().props;
  const { t } = useTranslation();

  return (
    <AppAnchor href={Routes.blog_post_path(post.slug!, { suffix })} td="none">
      <HoverLift h="100%">
        <Card shadow="sm" radius="md" h="100%">
          {post.cover_list_variant && (
            <Card.Section>
              <AspectRatio ratio={2 / 1}>
                <Image
                  fit="cover"
                  loading={lazy ? "lazy" : "eager"}
                  src={post.cover_list_variant!}
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
            <Text me="auto" mb="xs">
              {dayjs().to(post.created_at)}
            </Text>
            <Group gap={5}>
              <IconThumbUp size={14} />
              <Text>{post.likes_count}</Text>
            </Group>
            <Group gap={5}>
              <IconClock size={14} />
              <Text>
                ~
                {t(($) => $.common.time.minutes, {
                  count: 5,
                })}
              </Text>
            </Group>
          </Group>
        </Card>
      </HoverLift>
    </AppAnchor>
  );
}
