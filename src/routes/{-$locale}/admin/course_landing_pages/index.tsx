import { Group, Stack, TextInput, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminDeleteCourseLandingPageMutation,
  adminListCourseLandingPagesOptions,
  adminListCourseLandingPagesQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zAdminListCourseLandingPagesQuery } from "@/client/zod.gen";
import type { CourseLandingPage } from "@/client/types.gen";
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

export const Route = createFileRoute("/{-$locale}/admin/course_landing_pages/")({
  component: CourseLandingPagesList,
  validateSearch: zAdminListCourseLandingPagesQuery,
});

function CourseLandingPagesList(): ReactElement {
  const { t } = useTranslation();
  const search = Route.useSearch();
  const navigate = Route.useNavigate();
  const { data, isFetching } = useQuery(
    adminListCourseLandingPagesOptions({
      query: { page: search.page ?? 1, perPage: search.perPage, q: search.q },
    }),
  );
  const [term, setTerm] = useState(search.q ?? "");
  const applyTerm = useDebouncedCallback((text: string) => {
    void navigate({
      search: (prev) => ({ ...prev, page: undefined, q: text === "" ? undefined : text }),
    });
  }, SEARCH_DEBOUNCE_MS);

  const columns: DataTableColumn<CourseLandingPage>[] = [
    {
      accessor: "name",
      title: t(($) => $.models.attributes.course_landing_page.name),
    },
    {
      accessor: "slug",
      title: t(($) => $.models.attributes.course_landing_page.slug),
    },
    {
      accessor: "courseSlug",
      title: t(($) => $.models.attributes.course_landing_page.course_id),
    },
    {
      accessor: "state",
      render: (record) => {
        const { state } = record;
        return state
          ? t(($) => $.models.attributes.course_landing_page["state/values"][state])
          : "";
      },
      title: t(($) => $.models.attributes.course_landing_page.state),
    },
    {
      accessor: "main",
      render: (record) => (record.main ? "✓" : ""),
      title: t(($) => $.models.attributes.course_landing_page.main),
    },
    {
      accessor: "actions",
      render: (record) => (
        <RowActions
          describe={record.name}
          id={record.id}
          listKey={adminListCourseLandingPagesQueryKey()}
          remove={adminDeleteCourseLandingPageMutation()}
        >
          <EditLink
            to="/{-$locale}/admin/course_landing_pages/$id"
            params={{ id: String(record.id) }}
          />
        </RowActions>
      ),
      textAlign: "right",
      title: "",
    },
  ];

  return (
    <Stack>
      <Group justify="space-between">
        <Title order={2}>{t(($) => $.admin.resources.landingPages)}</Title>
        <ButtonLink to="/{-$locale}/admin/course_landing_pages/new">
          {t(($) => $.admin.crud.new)}
        </ButtonLink>
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
