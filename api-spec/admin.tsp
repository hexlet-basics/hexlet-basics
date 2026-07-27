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
// one engine (<CrudList> + <CrudForm>) drives them all. URLs keep their legacy
// `/admin/language_*` paths (ADR-0002); every model uses the domain name.
//
// Uniform response shape: writes return `T | ValidationError` (422). Index is
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
  file: HttpPart<File>;
}

@route("/admin/attachments")
@tag("Admin: Attachments")
interface AdminAttachments {
  /** Upload a file; returns the attachment to reference from a form. */
  @post
  @operationId("adminUploadAttachment")
  upload(@header contentType: "multipart/form-data", @multipartBody body: AttachmentUploadForm): {
    @statusCode statusCode: 201;
    @body attachment: Attachment;
  } | ValidationError;
}

// ---------------------------------------------------------------------------
// Course Categories (legacy `/admin/language_categories`).
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

@route("/admin/language_categories")
@tag("Admin: Course Categories")
interface AdminCourseCategories {
  /** List course categories (paginated). */
  @get
  @operationId("adminListCourseCategories")
  list(...ListQuery): Page<CourseCategory>;

  /** Get a single course category. A missing id surfaces as 404 via the
      central ent-error handler, not a typed union member. */
  @get
  @route("/{id}")
  @operationId("adminGetCourseCategory")
  get(@path id: int32): CourseCategory;

  /** Create a course category. A uniqueness violation (name/header/slug) is a
      DB constraint, surfaced as 409 by the central ent-error handler. */
  @post
  @operationId("adminCreateCourseCategory")
  create(@body body: CourseCategoryInput): {
    @statusCode statusCode: 201;
    @body category: CourseCategory;
  };

  /** Update a course category. 404 (missing) and 409 (uniqueness) both flow
      through the central ent-error handler. */
  @put
  @route("/{id}")
  @operationId("adminUpdateCourseCategory")
  update(@path id: int32, @body body: CourseCategoryInput): CourseCategory;

  /** Delete a course category. */
  @delete
  @route("/{id}")
  @operationId("adminDeleteCourseCategory")
  delete(@path id: int32): NoContent;
}

// ---------------------------------------------------------------------------
// Banners (legacy `/admin/banners`) — full CRUD.
// ---------------------------------------------------------------------------

model BannerInput {
  state: BannerState;
  background: BannerBackground;
  locale: BannerLocale;
  @minLength(1) body: string;
  url: string | null;
  startsAt: utcDateTime | null;
  finishesAt: utcDateTime | null;
}

@route("/admin/banners")
@tag("Admin: Banners")
interface AdminBanners {
  /** List banners (paginated). */
  @get @operationId("adminListBanners") list(...ListQuery): Page<Banner>;

  /** Get a single banner. A missing id surfaces as 404 via the central
      ent-error handler, not a typed union member. */
  @get @route("/{id}") @operationId("adminGetBanner") get(@path id: int32): Banner;

  /** Create a banner. A body violating the schema (e.g. an empty `body`)
      surfaces as 400 via the central handler; banners carry no uniqueness
      constraint, so there is no 409 path. */
  @post @operationId("adminCreateBanner") create(@body body: BannerInput): {
    @statusCode statusCode: 201;
    @body banner: Banner;
  };

  /** Update a banner. A missing id surfaces as 404 via the central handler. */
  @put @route("/{id}") @operationId("adminUpdateBanner") update(@path id: int32, @body body: BannerInput): Banner;

  /** Delete a banner. */
  @delete @route("/{id}") @operationId("adminDeleteBanner") delete(@path id: int32): NoContent;
}

// ---------------------------------------------------------------------------
// Blog posts (legacy `/admin/blog_posts`) — full CRUD + related_courses.
// ---------------------------------------------------------------------------

