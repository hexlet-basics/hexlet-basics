import { Group, Stack, TextInput, Title } from "@mantine/core";
import { useDebouncedCallback } from "@mantine/hooks";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { DataTable } from "mantine-datatable";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import {
  adminDeleteBannerMutation,
  adminListBannersOptions,
  adminListBannersQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zAdminListBannersQuery } from "@/client/zod.gen";
import { EditLink } from "@/components/admin/EditLink";
import { RowActions } from "@/components/admin/RowActions";
import { SEARCH_DEBOUNCE_MS } from "@/lib/duration";
import { ButtonLink } from "@/components/RouterLink";
import type { ReactElement } from "react";

export const Route = createFileRoute("/{-$locale}/admin/banners/")({
  component: BannersList,
  validateSearch: zAdminListBannersQuery,
});

function BannersList(): ReactElement {
  const { t } = useTranslation();
  const search = Route.useSearch();
  const navigate = Route.useNavigate();
  const { data, isFetching } = useQuery(
    adminListBannersOptions({
      query: { page: search.page ?? 1, perPage: search.perPage, q: search.q },
    }),
  );

  const [term, setTerm] = useState(search.q ?? "");
  const applyTerm = useDebouncedCallback((text: string) => {
    void navigate({
      search: (prev) => ({ ...prev, page: undefined, q: text === "" ? undefined : text }),
    });
  }, SEARCH_DEBOUNCE_MS);

  return (
    <Stack>
      <Group justify="space-between">
        <Title order={2}>{t(($) => $.admin.resources.banners)}</Title>
        <ButtonLink to="/{-$locale}/admin/banners/new">{t(($) => $.admin.crud.new)}</ButtonLink>
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
        idAccessor="id"
        fetching={isFetching}
        records={data?.items}
        totalRecords={data?.total}
        recordsPerPage={data?.perPage ?? search.perPage}
        page={search.page ?? 1}
        onPageChange={(page) => {
          void navigate({ search: (prev) => ({ ...prev, page }) });
        }}
        noRecordsText={t(($) => $.admin.crud.empty)}
        columns={[
          {
            accessor: "state",
            render: (banner) => t(($) => $.models.attributes.banner["state/values"][banner.state]),
            title: t(($) => $.models.attributes.banner.state),
          },
          { accessor: "locale", title: t(($) => $.models.attributes.banner.locale) },
          { accessor: "background", title: t(($) => $.models.attributes.banner.background) },
          { accessor: "body", title: t(($) => $.models.attributes.banner.body) },
          {
            accessor: "actions",
            render: (banner) => (
              <RowActions
                describe={banner.body}
                listKey={adminListBannersQueryKey()}
                remove={adminDeleteBannerMutation()}
                id={banner.id}
              >
                <EditLink to="/{-$locale}/admin/banners/$id" params={{ id: String(banner.id) }} />
              </RowActions>
            ),
            textAlign: "right",
            title: "",
          },
        ]}
      />
    </Stack>
  );
}
