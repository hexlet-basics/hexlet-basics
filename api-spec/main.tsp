import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi";
import "@typespec/openapi3";

using Http;
using Rest;
using OpenAPI;

/**
 * Hexlet Basics API.
 *
 * Naming: the domain concept is `Course`. The legacy Rails app called it
 * `Language`; we keep `Language` ONLY in URL paths (e.g. `/languages`) for
 * backward compatibility, but every type/operation/schema uses `Course`.
 */
@service(#{ title: "Hexlet Basics API" })
namespace HexletBasics;

/** How a course is meant to be learned. */
enum CourseLearnAs {
  firstLanguage: "first_language",
  secondLanguage: "second_language",
}

/** Publication/readiness state of a course. */
enum CourseProgress {
  completed: "completed",
  inDevelopment: "in_development",
  draft: "draft",
}

/** A course (legacy: `Language`). */
model Course {
  id: int32;
  slug: string;
  name: string | null;
  learnAs: CourseLearnAs | null;
  progress: CourseProgress | null;
  membersCount: int32;
  lessonsCount: int32;
  categoryId: int32 | null;
}

/** A catalog entry as shown on the courses index (legacy: LanguageLandingPageForLists). */
model CourseCatalogItem {
  id: int32;
  slug: string;
  header: string | null;
  name: string | null;
  locale: string | null;
  membersCount: int32;

  /** Estimated duration in hours (derived from lesson count). */
  duration: int32;

  /** Course cover image URL; null until assets are re-uploaded. */
  coverUrl: string | null;

  course: Course;
}

@route("/languages")
@tag("Courses")
interface Courses {
  /** List the published course catalog. */
  @get
  @operationId("listCourses")
  list(): CourseCatalogItem[];
}

// ---------------------------------------------------------------------------
// Admin
//
// The admin surface is a generated CRUD engine (see docs/FRONTEND_PLAN.md
// Wave 1): every resource exposes list/get/create/update/delete with a
// consistent shape so the contract drives both the ogen server validators and
// the hey-api zod client. URLs keep their legacy `/admin/language_*` paths for
// backward compatibility (ADR-0002); every model uses the domain name `Course`.
// ---------------------------------------------------------------------------

/** A page of results. Generic envelope reused by every admin list. */
@friendlyName("{name}Page", T)
model Page<T> {
  items: T[];

  /** Total rows across all pages (before pagination). */
  total: int32;

  /** 1-based page number this response represents. */
  page: int32;

  /** Page size used for this response. */
  perPage: int32;
}

/**
 * Field-level validation errors, keyed by field name (each value is the list
 * of messages for that field). Returned when a write fails validation —
 * including constraints the schema cannot express, like uniqueness.
 */
@error
model ValidationError {
  @statusCode statusCode: 422;
  errors: Record<string[]>;
}

/** A course category (legacy table: `language_categories`). */
model CourseCategory {
  id: int32;
  slug: string | null;
  name: string | null;
  header: string | null;
  description: string | null;
  locale: string | null;
  createdAt: utcDateTime;
}

/** Writable fields for creating or updating a course category. */
model CourseCategoryInput {
  @minLength(1) name: string;
  @minLength(1) header: string;
  @minLength(1) slug: string;
  description: string | null;
}

@route("/admin/language_categories")
@tag("Admin: Course Categories")
interface AdminCourseCategories {
  /** List course categories (paginated). */
  @get
  @operationId("adminListCourseCategories")
  list(@query page?: int32, @query perPage?: int32): Page<CourseCategory>;

  /** Get a single course category. */
  @get
  @operationId("adminGetCourseCategory")
  get(@path id: int32): CourseCategory;

  /** Create a course category. */
  @post
  @operationId("adminCreateCourseCategory")
  create(@body body: CourseCategoryInput): {
    @statusCode statusCode: 201;
    @body category: CourseCategory;
  } | ValidationError;

  /** Update a course category. */
  @put
  @operationId("adminUpdateCourseCategory")
  update(@path id: int32, @body body: CourseCategoryInput): CourseCategory | ValidationError;

  /** Delete a course category. */
  @delete
  @operationId("adminDeleteCourseCategory")
  delete(@path id: int32): {
    @statusCode statusCode: 204;
  };
}
