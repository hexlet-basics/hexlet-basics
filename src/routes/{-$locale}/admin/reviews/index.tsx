import { Group, Stack, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminDeleteReviewMutation,
  adminListReviewsOptions,
  adminListReviewsQueryKey,
} from "@/client/@tanstack/react-query.gen";
import type { Review } from "@/client/types.gen";
import { type CrudColumnDef, CrudList } from "@/components/admin/CrudList";
import { EditLink } from "@/components/admin/EditLink";
import { useCrudList } from "@/components/admin/useCrudList";
import { ButtonLink } from "@/components/RouterLink";

export const Route = createFileRoute("/{-$locale}/admin/reviews/")({
  component: ReviewsList,
});

function ReviewsList() {
  const { t } = useTranslation();

  const list = useCrudList({
    list: adminListReviewsOptions,
    listKey: adminListReviewsQueryKey,
    remove: adminDeleteReviewMutation,
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
      cell: ({ row }) => {
        const state = row.original.state;
        return state ? t(($) => $.models.attributes.review["state/values"][state]) : "";
      },
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
      edit: (row) => (
        <EditLink to="/{-$locale}/admin/reviews/$id" params={{ id: String(row.id) }} />
      ),
      describe: (row) => row.body ?? "",
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
