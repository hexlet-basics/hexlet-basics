import { Badge, Button, Card, Checkbox, Group, Stack, Text, Title } from "@mantine/core";
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
import type { Permission } from "@/client/types.gen";
import { useResourceMutation } from "@/hooks/useResourceMutation";

// The Permissions a role holds, edited on its page. The list is the catalog the
// contract's operations declare (ADR-0020), grouped by the resource each name
// starts with, so a new Permission appears here with no change to this screen.
// One save replaces the role's whole set, so submitting the screen is idempotent.
export function RolePermissions({ roleId }: { roleId: number }) {
  const { data } = useQuery(adminGetRolePermissionsOptions({ path: { roleId } }));
  const { data: catalog } = useQuery(adminListPermissionsOptions());

  if (!data || !catalog) {
    return null;
  }

  // Remount the form when a fresh set arrives so its state re-seeds.
  return (
    <PermissionsForm
      key={JSON.stringify(data.grants)}
      roleId={roleId}
      catalog={catalog.permissions}
      initialGrants={data.grants}
    />
  );
}

// ResourceOf is the part of a Permission name before the dot.
const resourceOf = (permission: Permission) => permission.name.split(".")[0] ?? permission.name;

function PermissionsForm({
  roleId,
  catalog,
  initialGrants,
}: {
  roleId: number;
  catalog: Permission[];
  initialGrants: string[];
}) {
  const { t } = useTranslation();
  const [grants, setGrants] = useState(initialGrants);

  const mutation = useResourceMutation({
    mutation: adminUpdateRolePermissionsMutation(),
    invalidate: [
      adminGetRolePermissionsQueryKey({ path: { roleId } }),
      adminGetRoleQueryKey({ path: { id: roleId } }),
    ],
    successMessage: t(($) => $.admin.crud.updated),
    errorMessage: t(($) => $.admin.crud.saveError),
  });

  const byResource = new Map<string, Permission[]>();
  for (const permission of catalog) {
    byResource.set(resourceOf(permission), [
      ...(byResource.get(resourceOf(permission)) ?? []),
      permission,
    ]);
  }
  const groups = [...byResource];

  const toggleGrant = (name: string) => {
    setGrants((current) => {
      const next = new Set(current);
      if (!next.delete(name)) {
        next.add(name);
      }
      return next;
    });
  };

  return (
    <Card withBorder p="xl" maw={720}>
      <Title order={4} mb="md">
        {t(($) => $.admin.permissions.title)}
      </Title>
      <Checkbox.Group value={grants} onChange={setGrants}>
        <Stack gap="md">
          {groups.map(([resource, permissions = []]) => (
            <Stack key={resource} gap="xs">
              <Text fw={600}>{resource}</Text>
              {permissions.map((permission) => (
                <Group key={permission.name} gap="sm">
                  <Checkbox value={permission.name} label={permission.name} />
                  {permission.capability && (
                    <Badge variant="outline" color="red">
                      {t(($) => $.admin.permissions.capability)}
                    </Badge>
                  )}
                  <Badge variant="light" color={permission.kind === "read" ? "gray" : "orange"}>
                    {permission.kind === "read"
                      ? t(($) => $.admin.permissions.kindRead)
                      : t(($) => $.admin.permissions.kindWrite)}
                  </Badge>
                </Group>
              ))}
            </Stack>
          ))}
        </Stack>
      </Checkbox.Group>
      <Group justify="flex-end" mt="md">
        <Button
          loading={mutation.isPending}
          onClick={() => mutation.mutate({ path: { roleId }, body: { grants } })}
        >
          {t(($) => $.admin.crud.save)}
        </Button>
      </Group>
    </Card>
  );
}
