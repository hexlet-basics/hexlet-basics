import { useMemo } from "react";
import { useTranslation } from "react-i18next";
import { type ZodType, z } from "zod";
import type { CrudFieldSpec } from "@/components/admin/CrudForm";

// Derive the CrudForm field list straight from the generated `zXxxInput` schema,
// so the mechanical, easy-to-desync parts (which fields exist, their order,
// enum→select with its exact option values, datetime detection, required vs
// nullable) are never hand-maintained. Only the human bits stay per-resource:
// the i18n label namespace and which string fields render as a textarea.
//
// Structure comes from `z.toJSONSchema` (a stable, documented API) rather than
// poking zod internals: enum → { enum: [...] }, `z.iso.datetime()` →
// { format: "date-time" }, `.nullable()` → { anyOf: [..., { type: "null" }] }.

interface JsonProp {
  type?: string | string[];
  format?: string;
  enum?: string[];
  anyOf?: JsonProp[];
}
interface JsonObject {
  properties?: Record<string, JsonProp>;
}

// The non-null branch carries the real type; a null branch marks the field
// optional.
const coreOf = (prop: JsonProp): JsonProp =>
  prop.anyOf?.find((branch) => branch.type !== "null") ?? prop;
const isNullable = (prop: JsonProp): boolean =>
  prop.anyOf?.some((branch) => branch.type === "null") ?? false;

// Field names are camelCase (startsAt); the legacy i18n attribute keys are snake
// (starts_at).
const toSnake = (name: string): string =>
  name.replace(/[A-Z]/g, (c) => `_${c.toLowerCase()}`);

export interface SchemaFieldsOptions {
  // i18n attributes namespace, e.g. "banner" → models.attributes.banner.<field>.
  namespace: string;
  // Field names rendered as a textarea instead of a single-line input.
  textarea?: string[];
}

export function useSchemaFields<T>(
  schema: ZodType,
  { namespace, textarea = [] }: SchemaFieldsOptions,
): CrudFieldSpec<T>[] {
  const { t } = useTranslation();

  return useMemo(() => {
    // Labels are looked up by a key COMPUTED from the field name, but the project
    // enables i18next's typed selector (`t(($) => $.a.b)`), whose only overload
    // rejects string keys. Runtime `t` still resolves plain string keys, so cast
    // to the string-key signature here — the one spot dynamic keys are needed.
    const translate = t as unknown as (key: string) => string;
    const json = z.toJSONSchema(schema, {
      unrepresentable: "any",
    }) as JsonObject;
    const attr = (field: string) =>
      translate(`models.attributes.${namespace}.${toSnake(field)}`);

    let autoFocusTaken = false;
    return Object.entries(json.properties ?? {}).map(([name, prop]) => {
      const core = coreOf(prop);

      let type: CrudFieldSpec<T>["type"] = "text";
      let options: CrudFieldSpec<T>["options"];
      if (Array.isArray(core.enum)) {
        type = "select";
        options = core.enum.map((value) => ({
          value,
          label: translate(
            `models.attributes.${namespace}.${toSnake(name)}/values.${value}`,
          ),
        }));
      } else if (core.format === "date-time") {
        type = "datetime";
      } else if (textarea.includes(name)) {
        type = "textarea";
      }

      // Focus the first free-text field (skip selects/dates), matching how the
      // hand-written forms behaved.
      const autoFocus =
        !autoFocusTaken && (type === "text" || type === "textarea");
      if (autoFocus) autoFocusTaken = true;

      return {
        name: name as keyof T & string,
        label: attr(name),
        type,
        options,
        required: !isNullable(prop),
        autoFocus,
      };
    });
  }, [schema, namespace, textarea, t]);
}
