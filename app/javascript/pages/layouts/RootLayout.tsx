import type { SharedProps } from "@/types";
import { usePage } from "@inertiajs/react";
import { usePostHog } from "posthog-js/react";
import { useEffect, type PropsWithChildren } from "react";

type Props = PropsWithChildren & {};
export default (props: Props) => {
  const { event } = usePage<SharedProps>().props;
  const posthog = usePostHog();

  // biome-ignore lint/correctness/useExhaustiveDependencies: <explanation>
  useEffect(() => {
    if (event) {
      switch (event.event_type) {
        case "UserSignedInEvent":
          // add carrot quest
          // posthog?.identify(event.data.id, {
          //   email: event.data.email,
          // });
          break;
        case "UserSignedUpEvent":
          break;
        default:
          break;
      }
    }
  }, [event]);

  return props.children;
};
