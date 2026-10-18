import { useCookies } from "react-cookie";
import { z } from "zod";
import { SECONDS_PER_YEAR } from "@/lib/duration";

// Where the learner left the divider between the theory and the workspace, as
// the theory's share of the width in percent. It lives in a cookie rather than
// in localStorage because the lesson page is server-rendered (ADR-0008): the
// server renders the split from the request's cookie, so the first paint is
// already the learner's and nothing moves once the page hydrates. One value for
// every lesson — a cookie per lesson would pile up and ride along on every
// request.
const COOKIE_NAME = "lesson_panes";
const COOKIE_TTL_SECONDS = SECONDS_PER_YEAR;

// The split before the learner drags the divider, and the bounds either pane
// may take (the theory pane's own `min`, mirrored for the workspace).
const DEFAULT_THEORY_PERCENT = 40;
// The divider position is stored to one decimal place.
const DECIMAL_STEP = 10;
const MIN_THEORY_PERCENT = 25;
const MAX_THEORY_PERCENT = 75;
const THEORY_PERCENT = z.coerce.number().min(MIN_THEORY_PERCENT).max(MAX_THEORY_PERCENT);

// The stored theory width in percent — the default when there is none or it is
// not a usable number — and a saver to call when a drag ends.
export function useTheoryPercent(): [number, (percent: number) => void] {
  const [cookies, setCookie] = useCookies([COOKIE_NAME]);
  const parsed = THEORY_PERCENT.safeParse(cookies[COOKIE_NAME]);

  const save = (percent: number): void => {
    // `path` is explicit: the library defaults to the current page, which would
    // keep one split per lesson URL after all.
    setCookie(COOKIE_NAME, String(Math.round(percent * DECIMAL_STEP) / DECIMAL_STEP), {
      maxAge: COOKIE_TTL_SECONDS,
      path: "/",
      sameSite: "lax",
    });
  };

  return [parsed.success ? parsed.data : DEFAULT_THEORY_PERCENT, save];
}
