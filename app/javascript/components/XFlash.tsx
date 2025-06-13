import { hasObjectKey } from "@/lib/utils";
import type { FlashKey, SharedProps } from "@/types";
import type { FlashVariants } from "@/types";
import { usePage } from "@inertiajs/react";
import { Alert } from '@mantine/core';
import XssContent from "./XssContent";
import { useState } from "react";

export default function XFlash() {
  const { flash } = usePage<SharedProps>().props;

  const [dismissed, setDismissed] = useState(false);
  const handleClose = () => setDismissed(true);

  const variants: FlashVariants = {
    error: "red",
    notice: "blue",
    success: "green",
  };

  const flashEntry = (Object.entries(variants) as [FlashKey, string][]).find(
    ([key]) => hasObjectKey(flash, key) && !dismissed
  );

  if (!flashEntry) return null;

  const [key, color] = flashEntry;
  const value = flash[key]

  return (
    <Alert
      mb="md"
      color={color}
      withCloseButton
      onClose={handleClose}
    >
      <XssContent>{value}</XssContent>
    </Alert>
  );
}
