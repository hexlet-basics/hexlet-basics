import type { AnyRouter } from "@tanstack/react-router";
import { z } from "zod";

// The search of the router's current location as a plain record, to compare with
// what a page should have put in the address.
export function searchOf(router: AnyRouter): Record<string, unknown> {
  return z.record(z.string(), z.unknown()).parse(router.state.location.search);
}
