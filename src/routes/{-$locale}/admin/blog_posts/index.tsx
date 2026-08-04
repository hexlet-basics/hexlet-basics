import { ActionIcon, Group, Stack, Title, Tooltip } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import type { PaginationState, SortingState } from "@tanstack/react-table";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import {
  adminDeleteBlogPostMutation,
  adminListBlogPostsOptions,
  adminListBlogPostsQueryKey,
} from "@/client/@tanstack/react-query.gen";
import type { BlogPost } from "@/client/types.gen";
import { type CrudColumnDef, CrudList } from "@/components/admin/CrudList";
import { ButtonLink } from "@/components/RouterLink";
import { useDeleteConfirmation } from "@/hooks/useDeleteConfirmation";
import { useResourceMutation } from "@/hooks/useResourceMutation";

export const Route = createFileRoute("/{-$locale}/admin/blog_posts/")({
  component: BlogPostsList,
});

function BlogPostsList() {
  const { t } = useTranslation();
  const confirmDelete = useDeleteConfirmation();

  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize: 25,
  });
  const [sorting, setSorting] = useState<SortingState>([]);

  const { data, isLoading } = useQuery(
    adminListBlogPostsOptions({
      query: { page: pagination.pageIndex + 1, perPage: pagination.pageSize },
    }),
  );

  const deleteMutation = useResourceMutation({
    mutation: adminDeleteBlogPostMutation(),
    invalidate: [adminListBlogPostsQueryKey()],
    successMessage: t(($) => $.admin.crud.deleted),
    errorMessage: t(($) => $.admin.crud.deleteError),
  });

  const columns: CrudColumnDef<BlogPost>[] = [
    {
      accessorKey: "name",
      header: t(($) => $.models.attributes.base.name),
      enableSorting: false,
    },
    {
      accessorKey: "slug",
      header: t(($) => $.models.attributes.base.slug),
      enableSorting: false,
    },
    {
      accessorKey: "state",
      header: t(($) => $.models.attributes.blog_post.state),
      enableSorting: false,
    },
    {
      accessorKey: "locale",
      header: t(($) => $.models.attributes.language_category.locale),
      enableSorting: false,
    },
    {
      accessorKey: "likesCount",
      header: t(($) => $.admin.blogPosts.likes),
      enableSorting: false,
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
        pagination={pagination}
        onPaginationChange={setPagination}
        sorting={sorting}
        onSortingChange={setSorting}
        isLoading={isLoading}
      />
    </Stack>
  );
}
