import type {
  ColumnDescriptor,
  EnumValueDescriptor,
  FieldDescriptor,
  FieldKind,
} from "../src/components/admin/fieldDescriptors.ts";

// The pure half of the admin field generator (the file half is
// gen-admin-fields.ts): from a parsed OpenAPI document and a resource's config it
// computes the descriptors of its form fields and list columns. Kept free of
// node APIs so a test can run it over a small fixture.

interface Schema {
  $ref?: string;
  type?: string;
  format?: string;
  enum?: string[];
  nullable?: boolean;
  allOf?: Schema[];
  items?: Schema;
  required?: string[];
  properties?: Record<string, Schema>;
}

export interface OpenApi {
  components: { schemas: Record<string, Schema> };
}

export interface ResourceConfig {
  // File name and export prefix of the generated module.
  name: string;
  // The `models.attributes.<model>` i18n scope of the resource's labels.
  model: string;
  // The request body schema (form fields) and the read model (list columns).
  input: string;
  item: string;
  // Field -> scope, for a field whose label lives under another scope than
  // `model` (the legacy catalog keeps the shared ones under `base`).
  labelScopes?: Record<string, string>;
  // Input properties that are never a form field (ids the server derives or a
  // later wave fills); no descriptor, so no label is owed for them.
  skipFields?: string[];
}

export interface ResourceDescriptors {
  model: string;
  fields: FieldDescriptor[];
  columns: ColumnDescriptor[];
}

function schemaNamed(spec: OpenApi, name: string): Schema {
  const schema = spec.components.schemas[name];
  if (!schema) {
    throw new Error(`schema ${name} is not in the OpenAPI document`);
  }
  return schema;
}

// A property is a plain schema, a bare $ref, or `allOf: [$ref]` + nullable (how
// TypeSpec spells a nullable reference). Returns the schema that carries the
// shape, plus the nullability read off whichever level declared it.
function resolve(spec: OpenApi, schema: Schema): { shape: Schema; nullable: boolean } {
  let nullable = schema.nullable === true;
  let current = schema;
  for (;;) {
    if (current.$ref) {
      current = schemaNamed(spec, current.$ref.replace("#/components/schemas/", ""));
    } else if (current.allOf?.length === 1 && current.allOf[0]) {
      current = current.allOf[0];
    } else {
      return { shape: current, nullable: nullable || current.nullable === true };
    }
    nullable ||= current.nullable === true;
  }
}

export function snakeCase(name: string): string {
  return name.replaceAll(/[A-Z]/g, (letter) => `_${letter.toLowerCase()}`);
}

function kindOf(shape: Schema): FieldKind | undefined {
  if (shape.enum) {
    return "select";
  }
  if (shape.type === "string") {
    return shape.format === "date-time" ? "datetime" : "text";
  }
  if (shape.type === "boolean") {
    return "boolean";
  }
  if (shape.type === "integer" || shape.type === "number") {
    return "number";
  }
  if (shape.type === "array" && shape.items?.type && !shape.items.enum) {
    return "multiselect";
  }
  return undefined;
}

function labelKeys(config: ResourceConfig, field: string) {
  const scope = config.labelScopes?.[field] ?? config.model;
  const base = `models.attributes.${scope}.${snakeCase(field)}`;
  return {
    labelKey: base,
    enumKey: (value: string) => `${base}/values.${value}`,
  };
}

function describeProperty(spec: OpenApi, config: ResourceConfig, name: string, schema: Schema) {
  const { shape, nullable } = resolve(spec, schema);
  const kind = kindOf(shape);
  if (!kind) {
    return;
  }
  const { labelKey, enumKey } = labelKeys(config, name);
  const enumValues: EnumValueDescriptor[] | undefined = shape.enum?.map((value) => ({
    value,
    labelKey: enumKey(value),
  }));
  return { name, kind, nullable, labelKey, ...(enumValues && { enumValues }) };
}

export function describeResource(spec: OpenApi, config: ResourceConfig): ResourceDescriptors {
  const input = schemaNamed(spec, config.input);
  const item = schemaNamed(spec, config.item);
  const required = new Set(input.required ?? []);

  const skipped = new Set(config.skipFields);
  const inputProperties = Object.entries(input.properties ?? {}).filter(
    ([name]) => !skipped.has(name),
  );
  const fields = inputProperties.map(([name, schema]) => {
    const described = describeProperty(spec, config, name, schema);
    if (!described) {
      throw new Error(`${config.input}.${name}: no form field kind for this schema`);
    }
    return { ...described, required: required.has(name) && !described.nullable };
  });

  // Columns skip what a cell cannot show as-is (nested objects, lists of objects):
  // those stay a hand-written cell in the route.
  const columns = Object.entries(item.properties ?? {}).flatMap(([name, schema]) => {
    const described = describeProperty(spec, config, name, schema);
    return described && described.kind !== "multiselect" ? [described] : [];
  });

  return { model: config.model, fields, columns };
}

const header = (name: string) =>
  `// DO NOT EDIT. Generated by scripts/gen-admin-fields.ts from api-spec/dist/openapi.yaml\n` +
  `// (\`mise run gen:admin\`, part of \`mise run gen\`); hand-written exceptions live in\n` +
  `// src/components/admin/resources/${name}.ts.\n`;

// The module of one resource: two typed constant arrays. Plain JSON literals, so
// the output depends on nothing but the descriptors (oxfmt then lays it out).
export function renderModule(name: string, data: ResourceDescriptors): string {
  return (
    header(name) +
    `import type { ColumnDescriptor, FieldDescriptor } from "../fieldDescriptors";\n\n` +
    `export const ${name}Fields: FieldDescriptor[] = ${JSON.stringify(data.fields, null, 2)};\n\n` +
    `export const ${name}Columns: ColumnDescriptor[] = ${JSON.stringify(data.columns, null, 2)};\n`
  );
}
