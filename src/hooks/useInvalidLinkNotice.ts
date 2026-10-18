import { notifications } from "@mantine/notifications";
import { useEffect } from "react";

// The refused-link notice is a flash, so it is shown once on arrival like the
// other flashes; the id keeps StrictMode's double effect from stacking two.
export function useInvalidLinkNotice(invalid: boolean | undefined, message: string): void {
  useEffect(() => {
    if (invalid) {
      notifications.show({ color: "red", id: "invalid-link", message });
    }
  }, [invalid, message]);
}
