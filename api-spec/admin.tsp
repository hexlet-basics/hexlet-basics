import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi";

using Http;
using Rest;
using OpenAPI;

namespace HexletBasics;

// ===========================================================================
// ADMIN — the generated CRUD engine surface (FRONTEND_PLAN Wave 1/2).
//
// Every resource exposes a consistent list/get/create/update/delete subset so
// one engine (<CrudList> + <CrudForm>) drives them all. Admin URLs use the
// domain vocabulary; only the public routes keep their legacy paths (ADR-0002).
//
// Uniform response shape: writes return `T`; a failed one is ProblemDetails (a 422 carries `errors`). Index is
// `Page<T>` (paginated) via `...ListQuery`.
//
// FILE UPLOADS: covers / images do NOT ride inside the JSON create/update
// bodies (that would break the uniform JSON engine). They go through the shared
// multipart `POST /admin/attachments` below, which returns an id/url the form
// then submits as `*AttachmentId`. This is the ADR-0005 blob model and differs
// from legacy's inline multipart on purpose.
// ===========================================================================

/** A stored blob (cover image, outcomes image) returned by the uploader. */
model Attachment {
  id: int32;
  url: string;
  filename: string;
  contentType: string;
  byteSize: int64;
}

/** Multipart form for a single file upload. */
model AttachmentUploadForm {
  // `bytes`, not `File`: HttpPart<File> emits a requestBody Encoding Object
  // (`contentType: "*/*"`) that ogen cannot generate. The field schema is the
  // same binary string either way; the stored type is sniffed from the bytes.
  file: HttpPart<bytes>;
}

@route("/admin/attachments")
@tag("Admin: Attachments")
@useAuth(AdminSession)
interface AdminAttachments {
  /** Upload a file; returns the attachment to reference from a form. */
  @post
  // The handler writes to blob storage, so the audit middleware must not hold a
  // database transaction open around it (ADR-0020, "external calls" addendum).
  @extension("x-unit-of-work", "standalone")
  @extension("x-access", "any-staff") @operationId("adminUploadAttachment")
  upload(@header contentType: "multipart/form-data", @multipartBody body: AttachmentUploadForm): {
    @statusCode statusCode: 201;
    @body attachment: Attachment;
  } | ApiError;
}

// ---------------------------------------------------------------------------
// Course Categories.
//
// KEPT STABLE: this slice already has a Go handler + ent + generated client.
// The input intentionally stays minimal. Aligning it to the full serializer
// (locale, landingPageIds, nested qna items) is a tracked follow-up, not part
// of this pass — see the summary.
// ---------------------------------------------------------------------------

/** Writable fields for creating or updating a course category. */
model CourseCategoryInput {
  @minLength(1) name: string;
  @minLength(1) header: string;
  @minLength(1) slug: string;
  description: string | null;
}

