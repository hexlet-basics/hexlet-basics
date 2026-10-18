import { Button, Group, Stack, TextInput, Title } from "@mantine/core";
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
import { EditLink } from "@/components/admin/EditLink";
import { ButtonLink } from "@/components/RouterLink";
import { useResourceMutation } from "@/hooks/useResourceMutation";
import { DataTable } from "mantine-datatable";
import type { DataTableColumn } from "mantine-datatable";
import { RowActions } from "@/components/admin/RowActions";
import { SEARCH_DEBOUNCE_MS } from "@/lib/duration";
import { useState } from "react";
import { useDebouncedCallback } from "@mantine/hooks";
import { useQuery } from "@tanstack/react-query";
import type { ReactElement } from "react";

export const Route = createFileRoute("/{-$locale}/admin/blog_posts/")({
  component: BlogPostsList,
  validateSearch: zAdminListBlogPostsQuery,
});

function BlogPostsList(): ReactElement {
  const { t } = useTranslation();
  const search = Route.useSearch();
  const navigate = Route.useNavigate();
  const { data, isFetching } = useQuery(
    adminListBlogPostsOptions({
      query: { page: search.page ?? 1, perPage: search.perPage, q: search.q },
    }),
  );
  const [term, setTerm] = useState(search.q ?? "");
  const applyTerm = useDebouncedCallback((text: string) => {
    void navigate({
      search: (prev) => ({ ...prev, page: undefined, q: text === "" ? undefined : text }),
    });
  }, SEARCH_DEBOUNCE_MS);

  // Legacy row action: enqueue the AI pick of the post's related courses.
  const suggestMutation = useResourceMutation({
    errorMessage: t(($) => $.admin.crud.saveError),
    invalidate: [adminListBlogPostsQueryKey()],
    mutation: adminSuggestBlogPostRelatedCoursesMutation(),
    successMessage: t(($) => $.admin.blogPosts.suggestQueued),
  });

  const columns: DataTableColumn<BlogPost>[] = [
    {
      accessor: "name",
      title: t(($) => $.models.attributes.base.name),
    },
    {
      accessor: "slug",
      title: t(($) => $.models.attributes.base.slug),
    },
    {
      accessor: "state",
      render: (record) => {
        const { state } = record;
        return state ? t(($) => $.models.attributes.blog_post["state/values"][state]) : "";
      },
      title: t(($) => $.models.attributes.blog_post.state),
    },
    {
      accessor: "locale",
      title: t(($) => $.models.attributes.course_category.locale),
    },
    {
      accessor: "likesCount",
      title: t(($) => $.admin.blogPosts.likes),
    },
    {
      accessor: "actions",
      render: (record) => (
        <RowActions
          describe={record.name ?? ""}
          id={record.id}
          listKey={adminListBlogPostsQueryKey()}
          remove={adminDeleteBlogPostMutation()}
        >
          <EditLink to="/{-$locale}/admin/blog_posts/$id" params={{ id: String(record.id) }} />
          <Button
            size="xs"
            variant="light"
            color="grape"
            loading={suggestMutation.isPending && suggestMutation.variables?.path.id === record.id}
            onClick={() => {
              suggestMutation.mutate({ path: { id: record.id } });
            }}
          >
            {t(($) => $.admin.blogPosts.suggestRelatedCourses)}
          </Button>
        </RowActions>
      ),
      textAlign: "right",
      title: "",
    },
  ];

  return (
    <Stack>
      <Group justify="space-between">
        <Title order={2}>{t(($) => $.admin.resources.blogPosts)}</Title>
        <ButtonLink to="/{-$locale}/admin/blog_posts/new">{t(($) => $.admin.crud.new)}</ButtonLink>
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
