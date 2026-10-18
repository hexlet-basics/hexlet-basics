import { Group, Stack, TextInput, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminDeleteReviewMutation,
  adminListReviewsOptions,
  adminListReviewsQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zAdminListReviewsQuery } from "@/client/zod.gen";
import type { Review } from "@/client/types.gen";
import { EditLink } from "@/components/admin/EditLink";
import { ButtonLink } from "@/components/RouterLink";
import { DataTable } from "mantine-datatable";
import type { DataTableColumn } from "mantine-datatable";
import { RowActions } from "@/components/admin/RowActions";
import { SEARCH_DEBOUNCE_MS } from "@/lib/duration";
import { useState } from "react";
import { useDebouncedCallback } from "@mantine/hooks";
import { useQuery } from "@tanstack/react-query";
import type { ReactElement } from "react";

export const Route = createFileRoute("/{-$locale}/admin/reviews/")({
  component: ReviewsList,
  validateSearch: zAdminListReviewsQuery,
});

function ReviewsList(): ReactElement {
  const { t } = useTranslation();
  const search = Route.useSearch();
  const navigate = Route.useNavigate();
  const { data, isFetching } = useQuery(
    adminListReviewsOptions({
      query: { page: search.page ?? 1, perPage: search.perPage, q: search.q },
    }),
  );
  const [term, setTerm] = useState(search.q ?? "");
  const applyTerm = useDebouncedCallback((text: string) => {
    void navigate({
      search: (prev) => ({ ...prev, page: undefined, q: text === "" ? undefined : text }),
    });
  }, SEARCH_DEBOUNCE_MS);

  const columns: DataTableColumn<Review>[] = [
    {
      accessor: "fullName",
      title: t(($) => $.models.attributes.review.first_name),
    },
    {
      accessor: "course",
      render: (record) => record.course.slug,
      title: t(($) => $.models.attributes.review.course_id),
    },
    {
      accessor: "state",
      render: (record) => {
        const { state } = record;
        return state ? t(($) => $.models.attributes.review["state/values"][state]) : "";
      },
      title: t(($) => $.models.attributes.review.state),
    },
    {
      accessor: "locale",
      title: t(($) => $.models.attributes.review.locale),
    },
    {
      accessor: "body",
      title: t(($) => $.models.attributes.review.body),
    },
    {
      accessor: "actions",
      render: (record) => (
        <RowActions
          describe={record.body ?? ""}
          id={record.id}
          listKey={adminListReviewsQueryKey()}
          remove={adminDeleteReviewMutation()}
        >
          <EditLink to="/{-$locale}/admin/reviews/$id" params={{ id: String(record.id) }} />
        </RowActions>
      ),
      textAlign: "right",
      title: "",
    },
  ];

  return (
    <Stack>
      <Group justify="space-between">
        <Title order={2}>{t(($) => $.admin.resources.reviews)}</Title>
        <ButtonLink to="/{-$locale}/admin/reviews/new">{t(($) => $.admin.crud.new)}</ButtonLink>
      </Group>

      <TextInput
        maw={360}
        placeholder={t(($) => $.admin.crud.search)}
        value={term}
        onChange={(event) => {
          setTerm(event.currentTarget.value);
          applyTerm(event.currentTarget.value.trim());
        }}
      />

      <DataTable
        withTableBorder
        striped
        highlightOnHover
        fetching={isFetching}
        records={data?.items}
        totalRecords={data?.total}
        recordsPerPage={data?.perPage ?? search.perPage}
        page={search.page ?? 1}
        onPageChange={(page) => {
          void navigate({ search: (prev) => ({ ...prev, page }) });
        }}
        noRecordsText={t(($) => $.admin.crud.empty)}
        columns={columns}
      />
    </Stack>
  );
}
