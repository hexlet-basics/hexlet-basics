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
