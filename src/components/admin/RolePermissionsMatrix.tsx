import { Badge, Button, Card, Checkbox, Group, Stack, Table, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import {
  adminGetRolePermissionsOptions,
  adminGetRolePermissionsQueryKey,
  adminGetRoleQueryKey,
  adminListPermissionsOptions,
  adminUpdateRolePermissionsMutation,
} from "@/client/@tanstack/react-query.gen";
import type { Permission, PermissionResource, RolePermissionInput } from "@/client/types.gen";
import { useResourceMutation } from "@/hooks/useResourceMutation";

// Every admin surface a role can be scoped to. Mirrors the contract's
// PermissionResource union — a new resource extends both in one PR.
const RESOURCES: PermissionResource[] = [
  "banners",
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

// The editable permissions on a role's edit page: the declared Permissions the
// contract's operations require (listed from the catalog, ADR-0017) and, for the
// resources not yet moved over, the create/read/update/delete matrix. One save
// sends both. The backend upserts each matrix row by (role, resource) and
// leaves unlisted resources alone, and replaces the role's grants with the set
// submitted, so submitting the full screen is idempotent.
export function RolePermissionsMatrix({ roleId }: { roleId: number }) {
  const { data } = useQuery(adminGetRolePermissionsOptions({ path: { roleId } }));
  const { data: catalog } = useQuery(adminListPermissionsOptions());

  if (!data || !catalog) return null;

  const initial = new Map(data.permissions.map((p) => [p.resource, p]));
  const rows: RolePermissionInput[] = RESOURCES.map((resource) => ({
    resource,
    canIndex: initial.get(resource)?.canIndex ?? false,
    canCreate: initial.get(resource)?.canCreate ?? false,
    canUpdate: initial.get(resource)?.canUpdate ?? false,
    canDestroy: initial.get(resource)?.canDestroy ?? false,
  }));

  // Remount the grid when a fresh matrix arrives so its state re-seeds.
  return (
    <MatrixForm
      key={JSON.stringify([rows, data.grants])}
      roleId={roleId}
      initialRows={rows}
      catalog={catalog.permissions}
      initialGrants={data.grants}
    />
  );
}

function MatrixForm({
  roleId,
  initialRows,
  catalog,
  initialGrants,
}: {
  roleId: number;
  initialRows: RolePermissionInput[];
  catalog: Permission[];
  initialGrants: string[];
}) {
  const { t } = useTranslation();
  const [rows, setRows] = useState(initialRows);
  const [grants, setGrants] = useState(new Set(initialGrants));

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

  const toggleGrant = (name: string) => {
    setGrants((current) => {
      const next = new Set(current);
      if (!next.delete(name)) next.add(name);
      return next;
    });
  };

  return (
    <Card withBorder p="xl" maw={720}>
      <Title order={4} mb="md">
        {t(($) => $.admin.permissions.title)}
      </Title>
      <Stack gap="xs" mb="xl">
        {catalog.map((permission) => (
          <Group key={permission.name} gap="sm">
            <Checkbox
              checked={grants.has(permission.name)}
              label={permission.name}
              onChange={() => toggleGrant(permission.name)}
            />
            <Badge variant="light" color={permission.kind === "read" ? "gray" : "orange"}>
              {permission.kind === "read"
                ? t(($) => $.admin.permissions.kindRead)
                : t(($) => $.admin.permissions.kindWrite)}
            </Badge>
          </Group>
        ))}
      </Stack>
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
          onClick={() =>
            mutation.mutate({ path: { roleId }, body: { permissions: rows, grants: [...grants] } })
          }
        >
          {t(($) => $.admin.crud.save)}
        </Button>
      </Group>
    </Card>
  );
}
