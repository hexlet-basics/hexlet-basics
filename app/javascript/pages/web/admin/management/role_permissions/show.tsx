import { useForm } from "@inertiajs/react";
import { Button, Checkbox, Group, Table } from "@mantine/core";
import { useCallback } from "react";
import { useTranslation } from "react-i18next";
import AdminLayout from "@/layouts/AdminLayout";
import * as Routes from "@/routes.js";
import type { PermissionResource, StaffRoleCrud } from "@/types";
import { Menu } from "../roles/shared/menu";

type PermissionFormRow = {
  resource: PermissionResource;
  can_index: boolean;
  can_create: boolean;
  can_update: boolean;
  can_destroy: boolean;
};

type Props = {
  role: StaffRoleCrud;
  resources: PermissionResource[];
};

type ActionKey = keyof Omit<PermissionFormRow, "resource">;

function resourceLabel(resource: PermissionResource): string {
  return resource
    .split("_")
    .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
    .join(" ");
}

export default function Show({ role, resources }: Props) {
  const { t } = useTranslation();

  const initialPermissions: PermissionFormRow[] = resources.map((resource) => {
    const existing = role.permissions.find((p) => p.resource === resource);
    return {
      resource,
      can_index: existing?.can_index ?? false,
      can_create: existing?.can_create ?? false,
      can_update: existing?.can_update ?? false,
      can_destroy: existing?.can_destroy ?? false,
    };
  });

  const { data, setData, transform, patch, processing } = useForm({
    permissions: initialPermissions,
  });

  const handleCheck = useCallback(
    (index: number, field: ActionKey, checked: boolean) => {
      setData(
        "permissions",
        data.permissions.map((row, i) =>
          i === index ? { ...row, [field]: checked } : row,
        ),
      );
    },
    [data.permissions, setData],
  );

  const handleCheckAll = useCallback(
    (index: number, checked: boolean) => {
      setData(
        "permissions",
        data.permissions.map((row, i) =>
          i === index
            ? {
                ...row,
                can_index: checked,
                can_create: checked,
                can_update: checked,
                can_destroy: checked,
              }
            : row,
        ),
      );
    },
    [data.permissions, setData],
  );

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    transform((d) => ({ data: d }));
    patch(Routes.admin_management_role_permission_path(role.id));
  };

  const columns: { key: ActionKey; label: string }[] = [
    { key: "can_index", label: t(($) => $.helpers.crud.list) },
    { key: "can_create", label: t(($) => $.helpers.crud.adding) },
    { key: "can_update", label: t(($) => $.helpers.crud.editing) },
    { key: "can_destroy", label: t(($) => $.helpers.crud.remove) },
  ];

  return (
    <AdminLayout
      header={t(($) => $.admin.management.role_permissions.show.header, {
        name: role.name,
      })}
    >
      <Menu data={role} />
      <form onSubmit={handleSubmit}>
        <Table striped highlightOnHover>
          <Table.Thead>
            <Table.Tr>
              <Table.Th>
                {t(($) => $.admin.management.role_permissions.show.resource)}
              </Table.Th>
              <Table.Th>
                {t(($) => $.admin.management.role_permissions.show.all)}
              </Table.Th>
              {columns.map((col) => (
                <Table.Th key={col.key}>{col.label}</Table.Th>
              ))}
            </Table.Tr>
          </Table.Thead>
          <Table.Tbody>
            {data.permissions.map((row, index) => {
              const checkedCount = columns.filter((col) => row[col.key]).length;
              const allChecked = checkedCount === columns.length;
              const someChecked = checkedCount > 0 && !allChecked;
              return (
                <Table.Tr key={row.resource}>
                  <Table.Td>{resourceLabel(row.resource)}</Table.Td>
                  <Table.Td>
                    <Checkbox
                      checked={allChecked}
                      indeterminate={someChecked}
                      onChange={(e) =>
                        handleCheckAll(index, e.currentTarget.checked)
                      }
                    />
                  </Table.Td>
                  {columns.map((col) => (
                    <Table.Td key={col.key}>
                      <Checkbox
                        checked={row[col.key]}
                        onChange={(e) =>
                          handleCheck(index, col.key, e.currentTarget.checked)
                        }
                      />
                    </Table.Td>
                  ))}
                </Table.Tr>
              );
            })}
          </Table.Tbody>
        </Table>
        <Group justify="flex-end" mt="md">
          <Button loading={processing} type="submit">
            {t(($) => $.helpers.submit.save)}
          </Button>
        </Group>
      </form>
    </AdminLayout>
  );
}