model BlogPostInput {
  name: string | null;
  slug: string | null;
  description: string | null;
  state: BlogPostState | null;

  /** Markdown/rich source; `richBodyHtml` is derived on read. */
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
interface AdminBlogPosts {
  @get @operationId("adminListBlogPosts") list(...ListQuery): Page<BlogPost>;

  // 404 comes from the central APIErrorHandler (ent not-found), like every other
  // admin get — so no explicit NotFoundError union here.
  @get @route("/{id}") @operationId("adminGetBlogPost") get(@path id: int32): BlogPost;

  @post @operationId("adminCreateBlogPost") create(@body body: BlogPostInput): {
    @statusCode statusCode: 201;
    @body post: BlogPost;
  } | ValidationError;

  @put @route("/{id}") @operationId("adminUpdateBlogPost") update(@path id: int32, @body body: BlogPostInput): BlogPost | ValidationError;

  @delete @route("/{id}") @operationId("adminDeleteBlogPost") delete(@path id: int32): NoContent;

  /** Set the related/promoted courses for a post. */
  @post @route("/{id}/related_courses") @operationId("adminSetBlogPostRelatedCourses")
  relatedCourses(@path id: int32, @body body: BlogPostRelatedCoursesInput): BlogPost | ValidationError;
}

// ---------------------------------------------------------------------------
// Reviews (legacy `/admin/reviews`) — full CRUD.
// ---------------------------------------------------------------------------

model ReviewInput {
  state: ReviewState | null;
  pinned: boolean | null;
  courseId: int32 | null;
  userId: int32 | null;
  body: string | null;
  firstName: string | null;
  lastName: string | null;
}

@route("/admin/reviews")
@tag("Admin: Reviews")
interface AdminReviews {
  @get @operationId("adminListReviews") list(...ListQuery): Page<Review>;

  /** Get a single review. A missing id surfaces as 404 via the central
      ent-error handler, not a typed union member. */
  @get @route("/{id}") @operationId("adminGetReview") get(@path id: int32): Review;

  @post @operationId("adminCreateReview") create(@body body: ReviewInput): {
    @statusCode statusCode: 201;
    @body review: Review;
  };

  /** Update a review. A missing id surfaces as 404 via the central ent-error
      handler. */
  @put @route("/{id}") @operationId("adminUpdateReview") update(@path id: int32, @body body: ReviewInput): Review;

  @delete @route("/{id}") @operationId("adminDeleteReview") delete(@path id: int32): NoContent;
}

// ---------------------------------------------------------------------------
// Courses (legacy `/admin/languages`) — index/get/create/update (no delete)
// + trigger review + create version.
// ---------------------------------------------------------------------------

model CourseInput {
  slug: string | null;
  learnAs: CourseLearnAs | null;
  progress: CourseProgress | null;
  hexletProgramLandingPage: string | null;
  repositoryUrl: string | null;

  /** Attachment id from `POST /admin/attachments`, or null. */
  coverAttachmentId: int32 | null;
}

@route("/admin/languages")
@tag("Admin: Courses")
interface AdminCourses {
  @get @operationId("adminListCourses") list(...ListQuery): Page<Course>;

  /** Get a single course. A missing id surfaces as 404 via the central
      ent-error handler. */
  @get @route("/{id}") @operationId("adminGetCourse") get(@path id: int32): Course;

  @post @operationId("adminCreateCourse") create(@body body: CourseInput): {
    @statusCode statusCode: 201;
    @body course: Course;
  };

  @put @route("/{id}") @operationId("adminUpdateCourse") update(@path id: int32, @body body: CourseInput): Course;

  /** Enqueue AI re-review of every current lesson version. */
  @post @route("/{id}/review") @operationId("adminReviewCourse") review(@path id: int32): NoContent | NotFoundError;

  /** Build a new version of the course exercises. */
  @post @route("/{id}/versions") @operationId("adminCreateCourseVersion")
  createVersion(@path id: int32): {
    @statusCode statusCode: 201;
    @body version: CourseVersion;
  } | NotFoundError;
}

// ---------------------------------------------------------------------------
// Course landing pages (legacy `/admin/language_landing_pages`) — full CRUD
// + nested QnA items.
// ---------------------------------------------------------------------------

model CourseLandingPageInput {
  courseId: int32;
  slug: string | null;
  name: string | null;
  main: boolean | null;
  listed: boolean | null;
  footer: boolean | null;
  footerName: string | null;
  state: LandingPageState | null;
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

@route("/admin/language_landing_pages")
@tag("Admin: Course Landing Pages")
interface AdminCourseLandingPages {
  @get @operationId("adminListCourseLandingPages") list(...ListQuery): Page<CourseLandingPage>;

  /** Get a single landing page. A missing id surfaces as 404 via the central
      ent-error handler, not a typed union member. */
  @get @route("/{id}") @operationId("adminGetCourseLandingPage") get(@path id: int32): CourseLandingPage;

  @post @operationId("adminCreateCourseLandingPage") create(@body body: CourseLandingPageInput): {
    @statusCode statusCode: 201;
    @body landingPage: CourseLandingPage;
  };

  /** Update a landing page. A missing id surfaces as 404 via the central
      ent-error handler. */
  @put @route("/{id}") @operationId("adminUpdateCourseLandingPage") update(@path id: int32, @body body: CourseLandingPageInput): CourseLandingPage;

  @delete @route("/{id}") @operationId("adminDeleteCourseLandingPage") delete(@path id: int32): NoContent;
}

// ---- Nested QnA items ------------------------------------------------------

model QnaItemInput {
  @minLength(1) question: string;
  @minLength(1) answer: string;
}

@route("/admin/language_landing_pages/{landingPageId}/qna_items")
@tag("Admin: Course Landing Pages")
interface AdminCourseLandingPageQnaItems {
  @get @operationId("adminListLandingPageQnaItems") list(@path landingPageId: int32): QnaItem[];

  @post @operationId("adminCreateLandingPageQnaItem") create(@path landingPageId: int32, @body body: QnaItemInput): {
    @statusCode statusCode: 201;
    @body item: QnaItem;
  };

  /** Update a QnA item. A missing id (or one under a different parent) surfaces
      as 404 via the central ent-error handler. */
  @put @route("/{id}") @operationId("adminUpdateLandingPageQnaItem") update(@path landingPageId: int32, @path id: int32, @body body: QnaItemInput): QnaItem;

  @delete @route("/{id}") @operationId("adminDeleteLandingPageQnaItem") delete(@path landingPageId: int32, @path id: int32): NoContent;
}

@route("/admin/language_categories/{categoryId}/qna_items")
@tag("Admin: Course Categories")
interface AdminCourseCategoryQnaItems {
  @get @operationId("adminListCategoryQnaItems") list(@path categoryId: int32): QnaItem[];

  @post @operationId("adminCreateCategoryQnaItem") create(@path categoryId: int32, @body body: QnaItemInput): {
    @statusCode statusCode: 201;
    @body item: QnaItem;
  };

  /** Update a QnA item. A missing id (or one under a different parent) surfaces
      as 404 via the central ent-error handler. */
  @put @route("/{id}") @operationId("adminUpdateCategoryQnaItem") update(@path categoryId: int32, @path id: int32, @body body: QnaItemInput): QnaItem;

  @delete @route("/{id}") @operationId("adminDeleteCategoryQnaItem") delete(@path categoryId: int32, @path id: int32): NoContent;
}

// ---------------------------------------------------------------------------
// Read-only admin lists (legacy index-only routes).
// ---------------------------------------------------------------------------

@route("/admin/language_lessons")
@tag("Admin: Course Lessons")
interface AdminCourseLessons {
  @get @operationId("adminListCourseLessons") list(...ListQuery): Page<CourseLessonListItem>;

  /** Enqueue AI review for a single lesson. */
  @post @route("/{id}/review") @operationId("adminReviewCourseLesson") review(@path id: int32): NoContent | NotFoundError;
}

@route("/admin/language_lesson_members")
@tag("Admin: Course Lesson Members")
interface AdminCourseLessonMembers {
  @get @operationId("adminListCourseLessonMembers") list(...ListQuery): Page<CourseLessonMember>;
}

@route("/admin/language_lesson_reviews")
@tag("Admin: Course Lesson Reviews")
interface AdminCourseLessonReviews {
  @get @operationId("adminListCourseLessonReviews") list(...ListQuery): Page<CourseLessonReview>;
}

@route("/admin/leads")
@tag("Admin: Leads")
interface AdminLeads {
  @get @operationId("adminListLeads") list(...ListQuery): Page<Lead>;
}

@route("/admin/messages")
@tag("Admin: Messages")
interface AdminMessages {
  @get @operationId("adminListMessages") list(...ListQuery): Page<LessonAssistantMessage>;
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
interface AdminUsers {
  @get @operationId("adminListUsers") list(...ListQuery): Page<UserCrud>;

  /** Typeahead search by name/email. */
  @get @route("/search") @operationId("adminSearchUsers") search(@query q: string): UserCrud[];

  /** Get a single user. A missing id surfaces as 404 via the central ent-error
      handler, not a typed union member. */
  @get @route("/{id}") @operationId("adminGetUser") get(@path id: int32): UserCrud;

  /** Create a user. A duplicate email is a DB unique constraint, surfaced as 409
      by the central ent-error handler. */
  @post @operationId("adminCreateUser") create(@body body: UserInput): {
    @statusCode statusCode: 201;
    @body user: UserCrud;
  };

  /** Update a user. 404 (missing) and 409 (duplicate email) both flow through
      the central ent-error handler. */
  @put @route("/{id}") @operationId("adminUpdateUser") update(@path id: int32, @body body: UserInput): UserCrud;

  @delete @route("/{id}") @operationId("adminDeleteUser") delete(@path id: int32): NoContent;
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
interface AdminRoles {
  @get @operationId("adminListRoles") list(...ListQuery): Page<StaffRole>;

  /** Get a role with its permission matrix. A missing id surfaces as 404 via
      the central ent-error handler. */
  @get @route("/{id}") @operationId("adminGetRole") get(@path id: int32): StaffRoleDetail;

  /** Create a role. A duplicate name is a DB unique constraint, surfaced as 409
      by the central ent-error handler. */
  @post @operationId("adminCreateRole") create(@body body: RoleInput): {
    @statusCode statusCode: 201;
    @body role: StaffRoleDetail;
  };

  @put @route("/{id}") @operationId("adminUpdateRole") update(@path id: int32, @body body: RoleInput): StaffRoleDetail;

  @delete @route("/{id}") @operationId("adminDeleteRole") delete(@path id: int32): NoContent;
}

/** One permission row to set on a role. */
model RolePermissionInput {
  resource: PermissionResource;
  canIndex: boolean;
  canCreate: boolean;
  canUpdate: boolean;
  canDestroy: boolean;
}

/** The full permission matrix submitted for a role. */
model RolePermissionsInput {
  permissions: RolePermissionInput[];
}

@route("/admin/management/role_permissions")
@tag("Admin: Roles")
interface AdminRolePermissions {
  /** The permission matrix for a role. */
  @get @route("/{roleId}") @operationId("adminGetRolePermissions") get(@path roleId: int32): StaffRoleDetail;

  /** Replace the permission matrix for a role. */
  @put @route("/{roleId}") @operationId("adminUpdateRolePermissions") update(@path roleId: int32, @body body: RolePermissionsInput): StaffRoleDetail;
}

model StaffMemberInput {
  userId: int32;
  roleId: int32;
  allowedLocales: string[];
}

@route("/admin/management/staff_members")
@tag("Admin: Staff Members")
interface AdminStaffMembers {
  @get @operationId("adminListStaffMembers") list(...ListQuery): Page<StaffMember>;

  @get @route("/{id}") @operationId("adminGetStaffMember") get(@path id: int32): StaffMember;

  @post @operationId("adminCreateStaffMember") create(@body body: StaffMemberInput): {
    @statusCode statusCode: 201;
    @body staffMember: StaffMember;
  };

  @put @route("/{id}") @operationId("adminUpdateStaffMember") update(@path id: int32, @body body: StaffMemberInput): StaffMember;

  @delete @route("/{id}") @operationId("adminDeleteStaffMember") delete(@path id: int32): NoContent;
}

@route("/admin/management/users")
@tag("Admin: Users")
interface AdminManagementUsers {
  @get @operationId("adminListManagementUsers") list(...ListQuery): Page<UserCrud>;

  /** Get a management user. A missing id surfaces as 404 via the central
      ent-error handler. */
  @get @route("/{id}") @operationId("adminGetManagementUser") get(@path id: int32): UserCrud;

  @put @route("/{id}") @operationId("adminUpdateManagementUser") update(@path id: int32, @body body: UserInput): UserCrud;
}
