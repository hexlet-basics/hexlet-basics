import { Badge, Button, Card, Group, Title } from "@mantine/core";
import { useTranslation } from "react-i18next";
import {
  adminArchiveBlogPostMutation,
  adminGetBlogPostQueryKey,
  adminListBlogPostsQueryKey,
  adminPublishBlogPostMutation,
  adminUnpublishBlogPostMutation,
} from "@/client/@tanstack/react-query.gen";
import type { BlogPost, BlogPostState } from "@/client/types.gen";
import { transitionsFrom } from "@/components/admin/resources/blogPost";
import { useResourceMutation } from "@/hooks/useResourceMutation";

// The Publication State of a post and the buttons that change it. A state is not
// an update: the form does not carry it, each button calls its own operation
// (its own Permission on the server). Only the moves legal for the current state
// are offered; the server stays the judge and refuses the rest with a 409.
export function BlogPostStatePanel({ post }: { post: BlogPost }) {
  const { t } = useTranslation();
  const invalidate = [
    adminListBlogPostsQueryKey(),
    adminGetBlogPostQueryKey({ path: { id: post.id } }),
  ];
  const options = {
    invalidate,
    successMessage: t(($) => $.admin.blogPosts.stateChanged),
    errorMessage: t(($) => $.admin.blogPosts.stateChangeError),
  };

  const publish = useResourceMutation({ mutation: adminPublishBlogPostMutation(), ...options });
  const unpublish = useResourceMutation({ mutation: adminUnpublishBlogPostMutation(), ...options });
  const archive = useResourceMutation({ mutation: adminArchiveBlogPostMutation(), ...options });

  const state: BlogPostState = post.state ?? "draft";
  const actions = {
    publish: { label: t(($) => $.admin.blogPosts.publish), mutation: publish, color: "green" },
    unpublish: { label: t(($) => $.admin.blogPosts.unpublish), mutation: unpublish, color: "gray" },
    archive: { label: t(($) => $.admin.blogPosts.archive), mutation: archive, color: "red" },
  };

  return (
    <Card withBorder p="xl" maw={720}>
      <Group justify="space-between">
        <Group>
          <Title order={4}>{t(($) => $.models.attributes.blog_post.state)}</Title>
          <Badge variant="light">
            {t(($) => $.models.attributes.blog_post["state/values"][state])}
          </Badge>
        </Group>
        <Group>
          {transitionsFrom(state).map((name) => (
            <Button
              key={name}
              variant="light"
              color={actions[name].color}
              loading={actions[name].mutation.isPending}
              onClick={() => actions[name].mutation.mutate({ path: { id: post.id } })}
            >
              {actions[name].label}
            </Button>
          ))}
        </Group>
      </Group>
    </Card>
  );
}
