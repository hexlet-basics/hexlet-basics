import { Group, Stack, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminDeleteReviewMutation,
  adminListReviewsOptions,
  adminListReviewsQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zAdminListReviewsQuery } from "@/client/zod.gen";
import type { Review } from "@/client/types.gen";
import { CrudList } from "@/components/admin/CrudList";
import type { CrudColumnDef } from "@/components/admin/CrudList";
import { EditLink } from "@/components/admin/EditLink";
import { useCrudList } from "@/components/admin/useCrudList";
import { ButtonLink } from "@/components/RouterLink";
import type { ReactElement } from "react";

export const Route = createFileRoute("/{-$locale}/admin/reviews/")({
  component: ReviewsList,
  validateSearch: zAdminListReviewsQuery,
});

function ReviewsList(): ReactElement {
  const { t } = useTranslation();

  const list = useCrudList({
    list: adminListReviewsOptions,
    listKey: adminListReviewsQueryKey,
    remove: adminDeleteReviewMutation,
    searchable: true,
  });

  const columns: CrudColumnDef<Review>[] = [
    {
      accessorKey: "fullName",
      header: t(($) => $.models.attributes.review.first_name),
    },
    {
      cell: ({ row }) => row.original.course.slug,
      header: t(($) => $.models.attributes.review.course_id),
      id: "course",
    },
    {
      accessorKey: "state",
      cell: ({ row }) => {
        const { state } = row.original;
        return state ? t(($) => $.models.attributes.review["state/values"][state]) : "";
      },
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
    list.actionsColumn({
      describe: (row) => row.body ?? "",
      edit: (row) => (
        <EditLink to="/{-$locale}/admin/reviews/$id" params={{ id: String(row.id) }} />
      ),
    }),
  ];

  return (
    <Stack>
      <Group justify="space-between">
        <Title order={2}>{t(($) => $.admin.resources.reviews)}</Title>
        <ButtonLink to="/{-$locale}/admin/reviews/new">{t(($) => $.admin.crud.new)}</ButtonLink>
      </Group>

      <CrudList columns={columns} {...list.props} />
    </Stack>
  );
}
