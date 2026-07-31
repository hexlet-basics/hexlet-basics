import { Button, Card, Checkbox, Group, Table, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import {
  adminGetRolePermissionsOptions,
  adminGetRolePermissionsQueryKey,
  adminGetRoleQueryKey,
  adminUpdateRolePermissionsMutation,
} from "@/client/@tanstack/react-query.gen";
import type { PermissionResource, RolePermissionInput } from "@/client/types.gen";
import { useResourceMutation } from "@/hooks/useResourceMutation";

// Every admin surface a role can be scoped to. Mirrors the contract's
// PermissionResource union — a new resource extends both in one PR.
const RESOURCES: PermissionResource[] = [
  "banners",
  "blog_posts",
  "language_categories",
  "language_landing_pages",
  "language_lesson_members",
  "language_lesson_reviews",
  "language_lessons",
  "languages",
  "leads",
  "messages",
  "reviews",
];

const ABILITIES = ["canIndex", "canCreate", "canUpdate", "canDestroy"] as const;

// The editable permission matrix on a role's edit page. The backend upserts
// each submitted row by (role, resource) and leaves unlisted resources alone,
// so submitting the full grid is idempotent.
export function RolePermissionsMatrix({ roleId }: { roleId: number }) {
  const { data } = useQuery(adminGetRolePermissionsOptions({ path: { roleId } }));

  if (!data) return null;

  const initial = new Map(data.permissions.map((p) => [p.resource, p]));
  const rows: RolePermissionInput[] = RESOURCES.map((resource) => ({
    resource,
    canIndex: initial.get(resource)?.canIndex ?? false,
    canCreate: initial.get(resource)?.canCreate ?? false,
    canUpdate: initial.get(resource)?.canUpdate ?? false,
    canDestroy: initial.get(resource)?.canDestroy ?? false,
  }));

  // Remount the grid when a fresh matrix arrives so its state re-seeds.
  return <MatrixForm key={JSON.stringify(rows)} roleId={roleId} initialRows={rows} />;
}

function MatrixForm({
  roleId,
  initialRows,
}: {
  roleId: number;
  initialRows: RolePermissionInput[];
}) {
  const { t } = useTranslation();
  const [rows, setRows] = useState(initialRows);

  const mutation = useResourceMutation({
    mutation: adminUpdateRolePermissionsMutation(),
    invalidate: [
      adminGetRolePermissionsQueryKey({ path: { roleId } }),
      adminGetRoleQueryKey({ path: { id: roleId } }),
    ],
    successMessage: t(($) => $.admin.crud.updated),
    errorMessage: t(($) => $.admin.crud.saveError),
  });

  const toggle = (resource: PermissionResource, ability: (typeof ABILITIES)[number]) => {
    setRows((current) =>
      current.map((row) =>
        row.resource === resource ? { ...row, [ability]: !row[ability] } : row,
      ),
    );
  };

  return (
    <Card withBorder p="xl" maw={720}>
      <Title order={4} mb="md">
        {t(($) => $.admin.permissions.title)}
      </Title>
      <Table>
        <Table.Thead>
          <Table.Tr>
            <Table.Th />
            <Table.Th>{t(($) => $.admin.permissions.canIndex)}</Table.Th>
            <Table.Th>{t(($) => $.admin.permissions.canCreate)}</Table.Th>
            <Table.Th>{t(($) => $.admin.permissions.canUpdate)}</Table.Th>
            <Table.Th>{t(($) => $.admin.permissions.canDestroy)}</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody>
          {rows.map((row) => (
            <Table.Tr key={row.resource}>
              <Table.Td>{row.resource}</Table.Td>
              {ABILITIES.map((ability) => (
                <Table.Td key={ability}>
                  <Checkbox
                    checked={row[ability]}
                    aria-label={`${row.resource} ${ability}`}
                    onChange={() => toggle(row.resource, ability)}
                  />
                </Table.Td>
              ))}
            </Table.Tr>
          ))}
        </Table.Tbody>
      </Table>
      <Group justify="flex-end" mt="md">
        <Button
          loading={mutation.isPending}
          onClick={() => mutation.mutate({ path: { roleId }, body: { permissions: rows } })}
        >
          {t(($) => $.admin.crud.save)}
        </Button>
      </Group>
    </Card>
  );
}
