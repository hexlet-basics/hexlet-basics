import { createFormHookContexts } from "@tanstack/react-form";

// The contexts the form's field and form components read, kept apart from the components
// so a module that holds components exports components only.
export const { fieldContext, formContext, useFieldContext, useFormContext } =
  createFormHookContexts();
