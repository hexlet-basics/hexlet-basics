import type { i18n as I18n } from "i18next";
import { useTranslation } from "react-i18next";
import type { CrudColumnDef } from "./CrudList";
import type { CrudFieldSpec } from "./CrudForm";
import type { ColumnDescriptor, FieldDescriptor, FieldKind } from "./fieldDescriptors";

// What a resource file may still say by hand about its generated fields: the
// exceptions the contract cannot express.
export interface FieldOverrides<T> {
  // Contract fields the form never shows (ids deferred to a later wave, fields
  // the server derives).
  hide?: (keyof T & string)[];
  // Show the fields in this order (names left out keep their contract order,
  // after the listed ones).
  order?: (keyof T & string)[];
  // Per-field changes on top of the generated spec: a custom widget (`type`),
  // a stricter `required` than the schema, enum/multiselect `options`.
  set?: Partial<Record<keyof T & string, Partial<Omit<CrudFieldSpec<T>, "name">>>>;
}

// A generated key the catalog lacks would render as the raw key; fail the render
// instead so a contract change without its label cannot ship quietly.
function labelFor(i18n: I18n, key: string): string {
  if (!i18n.exists(key)) {
    throw new Error(`missing i18n label ${key}`);
  }
  return i18n.t(key);
}

const widgets: Record<FieldKind, CrudFieldSpec<never>["type"]> = {
  boolean: "checkbox",
  datetime: "datetime",
  multiselect: "multiselect",
  number: "number",
  select: "select",
  text: undefined,
};

// Turns the generated descriptors of a resource into the field list CrudForm
// renders: labels resolve through i18n here (the descriptors carry only keys),
// then the resource's overrides apply. A hook, like every i18n-reading helper.
export function useGeneratedFields<T>(
  descriptors: FieldDescriptor[],
  overrides: FieldOverrides<T> = {},
): CrudFieldSpec<T>[] {
  const { i18n } = useTranslation();
  const label = (key: string): string => labelFor(i18n, key);
  const hidden = new Set<string>(overrides.hide);
  const order = overrides.order ?? [];

  const specs = descriptors
    .filter((descriptor) => !hidden.has(descriptor.name))
    .map((descriptor): CrudFieldSpec<T> => {
      const spec: CrudFieldSpec<T> = {
        label: label(descriptor.labelKey),
        name: descriptor.name as keyof T & string,
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
          label: label(option.labelKey),
          value: option.value,
        }));
      }
      return Object.assign(spec, overrides.set?.[spec.name]);
    });

  const rank = (name: string): number => {
    const index = order.indexOf(name as keyof T & string);
    return index === -1 ? order.length : index;
  };
  // Array.sort is stable, so unlisted fields keep the contract's order.
  return specs.toSorted((left, right) => rank(left.name) - rank(right.name));
}

// The plain columns of a list: a header from the generated label and the field
// read off the row as-is. A route lists the names it shows, in order, and adds
// its special cells (links, joins, badges) by hand around them.
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
