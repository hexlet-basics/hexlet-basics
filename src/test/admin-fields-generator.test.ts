import { parse } from "yaml";
import { describe, expect, it } from "vitest";
import { describeResource, renderModule } from "../../scripts/admin-fields";

// A miniature contract in the shape the TypeSpec pipeline emits: a request body
// (`WidgetInput`), the read model (`Widget`), and the schemas they $ref.
const spec = parse(`
components:
  schemas:
    WidgetState:
      type: string
      enum: [draft, published]
    Widget:
      type: object
      required: [id, title, state, startsAt, owner, tags]
      properties:
        id: { type: integer, format: int32 }
        title: { type: string }
        state: { $ref: '#/components/schemas/WidgetState' }
        startsAt: { type: string, format: date-time, nullable: true }
        owner: { $ref: '#/components/schemas/Owner' }
        tags: { type: array, items: { type: string } }
    Owner:
      type: object
      properties:
        name: { type: string }
    WidgetInput:
      type: object
      required: [title, state, startsAt, pinned, seats, nickname, locales, richBody]
      properties:
        title: { type: string, minLength: 1 }
        state:
          allOf:
            - $ref: '#/components/schemas/WidgetState'
          nullable: true
        startsAt: { type: string, format: date-time, nullable: true }
        pinned: { type: boolean }
        seats: { type: integer, format: int32 }
        nickname: { type: string, nullable: true }
        locales: { type: array, items: { type: string } }
        optionalNote: { type: string }
        richBody: { type: string }
`);

const widget = { name: "widget", model: "widget", input: "WidgetInput", item: "Widget" };

describe("admin field generator", () => {
  it("describes each input property by kind, required and nullable", () => {
    const { fields } = describeResource(spec, widget);

    expect(fields.map((f) => [f.name, f.kind, f.required, f.nullable])).toEqual([
      ["title", "text", true, false],
      ["state", "select", false, true],
      ["startsAt", "datetime", false, true],
      ["pinned", "boolean", true, false],
      ["seats", "number", true, false],
      ["nickname", "text", false, true],
      ["locales", "multiselect", true, false],
      ["optionalNote", "text", false, false],
      ["richBody", "text", true, false],
    ]);
  });

  it("labels fields by the models.attributes convention, snake-casing the name", () => {
    const { fields } = describeResource(spec, widget);

    expect(fields.find((f) => f.name === "startsAt")?.labelKey).toBe(
      "models.attributes.widget.starts_at",
    );
    expect(fields.find((f) => f.name === "richBody")?.labelKey).toBe(
      "models.attributes.widget.rich_body",
    );
  });

  it("carries enum values with their /values label keys, through $ref and allOf", () => {
    const { fields } = describeResource(spec, widget);

    expect(fields.find((f) => f.name === "state")?.enumValues).toEqual([
      { value: "draft", labelKey: "models.attributes.widget.state/values.draft" },
      { value: "published", labelKey: "models.attributes.widget.state/values.published" },
    ]);
  });

  it("lets a resource point a field's label at another model scope", () => {
    const { fields } = describeResource(spec, { ...widget, labelScopes: { title: "base" } });

    expect(fields.find((f) => f.name === "title")?.labelKey).toBe("models.attributes.base.title");
  });

  it("leaves out the input fields a resource declares it has no form field for", () => {
    const { fields } = describeResource(spec, { ...widget, skipFields: ["seats", "pinned"] });

    expect(fields.map((f) => f.name)).not.toContain("seats");
    expect(fields.map((f) => f.name)).not.toContain("pinned");
    expect(fields.map((f) => f.name)).toContain("title");
  });

  it("describes the scalar properties of the read model as list columns", () => {
    const { columns } = describeResource(spec, widget);

    expect(columns.map((c) => [c.name, c.kind])).toEqual([
      ["id", "number"],
      ["title", "text"],
      ["state", "select"],
      ["startsAt", "datetime"],
    ]);
  });

  it("fails loudly on a schema the contract does not define", () => {
    expect(() => describeResource(spec, { ...widget, input: "NopeInput" })).toThrow(/NopeInput/);
  });

  it("renders a module that names the generator and is stable across runs", () => {
    const data = describeResource(spec, widget);
    const source = renderModule("widget", data);

    expect(source).toContain("DO NOT EDIT");
    expect(source).toContain("export const widgetFields");
    expect(source).toContain("export const widgetColumns");
    expect(renderModule("widget", describeResource(spec, widget))).toBe(source);
  });
});
