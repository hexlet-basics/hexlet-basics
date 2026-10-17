import { ActionIcon, Button, Group, Stack, Title, Tooltip } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
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
import { useCrudPager } from "@/components/admin/useCrudPager";
import { ButtonLink } from "@/components/RouterLink";
import { useDeleteConfirmation } from "@/hooks/useDeleteConfirmation";
import { useResourceMutation } from "@/hooks/useResourceMutation";

export const Route = createFileRoute("/{-$locale}/admin/blog_posts/")({
  component: BlogPostsList,
});

function BlogPostsList() {
  const { t } = useTranslation();
  const confirmDelete = useDeleteConfirmation();

  const pager = useCrudPager();

  const { data, isLoading } = useQuery(
    adminListBlogPostsOptions({
      query: pager.query,
    }),
  );

  const deleteMutation = useResourceMutation({
    mutation: adminDeleteBlogPostMutation(),
    invalidate: [adminListBlogPostsQueryKey()],
    successMessage: t(($) => $.admin.crud.deleted),
    errorMessage: t(($) => $.admin.crud.deleteError),
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
    {
      id: "actions",
      header: "",
      cell: ({ row }) => (
        <Group gap="xs" justify="flex-end" wrap="nowrap">
          <ButtonLink
            to="/{-$locale}/admin/blog_posts/$id"
            params={{ id: String(row.original.id) }}
            size="xs"
            variant="light"
          >
            {t(($) => $.admin.crud.edit)}
          </ButtonLink>
          <Button
            size="xs"
            variant="light"
            color="grape"
            loading={
              suggestMutation.isPending && suggestMutation.variables?.path.id === row.original.id
            }
            onClick={() => suggestMutation.mutate({ path: { id: row.original.id } })}
          >
            {t(($) => $.admin.blogPosts.suggestRelatedCourses)}
          </Button>
          <Tooltip label={t(($) => $.admin.crud.delete)}>
            <ActionIcon
              color="red"
              variant="light"
              aria-label={t(($) => $.admin.crud.delete)}
              onClick={() =>
                confirmDelete({
                  description: row.original.name ?? "",
                  onConfirm: () => deleteMutation.mutate({ path: { id: row.original.id } }),
                })
              }
            >
              ✕
            </ActionIcon>
          </Tooltip>
        </Group>
      ),
    },
  ];

  return (
    <Stack>
      <Group justify="space-between">
        <Title order={2}>{t(($) => $.admin.resources.blogPosts)}</Title>
        <ButtonLink to="/{-$locale}/admin/blog_posts/new">{t(($) => $.admin.crud.new)}</ButtonLink>
      </Group>

      <CrudList
        columns={columns}
        data={data?.items ?? []}
        total={data?.total ?? 0}
        {...pager.props(data?.perPage)}
        isLoading={isLoading}
      />
    </Stack>
  );
}
