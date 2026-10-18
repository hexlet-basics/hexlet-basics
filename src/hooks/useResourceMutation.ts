import { notifications } from "@mantine/notifications";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { QueryKey, UseMutationOptions } from "@tanstack/react-query";
import { isValidationError } from "@/lib/server-validation";

// The three things every admin CRUD mutation repeats — invalidate the affected
// list/detail queries, toast on success, toast on failure — collapse into
// declarative fields, so route files stop hand-wiring onSuccess/onError. Adapted
// from the sibling 1mail stack. Typing is inferred from the passed generated
// mutation, so callers keep full type-safety on variables.
interface ResourceMutationOptions<TData, TError, TVars, TCtx> {
  // The generated react-query mutation, e.g. adminCreateBannerMutation().
  mutation: UseMutationOptions<TData, TError, TVars, TCtx>;
  // Query keys invalidated after success (list query, edited detail query).
  invalidate?: QueryKey[];
  // For a mutation a form submits: a 422 is then shown at the form's fields
  // (submitWithServerErrors), so it gets no toast of its own. Left off for
  // mutations with no form, where the toast is the only feedback.
  fieldErrors?: boolean;
  successMessage: string;
  errorMessage: string;
  // Extra success work — navigation, mostly — run after invalidation.
  onDone?: (data: TData, variables: TVars) => void;
}

export function useResourceMutation<TData, TError, TVars, TCtx>({
  mutation,
  invalidate,
  fieldErrors,
  successMessage,
  errorMessage,
  onDone,
}: ResourceMutationOptions<TData, TError, TVars, TCtx>) {
  const queryClient = useQueryClient();

  return useMutation({
    ...mutation,
    onSuccess: async (data, variables) => {
      await Promise.all(
        (invalidate ?? []).map((queryKey) => queryClient.invalidateQueries({ queryKey })),
      );
      notifications.show({ color: "green", message: successMessage });
      onDone?.(data, variables);
    },
    onError: (error) => {
      if (fieldErrors && isValidationError(error)) {
        return;
      }
      notifications.show({ color: "red", message: errorMessage });
    },
  });
}
