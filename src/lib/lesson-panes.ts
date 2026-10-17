import { createIsomorphicFn } from "@tanstack/react-start";
import { getCookie } from "@tanstack/react-start/server";
import Cookies from "js-cookie";
import { z } from "zod";

// Where the learner left the divider between the theory and the workspace, as
// the theory's share of the width in percent. It lives in a cookie rather than
// in localStorage because the lesson page is server-rendered (ADR-0008): the
// server reads the cookie and renders that split, so the first paint is already
// the learner's and nothing moves once the page hydrates. One value for every
// lesson — a cookie per lesson would pile up and ride along on every request.
const COOKIE_NAME = "lesson_panes";
const COOKIE_TTL_DAYS = 365;

// The split before the learner drags the divider, and the least either pane
// may take (the theory pane's `min`, mirrored for the workspace).
export const DEFAULT_THEORY_PERCENT = 40;
const THEORY_PERCENT = z.coerce.number().min(25).max(75);

const readCookie = createIsomorphicFn()
  .server(() => getCookie(COOKIE_NAME))
  .client(() => Cookies.get(COOKIE_NAME));

// The theory pane's width in percent: the stored one, or the default when there
// is none or it is not a usable number.
export function readTheoryPercent(): number {
  const parsed = THEORY_PERCENT.safeParse(readCookie());
  return parsed.success ? parsed.data : DEFAULT_THEORY_PERCENT;
}

// Browser-only: called when a drag ends.
export function writeTheoryPercent(percent: number): void {
  Cookies.set(COOKIE_NAME, String(Math.round(percent * 10) / 10), {
    expires: COOKIE_TTL_DAYS,
    sameSite: "lax",
  });
}
