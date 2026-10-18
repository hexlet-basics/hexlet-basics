import { mkdirSync, readFileSync, readdirSync, rmSync, writeFileSync } from "node:fs";
import path from "node:path";
import { describeResource, parseOpenApi, renderModule } from "./admin-fields.ts";
import type { ResourceConfig } from "./admin-fields.ts";

// Writes src/components/admin/generated/<resource>.ts for every admin resource
// Below, from api-spec/dist/openapi.yaml. Run by `mise run gen:admin` (part of
// `mise run gen`, so `gen:check` fails when the committed output is stale).
// Deliberately a standalone script, not a hey-api plugin: the descriptors are
// Ours, the client generator stays on its stock configuration.

// Labels live under `models.attributes.<model>`; the legacy catalog keeps the
// shared ones (name, slug, description) under `base`.
const base = { description: "base", name: "base", slug: "base" };

const resources: ResourceConfig[] = [
  { input: "BannerInput", item: "Banner", model: "banner", name: "banner" },
  {
    input: "BlogPostInput",
    item: "BlogPost",
    labelScopes: base,
    model: "blog_post",
    name: "blogPost",
    skipFields: ["coverAttachmentId"],
  },
  {
    input: "CourseInput",
    item: "Course",
    model: "course",
    name: "course",
    skipFields: ["repositoryUrl", "coverAttachmentId"],
  },
  {
    input: "CourseCategoryInput",
    item: "CourseCategory",
    model: "course_category",
    name: "courseCategory",
  },
  {
    input: "CourseLandingPageInput",
    item: "CourseLandingPage",
    model: "course_landing_page",
    name: "courseLandingPage",
    skipFields: ["outcomesImageAttachmentId"],
  },
  { input: "ReviewInput", item: "Review", model: "review", name: "review" },
  { input: "StaffMemberInput", item: "StaffMember", model: "staff_member", name: "staffMember" },
  {
    input: "RoleInput",
    item: "StaffRoleDetail",
    labelScopes: base,
    model: "staff_role",
    name: "staffRole",
  },
  { input: "UserInput", item: "UserCrud", model: "user", name: "user" },
];

const root = path.join(import.meta.dirname, "..");
const specPath = path.join(root, "api-spec/dist/openapi.yaml");
const outDir = path.join(root, "src/components/admin/generated");

const spec = parseOpenApi(readFileSync(specPath, "utf-8"));

// Start from an empty directory so a resource dropped from the list above does
// Not leave its module behind.
mkdirSync(outDir, { recursive: true });
for (const file of readdirSync(outDir)) {
  rmSync(path.join(outDir, file));
}

for (const config of resources) {
  writeFileSync(
    path.join(outDir, `${config.name}.ts`),
    renderModule(config.name, describeResource(spec, config)),
  );
}
