import { describe, expect, it } from "vitest";
import { describeResource, parseOpenApi, renderModule } from "../../scripts/admin-fields";

// A miniature contract in the shape the TypeSpec pipeline emits: a request body
// (`WidgetInput`), the read model (`Widget`), and the schemas they $ref.
const spec = parseOpenApi(`
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

const widget = { input: "WidgetInput", item: "Widget", model: "widget", name: "widget" };

describe("admin field generator", () => {
  it("describes each input property by kind, required and nullable", () => {
    expect.hasAssertions();
    const { fields } = describeResource(spec, widget);

    expect(
      fields.map((field) => [field.name, field.kind, field.required, field.nullable]),
    ).toStrictEqual([
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
  }, 60_000);

  it("labels fields by the models.attributes convention, snake-casing the name", () => {
    expect.hasAssertions();
    const { fields } = describeResource(spec, widget);

    expect(fields.find((field) => field.name === "startsAt")?.labelKey).toBe(
      "models.attributes.widget.starts_at",
    );
    expect(fields.find((field) => field.name === "richBody")?.labelKey).toBe(
      "models.attributes.widget.rich_body",
    );
  }, 60_000);

  it("carries enum values with their /values label keys, through $ref and allOf", () => {
    expect.hasAssertions();
    const { fields } = describeResource(spec, widget);

    expect(fields.find((field) => field.name === "state")?.enumValues).toStrictEqual([
      { labelKey: "models.attributes.widget.state/values.draft", value: "draft" },
      { labelKey: "models.attributes.widget.state/values.published", value: "published" },
    ]);
  }, 60_000);

  it("lets a resource point a field's label at another model scope", () => {
    expect.hasAssertions();
    const { fields } = describeResource(spec, { ...widget, labelScopes: { title: "base" } });

    expect(fields.find((field) => field.name === "title")?.labelKey).toBe(
      "models.attributes.base.title",
    );
  }, 60_000);

  it("leaves out the input fields a resource declares it has no form field for", () => {
    expect.hasAssertions();
    const { fields } = describeResource(spec, { ...widget, skipFields: ["seats", "pinned"] });

    expect(fields.map((field) => field.name)).not.toContain("seats");
    expect(fields.map((field) => field.name)).not.toContain("pinned");
    expect(fields.map((field) => field.name)).toContain("title");
  }, 60_000);

  it("describes the scalar properties of the read model as list columns", () => {
    expect.hasAssertions();
    const { columns } = describeResource(spec, widget);

    expect(columns.map((column) => [column.name, column.kind])).toStrictEqual([
      ["id", "number"],
      ["title", "text"],
      ["state", "select"],
      ["startsAt", "datetime"],
    ]);
  }, 60_000);

  it("fails loudly on a schema the contract does not define", () => {
    expect.hasAssertions();
    expect(() => describeResource(spec, { ...widget, input: "NopeInput" })).toThrow(/NopeInput/u);
  }, 60_000);

  it("renders a module that names the generator and is stable across runs", () => {
    expect.hasAssertions();
    const data = describeResource(spec, widget);
    const source = renderModule("widget", data);

    expect(source).toContain("DO NOT EDIT");
    expect(source).toContain("export const widgetFields");
    expect(source).toContain("export const widgetColumns");
    expect(renderModule("widget", describeResource(spec, widget))).toBe(source);
  }, 60_000);
});
