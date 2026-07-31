import { Button, Card, Group, Stack } from "@mantine/core";
import { useTranslation } from "react-i18next";
import type { ZodType } from "zod";
import { useAppForm } from "@/lib/form";

// One editable field, described declaratively. `name` is checked against the
// resource's input model, so a typo or a renamed contract field fails at compile
// time — the whole point of driving the form off the generated types.
export interface CrudFieldSpec<T> {
  name: keyof T & string;
  label: string;
  type?: "text" | "textarea" | "select" | "datetime" | "number" | "checkbox" | "multiselect";
  // Options for `type: "select"` / `"multiselect"`. Ignored otherwise.
  options?: { value: string; label: string }[];
  required?: boolean;
}

export interface CrudFormProps<T extends Record<string, unknown>> {
  fields: CrudFieldSpec<T>[];
  // The generated `zXxxInput` schema, typed to the form values on BOTH the output
  // and input side (`ZodType<T, T>`). zod v4 implements Standard Schema keyed to
  // its input type, and TanStack Form's validator wants `StandardSchemaV1<T>`;
  // fixing input=output=T lets `zXxxInput` satisfy it directly, with no cast.
  schema: ZodType<T, T>;
  defaultValues: T;
  // Submits the validated values. Fire-and-forget: the resource mutation
  // (useResourceMutation) toasts success/failure, so the form neither awaits nor
  // tracks a server error itself.
  onSubmit: (values: T) => void;
  submitLabel: string;
  // Invoked when the user cancels. The route owns navigation so the engine stays
  // resource-agnostic and the destination keeps TanStack Router's typed links.
  onCancel: () => void;
  isPending?: boolean;
}

// Generic admin form, the write half of the CRUD engine (Wave 1). Owns the
// Mantine card shell and submit/cancel; per resource you pass the generated Zod
// validator, typed default values, and a field list. The same component backs
// both create and edit — the route decides which mutation `onSubmit` runs and
// seeds `defaultValues` accordingly.
export function CrudForm<T extends Record<string, unknown>>({
  fields,
  schema,
  defaultValues,
  onSubmit,
  submitLabel,
  onCancel,
  isPending,
}: CrudFormProps<T>) {
  const { t } = useTranslation();

  const form = useAppForm({
    defaultValues,
    validators: { onSubmit: schema },
    onSubmit: ({ value }) => {
      onSubmit(value);
    },
  });

  return (
    <Card withBorder p="xl" maw={720}>
      <Stack
        component="form"
        onSubmit={(event) => {
          event.preventDefault();
          form.handleSubmit();
        }}
      >
        {fields.map((spec) => (
          <form.AppField key={spec.name} name={spec.name}>
            {(field) => {
              switch (spec.type) {
                case "textarea":
                  return <field.TextareaField label={spec.label} required={spec.required} />;
                case "select":
                  return (
                    <field.SelectField
                      label={spec.label}
                      data={spec.options ?? []}
                      required={spec.required}
                    />
                  );
                case "datetime":
                  return <field.DateTimeField label={spec.label} required={spec.required} />;
                case "number":
                  return <field.NumberField label={spec.label} required={spec.required} />;
                case "checkbox":
                  return <field.CheckboxField label={spec.label} />;
                case "multiselect":
                  return (
                    <field.MultiSelectField
                      label={spec.label}
                      data={spec.options ?? []}
                      required={spec.required}
                    />
                  );
                default:
                  return <field.TextField label={spec.label} required={spec.required} />;
              }
            }}
          </form.AppField>
        ))}

        <Group justify="flex-end" mt="md">
          <Button variant="default" onClick={onCancel}>
            {t(($) => $.admin.crud.cancel)}
          </Button>
          <Button type="submit" loading={isPending}>
            {submitLabel}
          </Button>
        </Group>
      </Stack>
    </Card>
  );
}
