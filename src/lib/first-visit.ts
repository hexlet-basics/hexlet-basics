import Cookies from "js-cookie";
import type { FirstVisit } from "@/client/types.gen";
import { zFirstVisit } from "@/client/zod.gen";

// First-visit attribution (ADR-0015). Legacy asked ahoy's visit for the UTM
// parameters, the landing page and the referrer when a lead came in; ahoy is
// not ported, so the browser writes them to a cookie on the visitor's first
// page and the lead form sends them along. Written once and never refreshed:
// what amoCRM should credit is where the visitor first came from.
const COOKIE_NAME = "first_visit";
// Ahoy's visitor cookie lived two years; the first visit is kept as long.
const COOKIE_TTL_DAYS = 730;

// Records the current page as the first visit unless one is already recorded.
// Browser-only: it reads the address bar and document.referrer.
export function recordFirstVisit(): void {
  if (Cookies.get(COOKIE_NAME) !== undefined) return;

  const params = new URLSearchParams(window.location.search);
  const visit: FirstVisit = {
    utmSource: params.get("utm_source"),
    utmMedium: params.get("utm_medium"),
    utmCampaign: params.get("utm_campaign"),
    utmContent: params.get("utm_content"),
    utmTerm: params.get("utm_term"),
    landingPage: window.location.href,
    referrer: document.referrer || null,
  };
  Cookies.set(COOKIE_NAME, JSON.stringify(visit), {
    expires: COOKIE_TTL_DAYS,
    path: "/",
    sameSite: "lax",
    secure: window.location.protocol === "https:",
  });
}

// The recorded first visit, or null when there is none or it does not parse
// (a cookie edited by hand or written by an older shape).
export function readFirstVisit(): FirstVisit | null {
  const raw = Cookies.get(COOKIE_NAME);
  if (raw === undefined) return null;
  try {
    const parsed = zFirstVisit.safeParse(JSON.parse(raw));
    return parsed.success ? parsed.data : null;
  } catch {
    return null;
  }
}
