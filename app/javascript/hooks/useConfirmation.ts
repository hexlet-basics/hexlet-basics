import type React from "react";
import { useTranslation } from "react-i18next";

interface ConfirmDialogOptions {
  message: string;
}

type UseConfirmationCallback = (event: React.MouseEvent<Element>) => void;

export default function useConfirmation(
  callback: UseConfirmationCallback,
  options?: ConfirmDialogOptions,
): UseConfirmationCallback {
  const { t: tCommon } = useTranslation("common");

  const requestConfirmation = (event: React.MouseEvent<Element>) => {
    event.preventDefault();

    const message = options?.message ?? tCommon("confirm");
    const isConfirmed = window.confirm(`${message}`);
    if (isConfirmed) {
      callback(event);
    }
  };

  return requestConfirmation;
}
