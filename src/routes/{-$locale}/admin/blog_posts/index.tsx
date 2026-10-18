import { Button, Group, Stack, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminDeleteBlogPostMutation,
  adminListBlogPostsOptions,
  adminListBlogPostsQueryKey,
  adminSuggestBlogPostRelatedCoursesMutation,
} from "@/client/@tanstack/react-query.gen";
import { zAdminListBlogPostsQuery } from "@/client/zod.gen";
import type { BlogPost } from "@/client/types.gen";
import { CrudList } from "@/components/admin/CrudList";
import type { CrudColumnDef } from "@/components/admin/CrudList";
import { EditLink } from "@/components/admin/EditLink";
import { useCrudList } from "@/components/admin/useCrudList";
import { ButtonLink } from "@/components/RouterLink";
import { useResourceMutation } from "@/hooks/useResourceMutation";
import type { ReactElement } from "react";

export const Route = createFileRoute("/{-$locale}/admin/blog_posts/")({
  validateSearch: zAdminListBlogPostsQuery,
  component: BlogPostsList,
});

function BlogPostsList(): ReactElement {
  const { t } = useTranslation();

  const list = useCrudList({
    list: adminListBlogPostsOptions,
    listKey: adminListBlogPostsQueryKey,
    remove: adminDeleteBlogPostMutation,
    searchable: true,
  });

  // Legacy row action: enqueue the AI pick of the post's related courses.
  const suggestMutation = useResourceMutation({
    errorMessage: t(($) => $.admin.crud.saveError),
    invalidate: [adminListBlogPostsQueryKey()],
    mutation: adminSuggestBlogPostRelatedCoursesMutation(),
    successMessage: t(($) => $.admin.blogPosts.suggestQueued),
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
      cell: ({ row }) => {
        const { state } = row.original;
        return state ? t(($) => $.models.attributes.blog_post["state/values"][state]) : "";
      },
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
      describe: (row) => row.name ?? "",
      edit: (row) => (
        <EditLink to="/{-$locale}/admin/blog_posts/$id" params={{ id: String(row.id) }} />
      ),
      extra: (row) => (
        <Button
          size="xs"
          variant="light"
          color="grape"
          loading={suggestMutation.isPending && suggestMutation.variables?.path.id === row.id}
          onClick={() => {
            suggestMutation.mutate({ path: { id: row.id } });
          }}
        >
          {t(($) => $.admin.blogPosts.suggestRelatedCourses)}
        </Button>
      ),
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
