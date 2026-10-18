import { describe, expect, it } from "vitest";
import en from "@/locales/en/translation";
import es from "@/locales/es/translation";
import ru from "@/locales/ru/translation";
import type { FieldDescriptor } from "@/components/admin/fieldDescriptors";

// The generated descriptors name their labels by i18n key, and nothing in the
// source mentions those keys as literals, so `lint:i18n` cannot see them. This
// is the check that a form field added in the contract has a label in every
// catalog (list columns are checked where a route uses them, see
// useGeneratedColumns).
const modules = import.meta.glob<Record<string, unknown>>("@/components/admin/generated/*.ts", {
  eager: true,
});

const catalogs = { en, es, ru };

function lookup(catalog: unknown, key: string): unknown {
  let node = (catalog as { translation: unknown }).translation;
  for (const part of key.split(".")) {
    node =
      typeof node === "object" && node !== null
        ? (node as Record<string, unknown>)[part]
        : undefined;
  }
  return node;
}

describe("generated admin descriptors", () => {
  const keys = Object.values(modules).flatMap((module) =>
    Object.entries(module)
      .filter(([exportName]) => exportName.endsWith("Fields"))
      .flatMap(([, descriptors]) =>
        (descriptors as FieldDescriptor[]).flatMap((descriptor) =>
          [
            descriptor.labelKey,
            (descriptor.enumValues ?? []).map((value) => value.labelKey),
          ].flat(),
        ),
      ),
  );

  it("covers the admin resources", () => {
    expect.hasAssertions();
    expect(Object.keys(modules).length).toBeGreaterThan(0);
  }, 60_000);

  it.each(Object.entries(catalogs))(
    "has every label key in the %s catalog",
    (_locale, catalog) => {
      expect.hasAssertions();
      const missing = [...new Set(keys)].filter((key) => typeof lookup(catalog, key) !== "string");
      expect(missing).toStrictEqual([]);
    },
    60_000,
  );
});
