import {
  Button,
  type ButtonProps,
  Checkbox,
  Input,
  MultiSelect,
  NumberInput,
  Select,
  TextInput,
  Textarea,
} from "@mantine/core";
import type {
  CheckboxProps,
  MultiSelectProps,
  NumberInputProps,
  SelectProps,
  TextInputProps,
  TextareaProps,
} from "@mantine/core";
import { DateTimePicker } from "@mantine/dates";
import type { DateTimePickerProps } from "@mantine/dates";
import { Link, RichTextEditor } from "@mantine/tiptap";
import { useEditor } from "@tiptap/react";
import StarterKit from "@tiptap/starter-kit";
import { createFormHook, createFormHookContexts } from "@tanstack/react-form";
import dayjs from "dayjs";
import type { SyntheticEvent } from "react";
import localizedFormat from "dayjs/plugin/localizedFormat";
import type { ReactElement } from "react";

// `L LT` (localized date + time) needs the localizedFormat plugin.
dayjs.extend(localizedFormat);

// Mantine's DateTimePicker plumbs its value as a fixed `YYYY-MM-DD HH:mm:ss`
// local string, but the whole app (and every generated `z.iso.datetime()`
// validator) speaks RFC3339. Converting here, once, lets the field value BE the
// API's ISO string — so resources validate against the generated `zXxxInput`
// directly instead of a hand-written form schema.
const MANTINE_DATETIME = "YYYY-MM-DD HH:mm:ss";
const isoToPicker = (iso: string | null): string | null =>
  iso ? dayjs(iso).format(MANTINE_DATETIME) : null;
const pickerToIso = (value: string | null): string | null =>
  value ? dayjs(value).toISOString() : null;

// Mantine input wrappers over headless TanStack Form (ADR-0008). `useAppForm`
// binds the field components below so forms read like plain Mantine markup while
// TanStack Form owns state and Zod owns validation. Validators are the Zod
// schemas generated from OpenAPI (`src/client/zod.gen.ts`), passed straight to
// `validators.onSubmit` — zod v4 implements Standard Schema, which TanStack Form
// consumes natively and maps issues back to fields by path.
export const { fieldContext, formContext, useFieldContext } = createFormHookContexts();

function fieldError(errors: unknown[], isTouched: boolean): string | undefined {
  if (!isTouched || errors.length === 0) {
    return undefined;
  }
  return errors
    .map((entry) =>
      typeof entry === "string" ? entry : ((entry as { message?: string } | null)?.message ?? ""),
    )
    .filter(Boolean)
    .join(", ");
}

function TextField(props: TextInputProps): ReactElement {
  const field = useFieldContext<string>();
  return (
    <TextInput
      {...props}
      // Nullable string columns arrive as null; render "" so the input stays
      // controlled (the value type is the generated Input, which allows null).
      value={field.state.value ?? ""}
      onChange={(event) => field.handleChange(event.currentTarget.value)}
      onBlur={field.handleBlur}
      error={fieldError(field.state.meta.errors, field.state.meta.isTouched)}
    />
  );
}

function TextareaField(props: TextareaProps): ReactElement {
  const field = useFieldContext<string>();
  return (
    <Textarea
      autosize
      minRows={3}
      {...props}
      value={field.state.value ?? ""}
      onChange={(event) => field.handleChange(event.currentTarget.value)}
      onBlur={field.handleBlur}
      error={fieldError(field.state.meta.errors, field.state.meta.isTouched)}
    />
  );
}

