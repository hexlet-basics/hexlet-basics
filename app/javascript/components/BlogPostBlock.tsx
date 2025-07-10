import { Link, usePage } from '@inertiajs/react';
import { Anchor, Card, Group, Image, Stack, Text, Title } from '@mantine/core';
import dayjs from 'dayjs';
import type { PropsWithChildren } from 'react';
import { useTranslation } from 'react-i18next';
import * as Routes from '@/routes.js';
import type { SharedProps } from '@/types';
import type { BlogPost } from '@/types/serializers';
import AppAnchor from './AppAnchor';

type Props = PropsWithChildren & {
  post: BlogPost;
  lazy?: boolean;
};

export default function BlogPostBlock({ post, lazy }: Props) {
  const { suffix } = usePage<SharedProps>().props;
  const { t: tCommon } = useTranslation('common');

  return (
    <Card pos="relative" shadow="sm" padding="lg" radius="md" withBorder>
      {post.cover_thumb_variant && (
        <Card.Section>
          <Image
            loading={lazy ? 'lazy' : 'eager'}
            src={post.cover_list_variant!}
            alt={`Cover for ${post.name}`}
          />
        </Card.Section>
      )}

      <Title order={2} my="md">
        {post.name}
      </Title>
      <Text mb="xs">{post.description}</Text>

      <Group gap="xs" c="dimmed" mt="auto">
        <Text>{dayjs().to(post.created_at)}</Text>
        <Group gap="xs">
          <i className="bi bi-hand-thumbs-up" />
          <Text>{post.likes_count}</Text>
        </Group>
        <Group gap="xs">
          <i className="bi bi-clock" />
          <Text>~{tCommon('time.minutes', { count: 5 })}</Text>
        </Group>
      </Group>

      <AppAnchor
        pos="absolute"
        inset={0}
        href={Routes.blog_post_path(post.slug!, { suffix })}
      ></AppAnchor>
    </Card>
  );
}
