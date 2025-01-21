import type { SharedProps } from "@/types";
import { usePage } from "@inertiajs/react";
import { Alert } from "react-bootstrap";

export default function XFlash() {
  const { flash } = usePage<SharedProps>().props;
  const variants = {
    error: "danger",
    notice: "primary",
    success: "success",
  };

  return (
    <>
      {Object.entries(variants).map(
        ([key, variant]) =>
          flash[key] && (
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
