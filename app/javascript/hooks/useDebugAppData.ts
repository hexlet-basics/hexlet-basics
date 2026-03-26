import { router, usePage } from "@inertiajs/react";
import { useEffect } from "react";

export function useDebugAppData() {
  const page = usePage();

  useEffect(() => {
    console.group("[Inertia][initial]");
    console.log("url:", page.url);
    console.log("component:", page.component);
    console.log("props:", page.props);
    console.log("flash:", page.flash);
    console.groupEnd();

    const unsubscribe = router.on("navigate", (event) => {
      console.group("[Inertia][navigate]");
      console.log("url:", event.detail.page.url);
      console.log("component:", event.detail.page.component);
      console.log("props:", event.detail.page.props);
      console.log("flash:", event.detail.page.flash);
      console.groupEnd();
    });

    return unsubscribe;
  }, [page]);
}
