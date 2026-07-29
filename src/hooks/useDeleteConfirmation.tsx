import { Text } from "@mantine/core";
import { modals } from "@mantine/modals";
import type { ReactNode } from "react";
import { useTranslation } from "react-i18next";

// Opens the standard admin delete-confirm modal (shared labels + red confirm),
// so route lists only supply what to delete and the callback. Adapted from 1mail.
export function useDeleteConfirmation() {
  const { t } = useTranslation();

  return ({ description, onConfirm }: { description?: ReactNode; onConfirm: () => void }) => {
    modals.openConfirmModal({
      title: t(($) => $.admin.crud.confirmDeleteTitle),
      children: <Text>{description}</Text>,
      labels: {
        confirm: t(($) => $.admin.crud.delete),
        cancel: t(($) => $.admin.crud.cancel),
      },
      confirmProps: { color: "red" },
      onConfirm,
    });
  };
}
