import { ActionIcon, Group, Tooltip } from "@mantine/core";
import type { QueryKey, UseMutationOptions } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { useDeleteConfirmation } from "@/hooks/useDeleteConfirmation";
import { useResourceMutation } from "@/hooks/useResourceMutation";
import type { ReactNode } from "react";

interface RowActionsProps<TData, TError> {
  // The row's own buttons (edit link, extra operations), shown before delete.
  children: ReactNode;
  // What the delete confirmation names, e.g. the row's title.
  describe: ReactNode;
  id: number;
  // The generated list query key a delete invalidates, and the generated delete
  // mutation. Both left out for a row that cannot be deleted.
  listKey?: QueryKey;
  remove?: UseMutationOptions<TData, TError, { path: { id: number } }>;
}

// The last cell of an admin list row: its buttons and the confirmed delete.
export function RowActions<TData, TError>({
  children,
  describe,
  id,
  listKey,
  remove,
}: RowActionsProps<TData, TError>): ReactNode {
  const { t } = useTranslation();
  const confirmDelete = useDeleteConfirmation();
  const deleteMutation = useResourceMutation({
    errorMessage: t(($) => $.admin.crud.deleteError),
    invalidate: listKey ? [listKey] : [],
    mutation: remove ?? {},
    successMessage: t(($) => $.admin.crud.deleted),
  });

  return (
    <Group gap="xs" justify="flex-end" wrap="nowrap">
      {children}
      {remove && (
        <Tooltip label={t(($) => $.admin.crud.delete)}>
          <ActionIcon
            color="red"
            variant="light"
            aria-label={t(($) => $.admin.crud.delete)}
            onClick={() => {
              confirmDelete({
                description: describe,
                onConfirm: () => {
                  deleteMutation.mutate({ path: { id } });
                },
              });
            }}
          >
            ✕
          </ActionIcon>
        </Tooltip>
      )}
    </Group>
  );
}
