import {
  Select,
  type SelectProps,
  Textarea,
  type TextareaProps,
  TextInput,
  type TextInputProps,
} from "@mantine/core";
import { DateTimePicker, type DateTimePickerProps } from "@mantine/dates";
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

function TextareaField(props: TextareaProps) {
  const field = useFieldContext<string>();
  return (
    <Textarea
      autosize
      minRows={3}
      {...props}
      value={field.state.value}
      onChange={(event) => field.handleChange(event.currentTarget.value)}
      onBlur={field.handleBlur}
      error={fieldError(field.state.meta.errors, field.state.meta.isTouched)}
    />
  );
}

// Enum/choice fields. The field value stays a plain string; callers pass Mantine's
// `data` (option list). `null` (cleared) maps to "" so the value stays controlled.
function SelectField(props: SelectProps) {
  const field = useFieldContext<string>();
  return (
    <Select
      {...props}
      value={field.state.value}
      onChange={(value) => field.handleChange(value ?? "")}
      onBlur={field.handleBlur}
      error={fieldError(field.state.meta.errors, field.state.meta.isTouched)}
    />
  );
}

// Date/time fields. The field value is Mantine's `YYYY-MM-DD HH:mm:ss` string (or
// "" when empty); converting to/from the API's RFC3339 is the resource module's
// job, so the engine stays format-agnostic.
function DateTimeField(props: DateTimePickerProps) {
  const field = useFieldContext<string>();
  return (
    <DateTimePicker
      clearable
      {...props}
      value={field.state.value || null}
      onChange={(value) => field.handleChange(value ?? "")}
      onBlur={field.handleBlur}
      error={fieldError(field.state.meta.errors, field.state.meta.isTouched)}
    />
  );
}

export const { useAppForm } = createFormHook({
  fieldComponents: { TextField, TextareaField, SelectField, DateTimeField },
  formComponents: {},
  fieldContext,
  formContext,
});
