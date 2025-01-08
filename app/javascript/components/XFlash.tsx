import type { SharedProps } from "@/types/types";
import { usePage } from "@inertiajs/react";
import { Alert } from "react-bootstrap";

export default function XFlash() {
  const { flash } = usePage<SharedProps>().props;
  const variants = [
    "primary",
    "secondary",
    "success",
    "danger",
    "warning",
    "info",
    "light",
    "dark",
  ];

  return (
    <>
      {variants.map(
        (variant) =>
          flash[variant] && (
            <Alert dismissible key={variant} className="my-3 border-0" variant={variant}>
              {flash[variant]}
            </Alert>
          ),
      )}
    </>
  );
}
