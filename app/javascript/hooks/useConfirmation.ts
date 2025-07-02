import type React from "react";
import { useTranslation } from "react-i18next";

interface ConfirmationOptions {
  message?: string;
  callback?: UseConfirmationCallback,
}

type UseConfirmationCallback = (event: React.MouseEvent<Element>) => void;

export default function useConfirmation(
  options?: ConfirmationOptions,
): UseConfirmationCallback {
  const { t: tCommon } = useTranslation("common");

  const requestConfirmation = (event: React.MouseEvent<Element>) => {
    const message = options?.message ?? tCommon("confirm");
    const isConfirmed = window.confirm(`${message}`);
    if (!isConfirmed) {
      event.preventDefault();
    } else {
      if (options?.callback) {
        event.preventDefault();
        return options.callback(event);
      }
    }
  };

  return requestConfirmation;
}