// Enum/choice fields. The field value stays a plain string; callers pass Mantine's
// `data` (option list). `null` (cleared) maps to "" so the value stays controlled.
function SelectField(props: SelectProps): ReactElement {
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

// Date/time fields. The field value is the API's RFC3339 string (or null); the
// Mantine picker-format conversion is hidden here, and `valueFormat="L LT"` shows
// a locale-aware date instead of a hard-coded display pattern.
function DateTimeField(props: DateTimePickerProps): ReactElement {
  const field = useFieldContext<string | null>();
  return (
    <DateTimePicker
      clearable
      valueFormat="L LT"
      {...props}
      value={isoToPicker(field.state.value)}
      onChange={(value) => field.handleChange(pickerToIso(value))}
      onBlur={field.handleBlur}
      error={fieldError(field.state.meta.errors, field.state.meta.isTouched)}
    />
  );
}

// Numeric fields (ids, counters). The generated inputs type these as
// `number | null`, and Mantine reports a cleared input as ""  — map that back
// to null so a blank field submits the contract's null, not zero.
function NumberField(props: NumberInputProps): ReactElement {
  const field = useFieldContext<number | null>();
  return (
    <NumberInput
      {...props}
      value={field.state.value ?? ""}
      onChange={(value) => field.handleChange(typeof value === "number" ? value : null)}
      onBlur={field.handleBlur}
      error={fieldError(field.state.meta.errors, field.state.meta.isTouched)}
    />
  );
}

// Boolean flags. Nullable contract booleans render as unchecked; toggling
// always submits a real boolean (legacy forms never sent null for checkboxes).
function CheckboxField(props: CheckboxProps): ReactElement {
  const field = useFieldContext<boolean | null>();
  return (
    <Checkbox
      {...props}
      checked={field.state.value ?? false}
      onChange={(event) => field.handleChange(event.currentTarget.checked)}
      onBlur={field.handleBlur}
      error={fieldError(field.state.meta.errors, field.state.meta.isTouched)}
    />
  );
}

// String-array fields (e.g. a staff member's allowed locales). The value IS the
// contract's string[]; callers pass Mantine's `data` option list.
function MultiSelectField(props: MultiSelectProps): ReactElement {
  const field = useFieldContext<string[]>();
  return (
    <MultiSelect
      {...props}
      value={field.state.value ?? []}
      onChange={(value) => field.handleChange(value)}
      onBlur={field.handleBlur}
      error={fieldError(field.state.meta.errors, field.state.meta.isTouched)}
    />
  );
}

// Trusted-HTML fields (the blog post body). The field value IS the stored
// HTML: tiptap parses it on mount and hands back `getHTML()` on every change,
// so the form round-trips exactly what the API stores. `immediatelyRender:
// false` keeps the editor SSR-safe under TanStack Start.
function RichTextField({ label }: { label: string }): ReactElement {
  const field = useFieldContext<string>();
  const editor = useEditor({
    content: field.state.value ?? "",
    extensions: [StarterKit, Link],
    immediatelyRender: false,
    onBlur: field.handleBlur,
    onUpdate: ({ editor: current }) => field.handleChange(current.getHTML()),
  });

  return (
    <Input.Wrapper
      label={label}
      error={fieldError(field.state.meta.errors, field.state.meta.isTouched)}
    >
      <RichTextEditor editor={editor}>
        <RichTextEditor.Toolbar sticky>
          <RichTextEditor.ControlsGroup>
            <RichTextEditor.Bold />
            <RichTextEditor.Italic />
            <RichTextEditor.Code />
            <RichTextEditor.ClearFormatting />
          </RichTextEditor.ControlsGroup>
          <RichTextEditor.ControlsGroup>
            <RichTextEditor.H2 />
            <RichTextEditor.H3 />
            <RichTextEditor.Blockquote />
            <RichTextEditor.CodeBlock />
          </RichTextEditor.ControlsGroup>
          <RichTextEditor.ControlsGroup>
            <RichTextEditor.BulletList />
            <RichTextEditor.OrderedList />
          </RichTextEditor.ControlsGroup>
          <RichTextEditor.ControlsGroup>
            <RichTextEditor.Link />
            <RichTextEditor.Unlink />
          </RichTextEditor.ControlsGroup>
        </RichTextEditor.Toolbar>
        <RichTextEditor.Content />
      </RichTextEditor>
    </Input.Wrapper>
  );
}

// The `onSubmit` of a native `<form>` driven by a TanStack form: stops the browser's
// own submission and lets the form validate and submit.
export const submitForm =
  (form: { handleSubmit: () => Promise<void> }) => (event: SyntheticEvent) => {
    event.preventDefault();
    void form.handleSubmit();
  };

// The submit button of a form: busy for as long as the submission runs. The forms
// await their mutation in `onSubmit`, so the button needs no pending flag handed
// down from it. Rendered inside `<form.AppForm>`, which supplies the context.
function SubmitButton({ children, ...props }: ButtonProps) {
  const form = useFormContext();
  return (
    <form.Subscribe selector={(state) => state.isSubmitting}>
      {(isSubmitting) => (
        <Button type="submit" loading={isSubmitting} {...props}>
          {children}
        </Button>
      )}
    </form.Subscribe>
  );
}

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
  formComponents: {},
  formContext,
});
