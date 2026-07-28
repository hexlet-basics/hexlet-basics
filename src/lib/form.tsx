import { TextInput, type TextInputProps } from "@mantine/core";
import { createFormHook, createFormHookContexts } from "@tanstack/react-form";

// Mantine input wrappers over headless TanStack Form (ADR-0008). `useAppForm`
// binds the field components below so forms read like plain Mantine markup while
// TanStack Form owns state and Zod owns validation. Validators are the Zod
// schemas generated from OpenAPI (`src/client/zod.gen.ts`), passed straight to
// `validators.onSubmit` — zod v4 implements Standard Schema, which TanStack Form
// consumes natively and maps issues back to fields by path.
export const { fieldContext, formContext, useFieldContext } =
  createFormHookContexts();

function fieldError(errors: unknown[], isTouched: boolean): string | undefined {
  if (!isTouched || errors.length === 0) return undefined;
  return errors
    .map((e) =>
      typeof e === "string"
        ? e
        : ((e as { message?: string } | null)?.message ?? ""),
    )
    .filter(Boolean)
    .join(", ");
}

function TextField(props: TextInputProps) {
  const field = useFieldContext<string>();
  return (
    <TextInput
      {...props}
      value={field.state.value}
      onChange={(event) => field.handleChange(event.currentTarget.value)}
      onBlur={field.handleBlur}
      error={fieldError(field.state.meta.errors, field.state.meta.isTouched)}
    />
  );
}

export const { useAppForm } = createFormHook({
  fieldComponents: { TextField },
  formComponents: {},
  fieldContext,
  formContext,
});
