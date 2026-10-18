import type { z } from "zod";

// A search parameter that falls back to "absent" when it is malformed, instead of
// failing the route. Zod's `.catch` supplies a fallback value; it is not a
// promise's `.catch`, which the promise rules cannot tell from it.
export function lenient<TSchema extends z.ZodType>(
  schema: TSchema,
): z.ZodCatch<z.ZodOptional<TSchema>> {
  // oxlint-disable-next-line promise/prefer-await-to-then, unicorn/prefer-top-level-await -- zod's catch, not a promise
  return schema.optional().catch(undefined);
}

// What a promise resolves to, or null when it rejects: for a read that is
// decoration, whose failure must not take the page down.
export async function orNull<TValue>(promise: Promise<TValue>): Promise<TValue | null> {
  try {
    return await promise;
  } catch {
    return null;
  }
}