@route("/admin/course_categories")
@tag("Admin: Course Categories")
@useAuth(AdminSession)
interface AdminCourseCategories {
  /** List course categories (paginated). */
  @get
  @extension("x-permission", #{ name: "course_categories.read", kind: "read" }) @operationId("adminListCourseCategories")
  list(...ListQuery): Page<CourseCategory> | ApiError;

  /** Get a single course category. A missing id surfaces as 404 via the
      central ent-error handler, not a typed union member. */
  @get
  @route("/{id}")
  @extension("x-permission", #{ name: "course_categories.read", kind: "read" }) @operationId("adminGetCourseCategory")
  get(@path id: int32): CourseCategory | ApiError;

  /** Create a course category. A uniqueness violation (name/header/slug) is a
      DB constraint, surfaced as 409 by the central ent-error handler. */
  @post
  @useAuth(AdminSession)
  @extension("x-permission", #{ name: "course_categories.create", kind: "write" }) @operationId("adminCreateCourseCategory")
  create(@body body: CourseCategoryInput): {
    @statusCode statusCode: 201;
    @body category: CourseCategory;
  } | ApiError;

  /** Update a course category. 404 (missing) and 409 (uniqueness) both flow
      through the central ent-error handler. */
  @put
  @useAuth(AdminSession)
  @route("/{id}")
  @extension("x-permission", #{ name: "course_categories.update", kind: "write" }) @operationId("adminUpdateCourseCategory")
  update(@path id: int32, @body body: CourseCategoryInput): CourseCategory | ApiError;

  /** Delete a course category. */
  @delete
  @useAuth(AdminSession)
  @route("/{id}")
  @extension("x-permission", #{ name: "course_categories.delete", kind: "write" }) @operationId("adminDeleteCourseCategory")
  delete(@path id: int32): NoContent | ApiError;
}

// ---------------------------------------------------------------------------
// Banners (legacy `/admin/banners`) — full CRUD.
// ---------------------------------------------------------------------------

/** The input of a banner: never its Publication State. A new banner is a draft and
    the state changes only through the publish, unpublish and archive operations. */
model BannerInput {
  background: BannerBackground;
  locale: BannerLocale;
  @minLength(1) body: string;
  url: string | null;
  startsAt: utcDateTime | null;
  finishesAt: utcDateTime | null;
}

@route("/admin/banners")
@tag("Admin: Banners")
@useAuth(AdminSession)
interface AdminBanners {
  /** List banners (paginated). */
  @get @extension("x-permission", #{ name: "banners.read", kind: "read" }) @operationId("adminListBanners") list(...SearchListQuery): Page<Banner> | ApiError;

  /** Get a single banner. A missing id surfaces as 404 via the central
      ent-error handler, not a typed union member. */
  @get @route("/{id}") @extension("x-permission", #{ name: "banners.read", kind: "read" }) @operationId("adminGetBanner") get(@path id: int32): Banner | ApiError;

  /** Create a banner. A body violating the schema (e.g. an empty `body`)
      surfaces as 400 via the central handler; banners carry no uniqueness
      constraint, so there is no 409 path. */
  @post @useAuth(AdminSession) @extension("x-permission", #{ name: "banners.create", kind: "write" }) @operationId("adminCreateBanner") create(@body body: BannerInput): {
    @statusCode statusCode: 201;
    @body banner: Banner;
  } | ApiError;

  /** Update a banner. A missing id surfaces as 404 via the central handler. */
  @put @useAuth(AdminSession) @route("/{id}") @extension("x-permission", #{ name: "banners.update", kind: "write" }) @operationId("adminUpdateBanner") update(@path id: int32, @body body: BannerInput): Banner | ApiError;

  /** Delete a banner. */
  @delete @useAuth(AdminSession) @route("/{id}") @extension("x-permission", #{ name: "banners.delete", kind: "write" }) @operationId("adminDeleteBanner") delete(@path id: int32): NoContent | ApiError;

  /** Publish a draft or archived banner. A state is not an update: each change of
      Publication State is its own operation with its own Permission. An illegal
      move, or repeating the current state, is a 409 with a localized `detail`. */
  @post @useAuth(AdminSession) @route("/{id}/publish") @extension("x-permission", #{ name: "banners.publish", kind: "write" }) @operationId("adminPublishBanner")
  publish(@path id: int32): Banner | ApiError;

  /** Take a published banner back to draft. */
  @post @useAuth(AdminSession) @route("/{id}/unpublish") @extension("x-permission", #{ name: "banners.unpublish", kind: "write" }) @operationId("adminUnpublishBanner")
  unpublish(@path id: int32): Banner | ApiError;

  /** Archive a draft or published banner. */
  @post @useAuth(AdminSession) @route("/{id}/archive") @extension("x-permission", #{ name: "banners.archive", kind: "write" }) @operationId("adminArchiveBanner")
  archive(@path id: int32): Banner | ApiError;
}

// ---------------------------------------------------------------------------
// Blog posts (legacy `/admin/blog_posts`) — full CRUD + related courses (set
// by hand or suggested by the AI job).
// ---------------------------------------------------------------------------

/** The input of a post: never its Publication State. A new post is a draft and
    the state changes only through the publish, unpublish and archive operations. */
model BlogPostInput {
  name: string | null;
  slug: string | null;
  description: string | null;

  /** Trusted editor HTML, stored and returned without sanitization. */
  richBody: string;

  /** Attachment id from `POST /admin/attachments`, or null to keep/clear. */
  coverAttachmentId: int32 | null;
}

/** The set of courses promoted on a blog post. */
model BlogPostRelatedCoursesInput {
  courseIds: int32[];
}

@route("/admin/blog_posts")
@tag("Admin: Blog Posts")
@useAuth(AdminSession)
interface AdminBlogPosts {
  @get @extension("x-permission", #{ name: "blog_posts.read", kind: "read" }) @operationId("adminListBlogPosts") list(...SearchListQuery): Page<BlogPost> | ApiError;

  // 404 comes from the shared ApiError response (a domain not-found), like
  // every other get.
  @get @route("/{id}") @extension("x-permission", #{ name: "blog_posts.read", kind: "read" }) @operationId("adminGetBlogPost") get(@path id: int32): BlogPost | ApiError;

  @post @useAuth(AdminSession) @extension("x-permission", #{ name: "blog_posts.create", kind: "write" }) @operationId("adminCreateBlogPost") create(@body body: BlogPostInput): {
    @statusCode statusCode: 201;
    @body post: BlogPost;
  } | ApiError;

  @put @useAuth(AdminSession) @route("/{id}") @extension("x-permission", #{ name: "blog_posts.update", kind: "write" }) @operationId("adminUpdateBlogPost") update(@path id: int32, @body body: BlogPostInput): BlogPost | ApiError;

  @delete @useAuth(AdminSession) @route("/{id}") @extension("x-permission", #{ name: "blog_posts.delete", kind: "write" }) @operationId("adminDeleteBlogPost") delete(@path id: int32): NoContent | ApiError;

  /** Set the related/promoted courses for a post by hand. */
  @post @useAuth(AdminSession) @route("/{id}/related_courses") @extension("x-permission", #{ name: "blog_posts.update", kind: "write" }) @operationId("adminSetBlogPostRelatedCourses")
  relatedCourses(@path id: int32, @body body: BlogPostRelatedCoursesInput): BlogPost | ApiError;

  /** Enqueue the AI pick of the post's related courses (legacy
      FindRelatedCoursesForBlogPostJob); the job replaces the set when done.
      A missing id surfaces as 404 via the central ent-error handler. */
  @post @useAuth(AdminSession) @route("/{id}/suggest_related_courses") @extension("x-permission", #{ name: "blog_posts.update", kind: "write" }) @operationId("adminSuggestBlogPostRelatedCourses")
  suggestRelatedCourses(@path id: int32): NoContent | ApiError;

  /** Publish a draft or archived post. A state is not an update: each change of
      Publication State is its own operation with its own Permission. An illegal
      move, or repeating the current state, is a 409 with a localized `detail`. */
  @post @useAuth(AdminSession) @route("/{id}/publish") @extension("x-permission", #{ name: "blog_posts.publish", kind: "write" }) @operationId("adminPublishBlogPost")
  publish(@path id: int32): BlogPost | ApiError;

  /** Take a published post back to draft. */
  @post @useAuth(AdminSession) @route("/{id}/unpublish") @extension("x-permission", #{ name: "blog_posts.unpublish", kind: "write" }) @operationId("adminUnpublishBlogPost")
  unpublish(@path id: int32): BlogPost | ApiError;

  /** Archive a draft or published post. */
  @post @useAuth(AdminSession) @route("/{id}/archive") @extension("x-permission", #{ name: "blog_posts.archive", kind: "write" }) @operationId("adminArchiveBlogPost")
  archive(@path id: int32): BlogPost | ApiError;
}

// ---------------------------------------------------------------------------
// Reviews (legacy `/admin/reviews`) — full CRUD.
// ---------------------------------------------------------------------------

/** The input of a review: never its Publication State (see BannerInput). */
model ReviewInput {
  pinned: boolean | null;
  courseId: int32 | null;
  userId: int32 | null;
  body: string | null;
  firstName: string | null;
  lastName: string | null;
}

@route("/admin/reviews")
@tag("Admin: Reviews")
@useAuth(AdminSession)
interface AdminReviews {
  @get @extension("x-permission", #{ name: "reviews.read", kind: "read" }) @extension("x-personal-data", "Review") @operationId("adminListReviews") list(...SearchListQuery): Page<Review> | ApiError;

  /** Get a single review. A missing id surfaces as 404 via the central
      ent-error handler, not a typed union member. */
  @get @route("/{id}") @extension("x-permission", #{ name: "reviews.read", kind: "read" }) @extension("x-personal-data", "Review") @operationId("adminGetReview") get(@path id: int32): Review | ApiError;

  @post @useAuth(AdminSession) @extension("x-permission", #{ name: "reviews.create", kind: "write" }) @operationId("adminCreateReview") create(@body body: ReviewInput): {
    @statusCode statusCode: 201;
    @body review: Review;
  } | ApiError;

  /** Update a review. A missing id surfaces as 404 via the central ent-error
      handler. */
  @put @useAuth(AdminSession) @route("/{id}") @extension("x-permission", #{ name: "reviews.update", kind: "write" }) @operationId("adminUpdateReview") update(@path id: int32, @body body: ReviewInput): Review | ApiError;

  @delete @useAuth(AdminSession) @route("/{id}") @extension("x-permission", #{ name: "reviews.delete", kind: "write" }) @operationId("adminDeleteReview") delete(@path id: int32): NoContent | ApiError;

  /** Publish a draft or archived review. A state is not an update: each change of
      Publication State is its own operation with its own Permission. An illegal
      move, or repeating the current state, is a 409 with a localized `detail`. */
  @post @useAuth(AdminSession) @route("/{id}/publish") @extension("x-permission", #{ name: "reviews.publish", kind: "write" }) @operationId("adminPublishReview")
  publish(@path id: int32): Review | ApiError;

  /** Take a published review back to draft. */
  @post @useAuth(AdminSession) @route("/{id}/unpublish") @extension("x-permission", #{ name: "reviews.unpublish", kind: "write" }) @operationId("adminUnpublishReview")
  unpublish(@path id: int32): Review | ApiError;

  /** Archive a draft or published review. */
  @post @useAuth(AdminSession) @route("/{id}/archive") @extension("x-permission", #{ name: "reviews.archive", kind: "write" }) @operationId("adminArchiveReview")
  archive(@path id: int32): Review | ApiError;
}

// ---------------------------------------------------------------------------
// Courses — index/get/create/update (no delete)
// + trigger review + create version.
// ---------------------------------------------------------------------------

/** The input of a course: never its Readiness. A new course is a draft and the
    Readiness changes only through the mark-as operations. */
model CourseInput {
  slug: string | null;
  learnAs: CourseLearnAs | null;
  hexletProgramLandingPage: string | null;
  repositoryUrl: string | null;

  /** Attachment id from `POST /admin/attachments`, or null. */
  coverAttachmentId: int32 | null;
}

/** A course reduced to what a select needs. */
model CourseOption {
  id: int32;

  /** The course's name, or its slug when it has none. */
  name: string;
}

@route("/admin/courses")
@tag("Admin: Courses")
@useAuth(AdminSession)
interface AdminCourses {
  /** Every course as id and name, unpaginated and oldest first, for selects.
      It is the post form's course picker, so it is classified under blog posts:
      an editor of posts picks the courses a post promotes without being allowed
      to edit courses. */
  @get @route("/options") @extension("x-permission", #{ name: "blog_posts.read", kind: "read" }) @operationId("adminListCourseOptions") listOptions(): CourseOption[] | ApiError;

  @get @extension("x-permission", #{ name: "courses.read", kind: "read" }) @operationId("adminListCourses") list(...SearchListQuery): Page<Course> | ApiError;

  /** Get a single course. A missing id surfaces as 404 via the central
      ent-error handler. */
  @get @route("/{id}") @extension("x-permission", #{ name: "courses.read", kind: "read" }) @operationId("adminGetCourse") get(@path id: int32): Course | ApiError;

  @post @useAuth(AdminSession) @extension("x-permission", #{ name: "courses.create", kind: "write" }) @operationId("adminCreateCourse") create(@body body: CourseInput): {
    @statusCode statusCode: 201;
    @body course: Course;
  } | ApiError;

  @put @useAuth(AdminSession) @route("/{id}") @extension("x-permission", #{ name: "courses.update", kind: "write" }) @operationId("adminUpdateCourse") update(@path id: int32, @body body: CourseInput): Course | ApiError;

  /** Mark a course as a draft. A state is not an update: each change of Course
      Readiness is its own operation with its own Permission. Any change between two
      distinct values is legal; repeating the current value is a 409 with a localized
      `detail`. */
  @post @useAuth(AdminSession) @route("/{id}/mark_draft") @extension("x-permission", #{ name: "courses.mark_draft", kind: "write" }) @operationId("adminMarkCourseDraft")
  markDraft(@path id: int32): Course | ApiError;

  /** Mark a course as in development. */
  @post @useAuth(AdminSession) @route("/{id}/mark_in_development") @extension("x-permission", #{ name: "courses.mark_in_development", kind: "write" }) @operationId("adminMarkCourseInDevelopment")
  markInDevelopment(@path id: int32): Course | ApiError;

  /** Mark a course as completed: it joins the public catalogue. */
  @post @useAuth(AdminSession) @route("/{id}/mark_completed") @extension("x-permission", #{ name: "courses.mark_completed", kind: "write" }) @operationId("adminMarkCourseCompleted")
  markCompleted(@path id: int32): Course | ApiError;

  /** Enqueue AI re-review of every current lesson version. */
  @post @useAuth(AdminSession) @route("/{id}/review") @extension("x-permission", #{ name: "courses.review", kind: "write" }) @operationId("adminReviewCourse") review(@path id: int32): NoContent | ApiError;

  /** Build a new version of the course exercises. */
  @post @useAuth(AdminSession) @route("/{id}/versions") @extension("x-permission", #{ name: "course_versions.create", kind: "write" }) @operationId("adminCreateCourseVersion")
  createVersion(@path id: int32): {
    @statusCode statusCode: 201;
    @body version: CourseVersion;
  } | ApiError;
}

// ---------------------------------------------------------------------------
// Course landing pages — full CRUD
// + nested QnA items.
// ---------------------------------------------------------------------------

/** The input of a landing page: never its Publication State (see BannerInput). */
model CourseLandingPageInput {
  courseId: int32;
  slug: string | null;
  name: string | null;
  main: boolean | null;
  listed: boolean | null;
  footer: boolean | null;
  footerName: string | null;
  order: string | null;
  landingPageToRedirectId: int32 | null;
  metaTitle: string | null;
  metaDescription: string | null;
  header: string | null;
  description: string | null;
  usedInHeader: string | null;
  usedInDescription: string | null;
  outcomesHeader: string | null;
  outcomesDescription: string | null;

  /** Attachment id from `POST /admin/attachments`, or null. */
  outcomesImageAttachmentId: int32 | null;
}

@route("/admin/course_landing_pages")
@tag("Admin: Course Landing Pages")
@useAuth(AdminSession)
interface AdminCourseLandingPages {
  @get @extension("x-permission", #{ name: "course_landing_pages.read", kind: "read" }) @operationId("adminListCourseLandingPages") list(...SearchListQuery): Page<CourseLandingPage> | ApiError;

  /** Get a single landing page. A missing id surfaces as 404 via the central
      ent-error handler, not a typed union member. */
  @get @route("/{id}") @extension("x-permission", #{ name: "course_landing_pages.read", kind: "read" }) @operationId("adminGetCourseLandingPage") get(@path id: int32): CourseLandingPage | ApiError;

  @post @useAuth(AdminSession) @extension("x-permission", #{ name: "course_landing_pages.create", kind: "write" }) @operationId("adminCreateCourseLandingPage") create(@body body: CourseLandingPageInput): {
    @statusCode statusCode: 201;
    @body landingPage: CourseLandingPage;
  } | ApiError;

  /** Update a landing page. A missing id surfaces as 404 via the central
      ent-error handler. */
  @put @useAuth(AdminSession) @route("/{id}") @extension("x-permission", #{ name: "course_landing_pages.update", kind: "write" }) @operationId("adminUpdateCourseLandingPage") update(@path id: int32, @body body: CourseLandingPageInput): CourseLandingPage | ApiError;

  @delete @useAuth(AdminSession) @route("/{id}") @extension("x-permission", #{ name: "course_landing_pages.delete", kind: "write" }) @operationId("adminDeleteCourseLandingPage") delete(@path id: int32): NoContent | ApiError;

  /** Publish a draft or archived landing page. A state is not an update: each change of
      Publication State is its own operation with its own Permission. An illegal
      move, or repeating the current state, is a 409 with a localized `detail`. */
  @post @useAuth(AdminSession) @route("/{id}/publish") @extension("x-permission", #{ name: "course_landing_pages.publish", kind: "write" }) @operationId("adminPublishCourseLandingPage")
  publish(@path id: int32): CourseLandingPage | ApiError;

  /** Take a published landing page back to draft. */
  @post @useAuth(AdminSession) @route("/{id}/unpublish") @extension("x-permission", #{ name: "course_landing_pages.unpublish", kind: "write" }) @operationId("adminUnpublishCourseLandingPage")
  unpublish(@path id: int32): CourseLandingPage | ApiError;

  /** Archive a draft or published landing page. */
  @post @useAuth(AdminSession) @route("/{id}/archive") @extension("x-permission", #{ name: "course_landing_pages.archive", kind: "write" }) @operationId("adminArchiveCourseLandingPage")
  archive(@path id: int32): CourseLandingPage | ApiError;
}

// ---- Nested QnA items ------------------------------------------------------

model QnaItemInput {
  @minLength(1) question: string;
  @minLength(1) answer: string;
}

@route("/admin/course_landing_pages/{landingPageId}/qna_items")
@tag("Admin: Course Landing Pages")
@useAuth(AdminSession)
interface AdminCourseLandingPageQnaItems {
  @get @extension("x-permission", #{ name: "course_landing_page_qna_items.read", kind: "read" }) @operationId("adminListLandingPageQnaItems") list(@path landingPageId: int32): QnaItem[] | ApiError;

  @post @useAuth(AdminSession) @extension("x-permission", #{ name: "course_landing_page_qna_items.create", kind: "write" }) @operationId("adminCreateLandingPageQnaItem") create(@path landingPageId: int32, @body body: QnaItemInput): {
    @statusCode statusCode: 201;
    @body item: QnaItem;
  } | ApiError;

  /** Update a QnA item. A missing id (or one under a different parent) surfaces
      as 404 via the central ent-error handler. */
  @put @useAuth(AdminSession) @route("/{id}") @extension("x-permission", #{ name: "course_landing_page_qna_items.update", kind: "write" }) @operationId("adminUpdateLandingPageQnaItem") update(@path landingPageId: int32, @path id: int32, @body body: QnaItemInput): QnaItem | ApiError;

  @delete @useAuth(AdminSession) @route("/{id}") @extension("x-permission", #{ name: "course_landing_page_qna_items.delete", kind: "write" }) @operationId("adminDeleteLandingPageQnaItem") delete(@path landingPageId: int32, @path id: int32): NoContent | ApiError;
}

@route("/admin/course_categories/{categoryId}/qna_items")
@tag("Admin: Course Categories")
@useAuth(AdminSession)
interface AdminCourseCategoryQnaItems {
  @get @extension("x-permission", #{ name: "course_category_qna_items.read", kind: "read" }) @operationId("adminListCategoryQnaItems") list(@path categoryId: int32): QnaItem[] | ApiError;

  @post @useAuth(AdminSession) @extension("x-permission", #{ name: "course_category_qna_items.create", kind: "write" }) @operationId("adminCreateCategoryQnaItem") create(@path categoryId: int32, @body body: QnaItemInput): {
    @statusCode statusCode: 201;
    @body item: QnaItem;
  } | ApiError;

  /** Update a QnA item. A missing id (or one under a different parent) surfaces
      as 404 via the central ent-error handler. */
  @put @useAuth(AdminSession) @route("/{id}") @extension("x-permission", #{ name: "course_category_qna_items.update", kind: "write" }) @operationId("adminUpdateCategoryQnaItem") update(@path categoryId: int32, @path id: int32, @body body: QnaItemInput): QnaItem | ApiError;

  @delete @useAuth(AdminSession) @route("/{id}") @extension("x-permission", #{ name: "course_category_qna_items.delete", kind: "write" }) @operationId("adminDeleteCategoryQnaItem") delete(@path categoryId: int32, @path id: int32): NoContent | ApiError;
}

// ---------------------------------------------------------------------------
// Read-only admin lists (legacy index-only routes).
// ---------------------------------------------------------------------------

@route("/admin/course_lessons")
@tag("Admin: Course Lessons")
@useAuth(AdminSession)
interface AdminCourseLessons {
  @get @extension("x-permission", #{ name: "course_lessons.read", kind: "read" }) @operationId("adminListCourseLessons") list(...SearchListQuery): Page<CourseLessonListItem> | ApiError;

  /** Enqueue AI review for a single lesson. */
  @post @useAuth(AdminSession) @route("/{id}/review") @extension("x-permission", #{ name: "course_lessons.review", kind: "write" }) @operationId("adminReviewCourseLesson") review(@path id: int32): NoContent | ApiError;
}

@route("/admin/lesson_progress")
@tag("Admin: Lesson Progress")
@useAuth(AdminSession)
interface AdminLessonProgress {
  @get @extension("x-permission", #{ name: "lesson_progress.read", kind: "read" }) @extension("x-personal-data", "LessonProgress") @operationId("adminListLessonProgress") list(...ListQuery): Page<LessonProgress> | ApiError;
}

@route("/admin/course_lesson_reviews")
@tag("Admin: Course Lesson Reviews")
@useAuth(AdminSession)
interface AdminCourseLessonReviews {
  @get @extension("x-permission", #{ name: "course_lesson_reviews.read", kind: "read" }) @operationId("adminListCourseLessonReviews") list(...ListQuery): Page<CourseLessonReview> | ApiError;
}

@route("/admin/leads")
@tag("Admin: Leads")
@useAuth(AdminSession)
interface AdminLeads {
  @get @extension("x-permission", #{ name: "leads.read", kind: "read" }) @extension("x-personal-data", "Lead") @operationId("adminListLeads") list(...SearchListQuery): Page<Lead> | ApiError;
}

@route("/admin/messages")
@tag("Admin: Messages")
@useAuth(AdminSession)
interface AdminMessages {
  @get @extension("x-permission", #{ name: "messages.read", kind: "read" }) @extension("x-personal-data", "AiMessage") @operationId("adminListMessages") list(...SearchListQuery): Page<LessonAssistantMessage> | ApiError;
}

/**
 * The capability Permissions (`audit.read`, `mcp.read`, ...) the signed-in staff
 * member holds through their Staff Role. The administrator flag adds none, so
 * the admin menu shows an entry behind a capability only to someone who holds it.
 */
model MyCapabilities {
  held: string[];
}

@route("/admin/capabilities")
@tag("Admin: Audit")
@useAuth(AdminSession)
interface AdminCapabilities {
  /** The capabilities the caller holds; open to any staff member. */
  @get @extension("x-access", "any-staff") @operationId("adminGetMyCapabilities")
  get(): MyCapabilities | ApiError;
}

/**
 * One record the admin catalogue search found: a `SearchResult` that also says
 * how published the record is, because the admin search finds drafts too.
 */
model AdminSearchResult {
  type: SearchResultType;

  /** The typed id `adminFetch` takes: `course:LANDING_SLUG`, `lesson:COURSE_SLUG/LESSON_SLUG` or `post:SLUG`. */
  id: string;

  title: string;

  /** A few words of the record's text around the match, plain text. */
  snippet: string;

  /** The record's page on the website once it is public: the locale prefix and the path, with no origin. */
  path: string;

  /** Publication state of the landing page, lesson or post (`draft`, `published`, `archived`, ...). */
  state: string;

  /** Readiness of the course (`draft`, `in_development`, `completed`); empty for a post. */
  readiness: string;
}

/** What is known about how published a record is. */
model AdminRecordMetadata {
  type: SearchResultType;
  slug: string;

  /** The course a lesson belongs to, or a course page's own course slug. */
  courseSlug?: string;

  /** Publication state of the landing page, lesson or post. */
  state: string;

  /** Readiness of the course; empty for a post. */
  readiness: string;

  /** Whether a landing page is on the course list. */
  listed?: boolean;

  /** State of the build a lesson was read from. */
  buildState?: string;

  /** Natural order of a lesson in its course. */
  position?: int32;
}

/** One record as an editor reads it, whatever state it is in. */
model AdminRecord {
  id: string;
  title: string;

  /** The record's text, as plain text or Markdown. */
  text: string;

  /** The record's page on the website once it is public, with the locale prefix and no origin. */
  path: string;

  metadata: AdminRecordMetadata;
}

@route("/admin/search")
@tag("Admin: Search")
@useAuth(AdminSession)
interface AdminCatalogueSearch {
  /**
   * Search every course, lesson and blog post by what they say, in the request
   * locale, whatever state they are in: draft, archived and unlisted landing
   * pages, courses not completed yet, lessons of the newest build and
   * unpublished posts, each with its state. Open to any staff member (and over
   * MCP to one holding mcp.read), who may read the content; the locale must be
   * one their Staff Role allows. The public `searchCatalogue` keeps its filters.
   */
  @get @extension("x-access", "any-staff") @operationId("adminSearch")
  search(
    @query @minLength(1) @maxLength(200) q: string,

    /** Only this kind of record; every kind when absent. */
    @query type?: SearchResultType,

    ...ListQuery,
  ): Page<AdminSearchResult> | ApiError;

  /** Read one record an `adminSearch` result names, in the request locale, whatever state it is in. */
  @get @route("/record") @extension("x-access", "any-staff") @operationId("adminFetch")
  fetch(@query @minLength(1) id: string): AdminRecord | ApiError;
}

/**
 * Filters of the audit trail listing. The table is partitioned by month on
 * `occurred_at`, so every query carries a time window: `from` defaults to 30
 * days before `to`, `to` defaults to now, and the window may not exceed 92 days
 * (a wider one is a 422). A narrower window is the way to look further back.
 */
model AuditRecordQuery {
  @query @minValue(1) page?: int32;
  @query @minValue(1) @maxValue(100) perPage?: int32;

  @query actorUserId?: int32;
  @query channel?: "web" | "mcp";
  @query operation?: string;
  @query entityType?: string;
  @query entityId?: string;

  /** Start of the window (inclusive). */
  @query from?: utcDateTime;

  /** End of the window (exclusive). */
  @query to?: utcDateTime;
}

@route("/admin/audit_records")
@tag("Admin: Audit")
@useAuth(AdminSession)
interface AdminAuditRecords {
  /**
   * The audit trail, newest first. Needs the capability Permission `audit.read`
   * on the Staff Role even for an administrator, and is itself recorded as a
   * read of personal data: records carry addresses and old values.
   */
  @get @extension("x-permission", #{ name: "audit.read", kind: "read" }) @extension("x-personal-data", "AuditRecord") @operationId("adminListAuditRecords")
  list(...AuditRecordQuery): Page<AuditRecord> | ApiError;
}

// ---------------------------------------------------------------------------
// Admin users (legacy `/admin/api/users`) — full CRUD + search.
// ---------------------------------------------------------------------------

model UserInput {
  @format("email") email: string;
  firstName: string | null;
  lastName: string | null;
  admin: boolean | null;
}

@route("/admin/api/users")
@tag("Admin: Users")
@useAuth(AdminSession)
interface AdminUsers {
  @get @extension("x-permission", #{ name: "users.read", kind: "read" }) @extension("x-personal-data", "User") @extension("x-requires-capability", "mcp.manage") @operationId("adminListUsers") list(...SearchListQuery): Page<UserCrud> | ApiError;

  /** Typeahead search by name/email. */
  @get @route("/search") @extension("x-permission", #{ name: "users.read", kind: "read" }) @extension("x-personal-data", "User") @extension("x-requires-capability", "mcp.manage") @operationId("adminSearchUsers") search(@query q: string): UserCrud[] | ApiError;

  /** Get a single user. A missing id surfaces as 404 via the central ent-error
      handler, not a typed union member. */
  @get @route("/{id}") @extension("x-permission", #{ name: "users.read", kind: "read" }) @extension("x-personal-data", "User") @extension("x-requires-capability", "mcp.manage") @operationId("adminGetUser") get(@path id: int32): UserCrud | ApiError;

  /** Create a user. A duplicate email is a DB unique constraint, surfaced as 409
      by the central ent-error handler. */
  @post @useAuth(AdminSession) @extension("x-permission", #{ name: "users.create", kind: "write" }) @extension("x-requires-capability", "mcp.manage") @operationId("adminCreateUser") create(@body body: UserInput): {
    @statusCode statusCode: 201;
    @body user: UserCrud;
  } | ApiError;

  /** Update a user. 404 (missing) and 409 (duplicate email) both flow through
      the central ent-error handler. */
  @put @useAuth(AdminSession) @route("/{id}") @extension("x-permission", #{ name: "users.update", kind: "write" }) @extension("x-requires-capability", "mcp.manage") @operationId("adminUpdateUser") update(@path id: int32, @body body: UserInput): UserCrud | ApiError;

  @delete @useAuth(AdminSession) @route("/{id}") @extension("x-permission", #{ name: "users.delete", kind: "write" }) @extension("x-requires-capability", "mcp.manage") @operationId("adminDeleteUser") delete(@path id: int32): NoContent | ApiError;
}

// ---------------------------------------------------------------------------
// Management (legacy `/admin/management/*`).
// ---------------------------------------------------------------------------

model RoleInput {
  @minLength(1) name: string;
  description: string | null;
}

@route("/admin/management/roles")
@tag("Admin: Roles")
@useAuth(AdminSession)
interface AdminRoles {
  @get @extension("x-permission", #{ name: "staff_roles.read", kind: "read" }) @extension("x-requires-capability", "mcp.manage") @operationId("adminListRoles") list(...ListQuery): Page<StaffRole> | ApiError;

  /** Get a role with its permission matrix. A missing id surfaces as 404 via
      the central ent-error handler. */
  @get @route("/{id}") @extension("x-permission", #{ name: "staff_roles.read", kind: "read" }) @extension("x-requires-capability", "mcp.manage") @operationId("adminGetRole") get(@path id: int32): StaffRoleDetail | ApiError;

  /** Create a role. A duplicate name is a DB unique constraint, surfaced as 409
      by the central ent-error handler. */
  @post @useAuth(AdminSession) @extension("x-permission", #{ name: "staff_roles.create", kind: "write" }) @extension("x-requires-capability", "mcp.manage") @operationId("adminCreateRole") create(@body body: RoleInput): {
    @statusCode statusCode: 201;
    @body role: StaffRoleDetail;
  } | ApiError;

  @put @useAuth(AdminSession) @route("/{id}") @extension("x-permission", #{ name: "staff_roles.update", kind: "write" }) @extension("x-requires-capability", "mcp.manage") @operationId("adminUpdateRole") update(@path id: int32, @body body: RoleInput): StaffRoleDetail | ApiError;

  @delete @useAuth(AdminSession) @route("/{id}") @extension("x-permission", #{ name: "staff_roles.delete", kind: "write" }) @extension("x-requires-capability", "mcp.manage") @operationId("adminDeleteRole") delete(@path id: int32): NoContent | ApiError;
}

/** The Permissions submitted for a role. */
model RolePermissionsInput {
  /**
   * The declared Permissions the role holds, replacing its current set. Names
   * must come from the catalog (`GET /admin/management/permissions`).
   */
  grants: string[];
}

/** Every Permission an admin operation can require. */
model PermissionCatalog {
  permissions: Permission[];
}

@route("/admin/management/permissions")
@tag("Admin: Roles")
@useAuth(AdminSession)
interface AdminPermissions {
  /** The catalog of Permissions, generated from the operations' declarations. */
  @get @extension("x-permission", #{ name: "staff_roles.read", kind: "read" }) @extension("x-requires-capability", "mcp.manage") @operationId("adminListPermissions") list(): PermissionCatalog | ApiError;
}

@route("/admin/management/role_permissions")
@tag("Admin: Roles")
@useAuth(AdminSession)
interface AdminRolePermissions {
  /** The permission matrix for a role. */
  @get @route("/{roleId}") @extension("x-permission", #{ name: "staff_roles.read", kind: "read" }) @extension("x-requires-capability", "mcp.manage") @operationId("adminGetRolePermissions") get(@path roleId: int32): StaffRoleDetail | ApiError;

  /** Replace the permission matrix for a role. */
  @put @useAuth(AdminSession) @route("/{roleId}") @extension("x-permission", #{ name: "staff_roles.update", kind: "write" }) @extension("x-requires-capability", "mcp.manage") @operationId("adminUpdateRolePermissions") update(@path roleId: int32, @body body: RolePermissionsInput): StaffRoleDetail | ApiError;
}

model StaffMemberInput {
  userId: int32;
  roleId: int32;
  allowedLocales: string[];
}

@route("/admin/management/staff_members")
@tag("Admin: Staff Members")
@useAuth(AdminSession)
interface AdminStaffMembers {
  @get @extension("x-permission", #{ name: "staff_members.read", kind: "read" }) @extension("x-personal-data", "StaffMember") @extension("x-requires-capability", "mcp.manage") @operationId("adminListStaffMembers") list(...ListQuery): Page<StaffMember> | ApiError;

  @get @route("/{id}") @extension("x-permission", #{ name: "staff_members.read", kind: "read" }) @extension("x-personal-data", "StaffMember") @extension("x-requires-capability", "mcp.manage") @operationId("adminGetStaffMember") get(@path id: int32): StaffMember | ApiError;

  @post @useAuth(AdminSession) @extension("x-permission", #{ name: "staff_members.create", kind: "write" }) @extension("x-requires-capability", "mcp.manage") @operationId("adminCreateStaffMember") create(@body body: StaffMemberInput): {
    @statusCode statusCode: 201;
    @body staffMember: StaffMember;
  } | ApiError;

  @put @useAuth(AdminSession) @route("/{id}") @extension("x-permission", #{ name: "staff_members.update", kind: "write" }) @extension("x-requires-capability", "mcp.manage") @operationId("adminUpdateStaffMember") update(@path id: int32, @body body: StaffMemberInput): StaffMember | ApiError;

  @delete @useAuth(AdminSession) @route("/{id}") @extension("x-permission", #{ name: "staff_members.delete", kind: "write" }) @extension("x-requires-capability", "mcp.manage") @operationId("adminDeleteStaffMember") delete(@path id: int32): NoContent | ApiError;
}

@route("/admin/management/users")
@tag("Admin: Users")
@useAuth(AdminSession)
interface AdminManagementUsers {
  @get @extension("x-permission", #{ name: "users.read", kind: "read" }) @extension("x-personal-data", "User") @extension("x-requires-capability", "mcp.manage") @operationId("adminListManagementUsers") list(...SearchListQuery): Page<UserCrud> | ApiError;

  /** Get a management user. A missing id surfaces as 404 via the central
      ent-error handler. */
  @get @route("/{id}") @extension("x-permission", #{ name: "users.read", kind: "read" }) @extension("x-personal-data", "User") @extension("x-requires-capability", "mcp.manage") @operationId("adminGetManagementUser") get(@path id: int32): UserCrud | ApiError;

  @put @useAuth(AdminSession) @route("/{id}") @extension("x-permission", #{ name: "users.update", kind: "write" }) @extension("x-requires-capability", "mcp.manage") @operationId("adminUpdateManagementUser") update(@path id: int32, @body body: UserInput): UserCrud | ApiError;
}
