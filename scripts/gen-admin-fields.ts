import { mkdirSync, readdirSync, readFileSync, rmSync, writeFileSync } from "node:fs";
import { join } from "node:path";
import { parse } from "yaml";
import { describeResource, renderModule } from "./admin-fields.ts";
import type { OpenApi, ResourceConfig } from "./admin-fields.ts";

// Writes src/components/admin/generated/<resource>.ts for every admin resource
// below, from api-spec/dist/openapi.yaml. Run by `mise run gen:admin` (part of
// `mise run gen`, so `gen:check` fails when the committed output is stale).
// Deliberately a standalone script, not a hey-api plugin: the descriptors are
// ours, the client generator stays on its stock configuration.

// Labels live under `models.attributes.<model>`; the legacy catalog keeps the
// shared ones (name, slug, description) under `base`.
const base = { name: "base", slug: "base", description: "base" };

const resources: ResourceConfig[] = [
  { name: "banner", model: "banner", input: "BannerInput", item: "Banner" },
  {
    name: "blogPost",
    model: "blog_post",
    input: "BlogPostInput",
    item: "BlogPost",
    skipFields: ["coverAttachmentId"],
    labelScopes: base,
  },
  {
    name: "course",
    model: "course",
    input: "CourseInput",
    item: "Course",
    skipFields: ["repositoryUrl", "coverAttachmentId"],
  },
  {
    name: "courseCategory",
    model: "course_category",
    input: "CourseCategoryInput",
    item: "CourseCategory",
  },
  {
    name: "courseLandingPage",
    model: "course_landing_page",
    input: "CourseLandingPageInput",
    item: "CourseLandingPage",
    skipFields: ["outcomesImageAttachmentId"],
  },
  { name: "review", model: "review", input: "ReviewInput", item: "Review" },
  { name: "staffMember", model: "staff_member", input: "StaffMemberInput", item: "StaffMember" },
  {
    name: "staffRole",
    model: "staff_role",
    input: "RoleInput",
    item: "StaffRoleDetail",
    labelScopes: base,
  },
  { name: "user", model: "user", input: "UserInput", item: "UserCrud" },
];

const root = join(import.meta.dirname, "..");
const specPath = join(root, "api-spec/dist/openapi.yaml");
const outDir = join(root, "src/components/admin/generated");

const spec = parse(readFileSync(specPath, "utf8")) as OpenApi;

// Start from an empty directory so a resource dropped from the list above does
// not leave its module behind.
mkdirSync(outDir, { recursive: true });
for (const file of readdirSync(outDir)) {
  rmSync(join(outDir, file));
}

for (const config of resources) {
  writeFileSync(
    join(outDir, `${config.name}.ts`),
    renderModule(config.name, describeResource(spec, config)),
  );
}
