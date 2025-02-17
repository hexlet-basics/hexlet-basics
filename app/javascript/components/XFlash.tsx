import { hasObjectKey } from "@/lib/utils";
import type { SharedProps } from "@/types";
import type { FlashVariants } from "@/types";
import { usePage } from "@inertiajs/react";
import { Alert } from "react-bootstrap";

export default function XFlash() {
  const { flash } = usePage<SharedProps>().props;

  const variants: FlashVariants = {
    error: "danger",
    notice: "primary",
    success: "success",
  };

  return (
    <>
      {Object.entries(variants).map(
        ([key, variant]) =>
          hasObjectKey(flash, key) && (
            <Alert
              dismissible
              key={key}
              className="my-3 border-0"
              variant={variant}
            >
              {flash[key]}
            </Alert>
          ),
      )}
    </>
  );
}
