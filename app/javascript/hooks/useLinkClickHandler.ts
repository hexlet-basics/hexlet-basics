import { router } from "@inertiajs/react";
import type { RequestPayload } from "@inertiajs/core";
import { useCallback } from "react";

const useLinkClickHandler = () => {
  return useCallback(
    (method: "get" | "post" | "patch" | "delete", data?: RequestPayload) => {
      return (event: React.MouseEvent<HTMLAnchorElement>) => {
        event.preventDefault(); // Prevent default link behavior

        const url = (event.currentTarget as HTMLAnchorElement).href; // Extract the URL from the clicked link

        router.visit(url, {
          method,
          data,
        });
      };
    },
    [],
  );
};

export default useLinkClickHandler;
