import type { i18n as I18n } from "i18next";
import { useTranslation } from "react-i18next";
import type { CrudColumnDef } from "./CrudList";
import type { CrudFieldSpec } from "./CrudForm";
import type { ColumnDescriptor, FieldDescriptor, FieldKind } from "./fieldDescriptors";

// What a resource file may still say by hand about its generated fields: the
// Exceptions the contract cannot express.
export interface FieldOverrides<T> {
  // Contract fields the form never shows (ids deferred to a later wave, fields
  // The server derives).
  hide?: (keyof T & string)[];
  // Show the fields in this order (names left out keep their contract order,
  // After the listed ones).
  order?: (keyof T & string)[];
  // Per-field changes on top of the generated spec: a custom widget (`type`),
  // A stricter `required` than the schema, enum/multiselect `options`.
  set?: Partial<Record<keyof T & string, Partial<Omit<CrudFieldSpec<T>, "name">>>>;
}

// A generated key the catalog lacks would render as the raw key; fail the render
// Instead so a contract change without its label cannot ship quietly.
function labelFor(i18n: I18n, key: string): string {
  if (!i18n.exists(key)) {
    throw new Error(`missing i18n label ${key}`);
  }
  return i18n.t(key);
}

const widgets: Record<FieldKind, CrudFieldSpec<never>["type"]> = {
  text: undefined,
  select: "select",
  datetime: "datetime",
  boolean: "checkbox",
  number: "number",
  multiselect: "multiselect",
};

// Turns the generated descriptors of a resource into the field list CrudForm
// Renders: labels resolve through i18n here (the descriptors carry only keys),
// Then the resource's overrides apply. A hook, like every i18n-reading helper.
export function useGeneratedFields<T>(
  descriptors: FieldDescriptor[],
  overrides: FieldOverrides<T> = {},
): CrudFieldSpec<T>[] {
  const { i18n } = useTranslation();
  const label = (key: string) => labelFor(i18n, key);
  const hidden = new Set<string>(overrides.hide);
  const order = overrides.order ?? [];

  const specs = descriptors
    .filter((descriptor) => !hidden.has(descriptor.name))
    .map((descriptor): CrudFieldSpec<T> => {
      const spec: CrudFieldSpec<T> = {
        name: descriptor.name as keyof T & string,
        label: label(descriptor.labelKey),
      };
      const type = widgets[descriptor.kind];
      if (type) {
        spec.type = type;
      }
      if (descriptor.required) {
        spec.required = true;
      }
      if (descriptor.enumValues) {
        spec.options = descriptor.enumValues.map((option) => ({
          value: option.value,
          label: label(option.labelKey),
        }));
      }
      return { ...spec, ...overrides.set?.[spec.name] };
    });

  const rank = (name: string) => {
    const index = order.indexOf(name as keyof T & string);
    return index === -1 ? order.length : index;
  };
  // Array.sort is stable, so unlisted fields keep the contract's order.
  return specs.toSorted((a, b) => rank(a.name) - rank(b.name));
}

// The plain columns of a list: a header from the generated label and the field
// Read off the row as-is. A route lists the names it shows, in order, and adds
// Its special cells (links, joins, badges) by hand around them.
export function useGeneratedColumns<T extends object>(
  descriptors: ColumnDescriptor[],
  names: (keyof T & string)[],
): CrudColumnDef<T>[] {
  const { i18n } = useTranslation();
  return names.map((name) => {
    const descriptor = descriptors.find((candidate) => candidate.name === name);
    if (!descriptor) {
      throw new Error(`no generated column ${name}`);
    }
    return { accessorKey: name, header: labelFor(i18n, descriptor.labelKey) };
  });
}
