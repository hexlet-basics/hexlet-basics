import { createFormHook } from "@tanstack/react-form";
import type { SyntheticEvent } from "react";
import {
  CheckboxField,
  DateTimeField,
  MultiSelectField,
  NumberField,
  RichTextField,
  SelectField,
  SubmitButton,
  TextField,
  TextareaField,
} from "@/lib/form-fields";
import { fieldContext, formContext } from "@/lib/form-context";

// The `onSubmit` of a native `<form>` driven by a TanStack form: stops the browser's
// own submission and lets the form validate and submit.
export const submitForm =
  (form: { handleSubmit: () => Promise<void> }) =>
  (event: SyntheticEvent): void => {
    event.preventDefault();
    void form.handleSubmit();
  };

// Mantine input wrappers over headless TanStack Form (ADR-0008). `useAppForm`
// binds the field components below so forms read like plain Mantine markup while
// TanStack Form owns state and Zod owns validation. Validators are the Zod
// schemas generated from OpenAPI (`src/client/zod.gen.ts`), passed straight to
// `validators.onSubmit` — zod v4 implements Standard Schema, which TanStack Form
// consumes natively and maps issues back to fields by path.
export const { useAppForm } = createFormHook({
  fieldComponents: {
    CheckboxField,
    DateTimeField,
    MultiSelectField,
    NumberField,
    RichTextField,
    SelectField,
    TextField,
    TextareaField,
  },
  fieldContext,
  formComponents: { SubmitButton },
  formContext,
});
