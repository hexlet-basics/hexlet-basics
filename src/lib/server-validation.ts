import { zValidationError } from "@/client/zod.gen";
import { isApiError } from "@/lib/api-error";

// The API answers a write that fails validation with 422 and a generated
// `ValidationError` body: messages keyed by field, covering rules the Zod
// schema cannot express (an email already taken). Parsing it with the generated
// `zValidationError` keeps this in step with the contract.
function validationErrors(error: unknown): Record<string, string[]> | undefined {
  if (!isApiError(error, 422)) return undefined;
  const parsed = zValidationError.safeParse(error.body);
  return parsed.success ? parsed.data.errors : undefined;
}

// For a mutation's onError: a 422 is shown at the fields by
// `submitWithServerErrors`, so only the other failures need a toast.
export const isValidationError = (error: unknown): boolean => validationErrors(error) !== undefined;

// A TanStack Form `validators.onSubmitAsync` that runs the submission and hands
// a 422 back in the form's own `{ fields }` shape, which is the library's way
// to put server errors on fields: they render through the same field `error`
// as the Zod ones and clear themselves on the next validation, with no state
// to reset by hand. Any other failure is left to the mutation's onError.
export function submitWithServerErrors<TValue>(submit: (value: TValue) => Promise<unknown>) {
  return async ({ value }: { value: TValue }) => {
    try {
      await submit(value);
    } catch (error) {
      const errors = validationErrors(error);
      if (!errors) return undefined;
      return {
        fields: Object.fromEntries(
          Object.entries(errors).map(([field, messages]) => [field, messages.join(", ")]),
        ),
      };
    }
    return undefined;
  };
}

// Every message of a 422, flattened, for a screen that has no form fields to put
// them on (a filter bar): empty when the failure is not a validation error.
export const validationMessages = (error: unknown): string[] =>
  Object.values(validationErrors(error) ?? {}).flat();
