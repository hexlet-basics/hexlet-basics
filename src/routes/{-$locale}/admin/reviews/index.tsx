import { ActionIcon, Group, Stack, Title, Tooltip } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminDeleteReviewMutation,
  adminListReviewsOptions,
  adminListReviewsQueryKey,
} from "@/client/@tanstack/react-query.gen";
import type { Review } from "@/client/types.gen";
import { type CrudColumnDef, CrudList } from "@/components/admin/CrudList";
import { useCrudPager } from "@/components/admin/useCrudPager";
import { ButtonLink } from "@/components/RouterLink";
import { useDeleteConfirmation } from "@/hooks/useDeleteConfirmation";
import { useResourceMutation } from "@/hooks/useResourceMutation";

export const Route = createFileRoute("/{-$locale}/admin/reviews/")({
  component: ReviewsList,
});

function ReviewsList() {
  const { t } = useTranslation();
  const confirmDelete = useDeleteConfirmation();

  const pager = useCrudPager();

  const { data, isLoading } = useQuery(
    adminListReviewsOptions({
      query: pager.query,
    }),
  );

  const deleteMutation = useResourceMutation({
    mutation: adminDeleteReviewMutation(),
    invalidate: [adminListReviewsQueryKey()],
    successMessage: t(($) => $.admin.crud.deleted),
    errorMessage: t(($) => $.admin.crud.deleteError),
  });

  const columns: CrudColumnDef<Review>[] = [
    {
      accessorKey: "fullName",
      header: t(($) => $.models.attributes.review.first_name),
    },
    {
      id: "course",
      header: t(($) => $.models.attributes.review.course_id),
      cell: ({ row }) => row.original.course.slug,
    },
    {
      accessorKey: "state",
      header: t(($) => $.models.attributes.review.state),
    },
    {
      accessorKey: "locale",
      header: t(($) => $.models.attributes.review.locale),
    },
    {
      accessorKey: "body",
      header: t(($) => $.models.attributes.review.body),
    },
    {
      id: "actions",
      header: "",
      cell: ({ row }) => (
        <Group gap="xs" justify="flex-end" wrap="nowrap">
          <ButtonLink
            to="/{-$locale}/admin/reviews/$id"
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
                  description: row.original.body ?? "",
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
        <Title order={2}>{t(($) => $.admin.resources.reviews)}</Title>
        <ButtonLink to="/{-$locale}/admin/reviews/new">{t(($) => $.admin.crud.new)}</ButtonLink>
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
