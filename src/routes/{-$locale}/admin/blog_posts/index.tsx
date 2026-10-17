import { Button, Group, Stack, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminDeleteBlogPostMutation,
  adminListBlogPostsOptions,
  adminListBlogPostsQueryKey,
  adminSuggestBlogPostRelatedCoursesMutation,
} from "@/client/@tanstack/react-query.gen";
import type { BlogPost } from "@/client/types.gen";
import { type CrudColumnDef, CrudList } from "@/components/admin/CrudList";
import { EditLink } from "@/components/admin/EditLink";
import { useCrudList } from "@/components/admin/useCrudList";
import { ButtonLink } from "@/components/RouterLink";
import { useResourceMutation } from "@/hooks/useResourceMutation";

export const Route = createFileRoute("/{-$locale}/admin/blog_posts/")({
  component: BlogPostsList,
});

function BlogPostsList() {
  const { t } = useTranslation();

  const list = useCrudList({
    list: adminListBlogPostsOptions,
    listKey: adminListBlogPostsQueryKey,
    remove: adminDeleteBlogPostMutation,
  });

  // Legacy row action: enqueue the AI pick of the post's related courses.
  const suggestMutation = useResourceMutation({
    mutation: adminSuggestBlogPostRelatedCoursesMutation(),
    invalidate: [adminListBlogPostsQueryKey()],
    successMessage: t(($) => $.admin.blogPosts.suggestQueued),
    errorMessage: t(($) => $.admin.crud.saveError),
  });

  const columns: CrudColumnDef<BlogPost>[] = [
    {
      accessorKey: "name",
      header: t(($) => $.models.attributes.base.name),
    },
    {
      accessorKey: "slug",
      header: t(($) => $.models.attributes.base.slug),
    },
    {
      accessorKey: "state",
      header: t(($) => $.models.attributes.blog_post.state),
    },
    {
      accessorKey: "locale",
      header: t(($) => $.models.attributes.course_category.locale),
    },
    {
      accessorKey: "likesCount",
      header: t(($) => $.admin.blogPosts.likes),
    },
    list.actionsColumn({
      edit: (row) => (
        <EditLink to="/{-$locale}/admin/blog_posts/$id" params={{ id: String(row.id) }} />
      ),
      extra: (row) => (
        <Button
          size="xs"
          variant="light"
          color="grape"
          loading={suggestMutation.isPending && suggestMutation.variables?.path.id === row.id}
          onClick={() => suggestMutation.mutate({ path: { id: row.id } })}
        >
          {t(($) => $.admin.blogPosts.suggestRelatedCourses)}
        </Button>
      ),
      describe: (row) => row.name ?? "",
    }),
  ];

  return (
    <Stack>
      <Group justify="space-between">
        <Title order={2}>{t(($) => $.admin.resources.blogPosts)}</Title>
        <ButtonLink to="/{-$locale}/admin/blog_posts/new">{t(($) => $.admin.crud.new)}</ButtonLink>
      </Group>

      <CrudList columns={columns} {...list.props} />
    </Stack>
  );
}
