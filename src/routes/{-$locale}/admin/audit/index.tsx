import {
  Alert,
  Badge,
  Button,
  Code,
  Group,
  Modal,
  NativeSelect,
  Stack,
  Table,
  Text,
  Title,
} from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import type { OnChangeFn, PaginationState } from "@tanstack/react-table";
import { useState } from "react";
import { z } from "zod";
import { useTranslation } from "react-i18next";
import dayjs from "dayjs";
import { adminListAuditRecordsOptions } from "@/client/@tanstack/react-query.gen";
import type { AuditRecord } from "@/client/types.gen";
import { zAdminListAuditRecordsQuery } from "@/client/zod.gen";
import { type CrudColumnDef, CrudList } from "@/components/admin/CrudList";
import { submitForm, useAppForm } from "@/lib/form";
import { isApiError } from "@/lib/api-error";
import { validationMessages } from "@/lib/server-validation";

// The audit trail (ADR-0020): every change a staff member made and every read of
// personal data, newest first. The API needs the capability Permission
// `audit.read` on the caller's Staff Role (an administrator included), refuses a
// window wider than 92 days, and records that the trail was read.
export const Route = createFileRoute("/{-$locale}/admin/audit/")({
  // The generated query schema of the list operation: the page, the filters and
  // the window live in the address, so a filtered view can be shared and survives
  // a reload.
  validateSearch: zAdminListAuditRecordsQuery,
  component: AuditTrail,
});

type AuditSearch = z.input<typeof zAdminListAuditRecordsQuery>;

const PER_PAGE = 25;

// What the filter bar holds before it is applied: the raw text of the inputs,
// and for the window the day pickers' `YYYY-MM-DD` strings.
const NO_FILTERS = {
  actorUserId: "",
  channel: "",
  operation: "",
  entityType: "",
  entityId: "",
  from: "",
  to: "",
};
type Filters = typeof NO_FILTERS;

// A picked day starts at the visitor's local midnight, so "from 5 May" includes
// the morning of the 5th wherever they are.
const startOfDay = (day: string) => dayjs(day).toISOString();

// The API's `to` is exclusive, so the last day picked is included by ending the
// window at the next local midnight.
const endOfDay = (day: string) => dayjs(day).add(1, "day").toISOString();

const filtersOf = (search: AuditSearch): Filters => ({
  actorUserId: search.actorUserId?.toString() ?? "",
  channel: search.channel ?? "",
  operation: search.operation ?? "",
  entityType: search.entityType ?? "",
  entityId: search.entityId ?? "",
  from: search.from ? dayjs(search.from).format("YYYY-MM-DD") : "",
  to: search.to ? dayjs(search.to).subtract(1, "day").format("YYYY-MM-DD") : "",
});

// A blank input is no filter: it stays out of the address and the request.
const searchOf = (filters: Filters): AuditSearch => ({
  actorUserId: filters.actorUserId === "" ? undefined : Number(filters.actorUserId),
  // "" (any channel) is no member of the enum, so it parses to nothing.
  channel: zAdminListAuditRecordsQuery.shape.channel.safeParse(filters.channel).data,
  operation: filters.operation || undefined,
  entityType: filters.entityType || undefined,
  entityId: filters.entityId || undefined,
  from: filters.from === "" ? undefined : startOfDay(filters.from),
  to: filters.to === "" ? undefined : endOfDay(filters.to),
});

const CHANNEL_COLOR = { web: "blue", mcp: "grape" } as const;
const OUTCOME_COLOR = { succeeded: "green", forbidden: "orange", failed: "red" } as const;

