// The shape of what `scripts/gen-admin-fields.ts` writes into `./generated/*.ts`:
// One descriptor per form field (from the resource's `XxxInput` schema) and per
// List column (from the read model), derived from the OpenAPI contract so a
// Field added in TypeSpec shows up here after `mise run gen`.

export type FieldKind =
  | "text"
  | "select"
  | "datetime"
  | "boolean"
  | "number"
  // An array of scalars; the generated descriptor cannot know the options.
  | "multiselect";

export interface EnumValueDescriptor {
  value: string;
  // I18n key of the value's label: `models.attributes.<model>.<field>/values.<value>`.
  labelKey: string;
}

export interface FieldDescriptor {
  // The contract's camelCase property name.
  name: string;
  kind: FieldKind;
  // In the schema's `required` list and not nullable: the form must send a value.
  required: boolean;
  nullable: boolean;
  // I18n key of the label: `models.attributes.<model>.<snake_case_field>`.
  labelKey: string;
  // Set for `kind: "select"`.
  enumValues?: EnumValueDescriptor[];
}

export interface ColumnDescriptor {
  name: string;
  kind: FieldKind;
  nullable: boolean;
  labelKey: string;
  enumValues?: EnumValueDescriptor[];
}
