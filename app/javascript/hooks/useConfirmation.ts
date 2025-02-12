import type React from "react";
import { useTranslation } from "react-i18next";

interface ConfirmDialogOptions {
  message: string;
}

type UseConfirmationResult = (event: React.SyntheticEvent) => void;

export default function useConfirmation(
  callback: (cb: React.SyntheticEvent) => void,
  options?: ConfirmDialogOptions,
): UseConfirmationResult {
  const { t: tCommon } = useTranslation("common");

  const requestConfirmation = (event: React.SyntheticEvent) => {
    event.preventDefault();

    const message = options?.message ?? tCommon("confirm");
    const isConfirmed = window.confirm(`${message}`);
    if (isConfirmed) {
      callback(event);
    }
  };

  return requestConfirmation;
}