function AuditTrail() {
  const { t } = useTranslation();
  const search = Route.useSearch();
  const navigate = Route.useNavigate();
  const [selected, setSelected] = useState<AuditRecord | null>(null);

  const pagination: PaginationState = {
    pageIndex: (search.page ?? 1) - 1,
    pageSize: search.perPage ?? PER_PAGE,
  };
  const onPaginationChange: OnChangeFn<PaginationState> = (updater) => {
    const next = typeof updater === "function" ? updater(pagination) : updater;
    void navigate({
      search: (prev) => ({ ...prev, page: next.pageIndex + 1, perPage: next.pageSize }),
    });
  };

  const { data, isLoading, error } = useQuery(
    adminListAuditRecordsOptions({ query: { ...search, perPage: pagination.pageSize } }),
  );

  // Applying starts from the first page: the old one may not exist in the narrowed trail.
  const apply = (filters: Filters) =>
    navigate({ search: { ...searchOf(filters), perPage: search.perPage } });

  const form = useAppForm({
    defaultValues: filtersOf(search),
    onSubmit: ({ value }) => apply(value),
  });
  const reset = () => {
    form.reset(NO_FILTERS);
    void apply(NO_FILTERS);
  };

  const columns: CrudColumnDef<AuditRecord>[] = [
    {
      id: "occurredAt",
      header: t(($) => $.admin.audit.columns.occurredAt),
      cell: ({ row }) => new Date(row.original.occurredAt).toLocaleString(),
    },
    {
      id: "channel",
      header: t(($) => $.admin.audit.columns.channel),
      cell: ({ row }) => (
        <Badge color={CHANNEL_COLOR[row.original.channel]} variant="light">
          {t(($) => $.admin.audit.channel[row.original.channel])}
        </Badge>
      ),
    },
    {
      id: "outcome",
      header: t(($) => $.admin.audit.columns.outcome),
      cell: ({ row }) => (
        <Badge color={OUTCOME_COLOR[row.original.outcome]}>
          {t(($) => $.admin.audit.outcome[row.original.outcome])}
        </Badge>
      ),
    },
    {
      accessorKey: "operation",
      header: t(($) => $.admin.audit.columns.operation),
    },
    {
      id: "entity",
      header: t(($) => $.admin.audit.columns.entity),
      cell: ({ row }) =>
        [row.original.entityType, row.original.entityId].filter(Boolean).join(" #"),
    },
    {
      id: "actor",
      header: t(($) => $.admin.audit.columns.actor),
      cell: ({ row }) => row.original.actorUserId ?? "—",
    },
    {
      id: "details",
      header: "",
      cell: ({ row }) => (
        <Button variant="subtle" size="compact-sm" onClick={() => setSelected(row.original)}>
          {t(($) => $.admin.audit.details.open)}
        </Button>
      ),
    },
  ];

  const forbidden = isApiError(error, 403);
  const problems = validationMessages(error);

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.resources.audit)}</Title>

      <form onSubmit={submitForm(form)}>
        <Stack gap="xs">
          <Group align="flex-end" gap="sm">
            <form.AppField name="from">
              {(field) => (
                <field.TextField label={t(($) => $.admin.audit.filters.from)} type="date" />
              )}
            </form.AppField>
            <form.AppField name="to">
              {(field) => (
                <field.TextField label={t(($) => $.admin.audit.filters.to)} type="date" />
              )}
            </form.AppField>
            <form.AppField name="actorUserId">
              {(field) => (
                <field.TextField
                  label={t(($) => $.admin.audit.filters.actorUserId)}
                  type="number"
                  min={1}
                />
              )}
            </form.AppField>
            <form.Field name="channel">
              {(field) => (
                <NativeSelect
                  label={t(($) => $.admin.audit.filters.channel)}
                  data={[
                    { value: "", label: t(($) => $.admin.audit.filters.anyChannel) },
                    { value: "web", label: t(($) => $.admin.audit.channel.web) },
                    { value: "mcp", label: t(($) => $.admin.audit.channel.mcp) },
                  ]}
                  value={field.state.value}
                  onChange={(event) => field.handleChange(event.currentTarget.value)}
                />
              )}
            </form.Field>
          </Group>
          <Group align="flex-end" gap="sm">
            <form.AppField name="operation">
              {(field) => (
                <field.TextField
                  label={t(($) => $.admin.audit.filters.operation)}
                  placeholder="AdminUpdateBlogPost"
                />
              )}
            </form.AppField>
            <form.AppField name="entityType">
              {(field) => (
                <field.TextField
                  label={t(($) => $.admin.audit.filters.entityType)}
                  placeholder="BlogPost"
                />
              )}
            </form.AppField>
            <form.AppField name="entityId">
              {(field) => <field.TextField label={t(($) => $.admin.audit.filters.entityId)} />}
            </form.AppField>
            <Button type="submit">{t(($) => $.admin.audit.filters.apply)}</Button>
            <Button type="button" variant="default" onClick={reset}>
              {t(($) => $.admin.audit.filters.reset)}
            </Button>
          </Group>
          <Text size="sm" c="dimmed">
            {t(($) => $.admin.audit.filters.windowHint)}
          </Text>
        </Stack>
      </form>

      {forbidden && (
        <Alert color="red" title={t(($) => $.admin.audit.forbidden.title)}>
          {t(($) => $.admin.audit.forbidden.body)}
        </Alert>
      )}
      {problems.length > 0 && (
        <Alert color="orange" role="alert">
          {problems.join(" ")}
        </Alert>
      )}

      <CrudList
        columns={columns}
        data={data?.items ?? []}
        total={data?.total ?? 0}
        pagination={pagination}
        onPaginationChange={onPaginationChange}
        isLoading={isLoading}
      />

      <Modal
        opened={selected !== null}
        onClose={() => setSelected(null)}
        title={selected?.operation}
        size="xl"
      >
        {selected && <RecordDetails record={selected} />}
      </Modal>
    </Stack>
  );
}

