import { router, usePage } from "@inertiajs/react";
import { useEffect } from "react";

// падает на SSR
// import { gon } from '../gon';

function readPageScript() {
  const el = document.querySelector<HTMLScriptElement>(
    'script[data-page="app"]',
  );
  if (!el) {
    console.warn('Script with data-page="app" not found');
    return null;
  }

  try {
    return JSON.parse(el.textContent || "");
  } catch (e) {
    console.error("Failed to parse data-page JSON", e);
    return null;
  }
}

export function useDebugAppData() {
  const { flash } = usePage();
  useEffect(() => {
    // Лог при первом рендере
    const initial = readPageScript();
    // if (initial) {
    //   console.log("[Inertia][initial] props:", initial.props ?? initial);
    // }

    // Подписка на навигацию Inertia
    const unsubscribe = router.on("navigate", (event) => {
      const data = readPageScript();
      console.group("[Inertia][navigate]");
      console.log("url:", event.detail.page.url);
      console.log("component:", event.detail.page.component);
      console.log("props from script:", data?.props ?? data);
      console.log("flash", flash);
      // console.log('gon:', gon);
      console.groupEnd();
    });

    // Отписка при размонтировании
    return unsubscribe;
  }, [flash]);
}
