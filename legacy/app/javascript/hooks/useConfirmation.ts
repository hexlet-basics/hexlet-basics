// @ts-nocheck
import { openConfirmModal } from "@mantine/modals";
import type React from "react";
import { useTranslation } from "react-i18next";

type UseConfirmationCallback<T = unknown> = (
  event: React.MouseEvent<Element>,
  data?: T,
) => void;

interface ConfirmationOptions<T = unknown> {
  message?: string;
  callback?: UseConfirmationCallback<T>;
}

export default function useConfirmation<T = unknown>(
  options?: ConfirmationOptions<T>,
): UseConfirmationCallback<T> {
  const { t } = useTranslation();

  const requestConfirmation: UseConfirmationCallback<T> = (event, data) => {
    event.preventDefault();
    const title = options?.message ?? t(($) => $.common.confirm);

    openConfirmModal({
      centered: true,
      title: title,
      labels: {
        confirm: t(($) => $.common.boolean.yes),
        cancel: t(($) => $.common.boolean.no),
      },
      onConfirm: () => {
        if (options?.callback) {
          options.callback(event, data);
        }
      },
    });
  };

  return requestConfirmation;
}