// One record in full: where it came from, what it changed (old and new value of
// each field; a secret field names itself and nothing else) and what it read
// (a count and ids, never values).
function RecordDetails({ record }: { record: AuditRecord }) {
  const { t } = useTranslation();
  const meta: [string, string | number | null][] = [
    [t(($) => $.admin.audit.details.requestId), record.requestId],
    [t(($) => $.admin.audit.details.address), record.address],
    [t(($) => $.admin.audit.details.userAgent), record.userAgent],
    [t(($) => $.admin.audit.details.connectedApp), record.connectedAppId],
    [t(($) => $.admin.audit.details.token), record.tokenId],
    [t(($) => $.admin.audit.details.locale), record.locale],
  ];

  return (
    <Stack>
      <Group gap="xs">
        <Badge color={CHANNEL_COLOR[record.channel]} variant="light">
          {t(($) => $.admin.audit.channel[record.channel])}
        </Badge>
        <Badge color={OUTCOME_COLOR[record.outcome]}>
          {t(($) => $.admin.audit.outcome[record.outcome])}
        </Badge>
      </Group>

      <Table withRowBorders={false}>
        <Table.Tbody>
          {meta
            .filter(([, value]) => value !== null)
            .map(([label, value]) => (
              <Table.Tr key={label}>
                <Table.Th w={160}>{label}</Table.Th>
                <Table.Td>{value}</Table.Td>
              </Table.Tr>
            ))}
        </Table.Tbody>
      </Table>

      {record.changes.length > 0 && (
        <Stack gap="xs">
          <Title order={5}>{t(($) => $.admin.audit.details.changes)}</Title>
          {record.changes.map((change, index) => (
            // A record may change the same entity twice, so the position is part of the key.
            <Stack key={`${change.entity}-${change.id ?? ""}-${index}`} gap={4}>
              <Text fw={500}>
                {change.entity}
                {change.id ? ` #${change.id}` : ""} · {change.op}
              </Text>
              {change.fields && change.fields.length > 0 && (
                <Table withTableBorder>
                  <Table.Thead>
                    <Table.Tr>
                      <Table.Th>{t(($) => $.admin.audit.details.field)}</Table.Th>
                      <Table.Th>{t(($) => $.admin.audit.details.old)}</Table.Th>
                      <Table.Th>{t(($) => $.admin.audit.details.new)}</Table.Th>
                    </Table.Tr>
                  </Table.Thead>
                  <Table.Tbody>
                    {change.fields.map((fieldChange) => (
                      <Table.Tr key={fieldChange.name}>
                        <Table.Td>{fieldChange.name}</Table.Td>
                        {fieldChange.secret ? (
                          <Table.Td colSpan={2}>
                            <Text size="sm" c="dimmed">
                              {t(($) => $.admin.audit.details.secret)}
                            </Text>
                          </Table.Td>
                        ) : (
                          <>
                            <Table.Td>
                              <Code>{JSON.stringify(fieldChange.old ?? null)}</Code>
                            </Table.Td>
                            <Table.Td>
                              <Code>{JSON.stringify(fieldChange.new ?? null)}</Code>
                            </Table.Td>
                          </>
                        )}
                      </Table.Tr>
                    ))}
                  </Table.Tbody>
                </Table>
              )}
            </Stack>
          ))}
        </Stack>
      )}

      {record.reads && (
        <Stack gap="xs">
          <Title order={5}>{t(($) => $.admin.audit.details.reads)}</Title>
          <Text>{t(($) => $.admin.audit.details.readCount, { count: record.reads.count })}</Text>
          {record.reads.ids && record.reads.ids.length > 0 && (
            <Code block>{record.reads.ids.join(", ")}</Code>
          )}
        </Stack>
      )}
    </Stack>
  );
}
