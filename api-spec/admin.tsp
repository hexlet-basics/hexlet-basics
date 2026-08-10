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
@useAuth([AdminSession, XsrfToken])
interface AdminAttachments {
  /** Upload a file; returns the attachment to reference from a form. */
  @post
  @operationId("adminUploadAttachment")
  upload(@header contentType: "multipart/form-data", @multipartBody body: AttachmentUploadForm): {
    @statusCode statusCode: 201;
    @body attachment: Attachment;
  } | ValidationError | UnauthorizedError | ForbiddenError | ApiError;
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
  @operationId("adminListCourseCategories")
  list(...ListQuery): Page<CourseCategory> | UnauthorizedError | ForbiddenError | ApiError;

  /** Get a single course category. A missing id surfaces as 404 via the
      central ent-error handler, not a typed union member. */
  @get
  @route("/{id}")
  @operationId("adminGetCourseCategory")
  get(@path id: int32): CourseCategory | UnauthorizedError | ForbiddenError | ApiError;

  /** Create a course category. A uniqueness violation (name/header/slug) is a
      DB constraint, surfaced as 409 by the central ent-error handler. */
  @post
  @useAuth([AdminSession, XsrfToken])
  @operationId("adminCreateCourseCategory")
  create(@body body: CourseCategoryInput): {
    @statusCode statusCode: 201;
    @body category: CourseCategory;
  } | UnauthorizedError | ForbiddenError | ApiError;

  /** Update a course category. 404 (missing) and 409 (uniqueness) both flow
      through the central ent-error handler. */
  @put
  @useAuth([AdminSession, XsrfToken])
  @route("/{id}")
  @operationId("adminUpdateCourseCategory")
  update(@path id: int32, @body body: CourseCategoryInput): CourseCategory | UnauthorizedError | ForbiddenError | ApiError;

  /** Delete a course category. */
  @delete
  @useAuth([AdminSession, XsrfToken])
  @route("/{id}")
  @operationId("adminDeleteCourseCategory")
  delete(@path id: int32): NoContent | UnauthorizedError | ForbiddenError | ApiError;
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
@useAuth(AdminSession)
interface AdminBanners {
  /** List banners (paginated). */
  @get @operationId("adminListBanners") list(...ListQuery): Page<Banner> | UnauthorizedError | ForbiddenError | ApiError;

  /** Get a single banner. A missing id surfaces as 404 via the central
      ent-error handler, not a typed union member. */
  @get @route("/{id}") @operationId("adminGetBanner") get(@path id: int32): Banner | UnauthorizedError | ForbiddenError | ApiError;

  /** Create a banner. A body violating the schema (e.g. an empty `body`)
      surfaces as 400 via the central handler; banners carry no uniqueness
      constraint, so there is no 409 path. */
  @post @useAuth([AdminSession, XsrfToken]) @operationId("adminCreateBanner") create(@body body: BannerInput): {
    @statusCode statusCode: 201;
    @body banner: Banner;
  } | UnauthorizedError | ForbiddenError | ApiError;

  /** Update a banner. A missing id surfaces as 404 via the central handler. */
  @put @useAuth([AdminSession, XsrfToken]) @route("/{id}") @operationId("adminUpdateBanner") update(@path id: int32, @body body: BannerInput): Banner | UnauthorizedError | ForbiddenError | ApiError;

  /** Delete a banner. */
  @delete @useAuth([AdminSession, XsrfToken]) @route("/{id}") @operationId("adminDeleteBanner") delete(@path id: int32): NoContent | UnauthorizedError | ForbiddenError | ApiError;
}

// ---------------------------------------------------------------------------
// Blog posts (legacy `/admin/blog_posts`) — full CRUD + related_courses.
// ---------------------------------------------------------------------------

model BlogPostInput {
  name: string | null;
  slug: string | null;
  description: string | null;
  state: BlogPostState | null;

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
  @get @operationId("adminListBlogPosts") list(...ListQuery): Page<BlogPost> | UnauthorizedError | ForbiddenError | ApiError;

  // 404 comes from the shared ApiError response (ent not-found), like every
  // other admin get — so no operation-specific NotFoundError member is needed.
  @get @route("/{id}") @operationId("adminGetBlogPost") get(@path id: int32): BlogPost | UnauthorizedError | ForbiddenError | ApiError;

  @post @useAuth([AdminSession, XsrfToken]) @operationId("adminCreateBlogPost") create(@body body: BlogPostInput): {
    @statusCode statusCode: 201;
    @body post: BlogPost;
  } | ValidationError | UnauthorizedError | ForbiddenError | ApiError;

  @put @useAuth([AdminSession, XsrfToken]) @route("/{id}") @operationId("adminUpdateBlogPost") update(@path id: int32, @body body: BlogPostInput): BlogPost | ValidationError | UnauthorizedError | ForbiddenError | ApiError;

  @delete @useAuth([AdminSession, XsrfToken]) @route("/{id}") @operationId("adminDeleteBlogPost") delete(@path id: int32): NoContent | UnauthorizedError | ForbiddenError | ApiError;

  /** Set the related/promoted courses for a post. */
  @post @useAuth([AdminSession, XsrfToken]) @route("/{id}/related_courses") @operationId("adminSetBlogPostRelatedCourses")
  relatedCourses(@path id: int32, @body body: BlogPostRelatedCoursesInput): BlogPost | ValidationError | UnauthorizedError | ForbiddenError | ApiError;
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
@useAuth(AdminSession)
interface AdminReviews {
  @get @operationId("adminListReviews") list(...ListQuery): Page<Review> | UnauthorizedError | ForbiddenError | ApiError;

  /** Get a single review. A missing id surfaces as 404 via the central
      ent-error handler, not a typed union member. */
  @get @route("/{id}") @operationId("adminGetReview") get(@path id: int32): Review | UnauthorizedError | ForbiddenError | ApiError;

  @post @useAuth([AdminSession, XsrfToken]) @operationId("adminCreateReview") create(@body body: ReviewInput): {
    @statusCode statusCode: 201;
    @body review: Review;
  } | UnauthorizedError | ForbiddenError | ApiError;

  /** Update a review. A missing id surfaces as 404 via the central ent-error
      handler. */
  @put @useAuth([AdminSession, XsrfToken]) @route("/{id}") @operationId("adminUpdateReview") update(@path id: int32, @body body: ReviewInput): Review | UnauthorizedError | ForbiddenError | ApiError;

  @delete @useAuth([AdminSession, XsrfToken]) @route("/{id}") @operationId("adminDeleteReview") delete(@path id: int32): NoContent | UnauthorizedError | ForbiddenError | ApiError;
}

// ---------------------------------------------------------------------------
// Courses — index/get/create/update (no delete)
// + trigger review + create version.
// ---------------------------------------------------------------------------

model CourseInput {
  slug: string | null;
  learnAs: CourseLearnAs | null;
  progress: CourseReadiness | null;
  hexletProgramLandingPage: string | null;
  repositoryUrl: string | null;

  /** Attachment id from `POST /admin/attachments`, or null. */
  coverAttachmentId: int32 | null;
}

@route("/admin/courses")
@tag("Admin: Courses")
@useAuth(AdminSession)
interface AdminCourses {
  @get @operationId("adminListCourses") list(...ListQuery): Page<Course> | UnauthorizedError | ForbiddenError | ApiError;

  /** Get a single course. A missing id surfaces as 404 via the central
      ent-error handler. */
  @get @route("/{id}") @operationId("adminGetCourse") get(@path id: int32): Course | UnauthorizedError | ForbiddenError | ApiError;

  @post @useAuth([AdminSession, XsrfToken]) @operationId("adminCreateCourse") create(@body body: CourseInput): {
    @statusCode statusCode: 201;
    @body course: Course;
  } | UnauthorizedError | ForbiddenError | ApiError;

  @put @useAuth([AdminSession, XsrfToken]) @route("/{id}") @operationId("adminUpdateCourse") update(@path id: int32, @body body: CourseInput): Course | UnauthorizedError | ForbiddenError | ApiError;

  /** Enqueue AI re-review of every current lesson version. */
  @post @useAuth([AdminSession, XsrfToken]) @route("/{id}/review") @operationId("adminReviewCourse") review(@path id: int32): NoContent | NotFoundError | UnauthorizedError | ForbiddenError | ApiError;

  /** Build a new version of the course exercises. */
  @post @useAuth([AdminSession, XsrfToken]) @route("/{id}/versions") @operationId("adminCreateCourseVersion")
  createVersion(@path id: int32): {
    @statusCode statusCode: 201;
    @body version: CourseVersion;
  } | NotFoundError | UnauthorizedError | ForbiddenError | ApiError;
}

// ---------------------------------------------------------------------------
// Course landing pages — full CRUD
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

@route("/admin/course_landing_pages")
@tag("Admin: Course Landing Pages")
@useAuth(AdminSession)
interface AdminCourseLandingPages {
  @get @operationId("adminListCourseLandingPages") list(...ListQuery): Page<CourseLandingPage> | UnauthorizedError | ForbiddenError | ApiError;

  /** Get a single landing page. A missing id surfaces as 404 via the central
      ent-error handler, not a typed union member. */
  @get @route("/{id}") @operationId("adminGetCourseLandingPage") get(@path id: int32): CourseLandingPage | UnauthorizedError | ForbiddenError | ApiError;

  @post @useAuth([AdminSession, XsrfToken]) @operationId("adminCreateCourseLandingPage") create(@body body: CourseLandingPageInput): {
    @statusCode statusCode: 201;
    @body landingPage: CourseLandingPage;
  } | UnauthorizedError | ForbiddenError | ApiError;

  /** Update a landing page. A missing id surfaces as 404 via the central
      ent-error handler. */
  @put @useAuth([AdminSession, XsrfToken]) @route("/{id}") @operationId("adminUpdateCourseLandingPage") update(@path id: int32, @body body: CourseLandingPageInput): CourseLandingPage | UnauthorizedError | ForbiddenError | ApiError;

  @delete @useAuth([AdminSession, XsrfToken]) @route("/{id}") @operationId("adminDeleteCourseLandingPage") delete(@path id: int32): NoContent | UnauthorizedError | ForbiddenError | ApiError;
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
  @get @operationId("adminListLandingPageQnaItems") list(@path landingPageId: int32): QnaItem[] | UnauthorizedError | ForbiddenError | ApiError;

  @post @useAuth([AdminSession, XsrfToken]) @operationId("adminCreateLandingPageQnaItem") create(@path landingPageId: int32, @body body: QnaItemInput): {
    @statusCode statusCode: 201;
    @body item: QnaItem;
  } | UnauthorizedError | ForbiddenError | ApiError;

  /** Update a QnA item. A missing id (or one under a different parent) surfaces
      as 404 via the central ent-error handler. */
  @put @useAuth([AdminSession, XsrfToken]) @route("/{id}") @operationId("adminUpdateLandingPageQnaItem") update(@path landingPageId: int32, @path id: int32, @body body: QnaItemInput): QnaItem | UnauthorizedError | ForbiddenError | ApiError;

  @delete @useAuth([AdminSession, XsrfToken]) @route("/{id}") @operationId("adminDeleteLandingPageQnaItem") delete(@path landingPageId: int32, @path id: int32): NoContent | UnauthorizedError | ForbiddenError | ApiError;
}

@route("/admin/course_categories/{categoryId}/qna_items")
@tag("Admin: Course Categories")
@useAuth(AdminSession)
interface AdminCourseCategoryQnaItems {
  @get @operationId("adminListCategoryQnaItems") list(@path categoryId: int32): QnaItem[] | UnauthorizedError | ForbiddenError | ApiError;

  @post @useAuth([AdminSession, XsrfToken]) @operationId("adminCreateCategoryQnaItem") create(@path categoryId: int32, @body body: QnaItemInput): {
    @statusCode statusCode: 201;
    @body item: QnaItem;
  } | UnauthorizedError | ForbiddenError | ApiError;

  /** Update a QnA item. A missing id (or one under a different parent) surfaces
      as 404 via the central ent-error handler. */
  @put @useAuth([AdminSession, XsrfToken]) @route("/{id}") @operationId("adminUpdateCategoryQnaItem") update(@path categoryId: int32, @path id: int32, @body body: QnaItemInput): QnaItem | UnauthorizedError | ForbiddenError | ApiError;

  @delete @useAuth([AdminSession, XsrfToken]) @route("/{id}") @operationId("adminDeleteCategoryQnaItem") delete(@path categoryId: int32, @path id: int32): NoContent | UnauthorizedError | ForbiddenError | ApiError;
}

// ---------------------------------------------------------------------------
// Read-only admin lists (legacy index-only routes).
// ---------------------------------------------------------------------------

@route("/admin/course_lessons")
@tag("Admin: Course Lessons")
@useAuth(AdminSession)
interface AdminCourseLessons {
  @get @operationId("adminListCourseLessons") list(...ListQuery): Page<CourseLessonListItem> | UnauthorizedError | ForbiddenError | ApiError;

  /** Enqueue AI review for a single lesson. */
  @post @useAuth([AdminSession, XsrfToken]) @route("/{id}/review") @operationId("adminReviewCourseLesson") review(@path id: int32): NoContent | NotFoundError | UnauthorizedError | ForbiddenError | ApiError;
}

@route("/admin/lesson_progress")
@tag("Admin: Lesson Progress")
@useAuth(AdminSession)
interface AdminLessonProgress {
  @get @operationId("adminListLessonProgress") list(...ListQuery): Page<LessonProgress> | UnauthorizedError | ForbiddenError | ApiError;
}

@route("/admin/course_lesson_reviews")
@tag("Admin: Course Lesson Reviews")
@useAuth(AdminSession)
interface AdminCourseLessonReviews {
  @get @operationId("adminListCourseLessonReviews") list(...ListQuery): Page<CourseLessonReview> | UnauthorizedError | ForbiddenError | ApiError;
}

@route("/admin/leads")
@tag("Admin: Leads")
@useAuth(AdminSession)
interface AdminLeads {
  @get @operationId("adminListLeads") list(...ListQuery): Page<Lead> | UnauthorizedError | ForbiddenError | ApiError;
}

@route("/admin/messages")
@tag("Admin: Messages")
@useAuth(AdminSession)
interface AdminMessages {
  @get @operationId("adminListMessages") list(...ListQuery): Page<LessonAssistantMessage> | UnauthorizedError | ForbiddenError | ApiError;
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
  @get @operationId("adminListUsers") list(...ListQuery): Page<UserCrud> | UnauthorizedError | ForbiddenError | ApiError;

  /** Typeahead search by name/email. */
  @get @route("/search") @operationId("adminSearchUsers") search(@query q: string): UserCrud[] | UnauthorizedError | ForbiddenError | ApiError;

  /** Get a single user. A missing id surfaces as 404 via the central ent-error
      handler, not a typed union member. */
  @get @route("/{id}") @operationId("adminGetUser") get(@path id: int32): UserCrud | UnauthorizedError | ForbiddenError | ApiError;

  /** Create a user. A duplicate email is a DB unique constraint, surfaced as 409
      by the central ent-error handler. */
  @post @useAuth([AdminSession, XsrfToken]) @operationId("adminCreateUser") create(@body body: UserInput): {
    @statusCode statusCode: 201;
    @body user: UserCrud;
  } | UnauthorizedError | ForbiddenError | ApiError;

  /** Update a user. 404 (missing) and 409 (duplicate email) both flow through
      the central ent-error handler. */
  @put @useAuth([AdminSession, XsrfToken]) @route("/{id}") @operationId("adminUpdateUser") update(@path id: int32, @body body: UserInput): UserCrud | UnauthorizedError | ForbiddenError | ApiError;

  @delete @useAuth([AdminSession, XsrfToken]) @route("/{id}") @operationId("adminDeleteUser") delete(@path id: int32): NoContent | UnauthorizedError | ForbiddenError | ApiError;
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
  @get @operationId("adminListRoles") list(...ListQuery): Page<StaffRole> | UnauthorizedError | ForbiddenError | ApiError;

  /** Get a role with its permission matrix. A missing id surfaces as 404 via
      the central ent-error handler. */
  @get @route("/{id}") @operationId("adminGetRole") get(@path id: int32): StaffRoleDetail | UnauthorizedError | ForbiddenError | ApiError;

  /** Create a role. A duplicate name is a DB unique constraint, surfaced as 409
      by the central ent-error handler. */
  @post @useAuth([AdminSession, XsrfToken]) @operationId("adminCreateRole") create(@body body: RoleInput): {
    @statusCode statusCode: 201;
    @body role: StaffRoleDetail;
  } | UnauthorizedError | ForbiddenError | ApiError;

  @put @useAuth([AdminSession, XsrfToken]) @route("/{id}") @operationId("adminUpdateRole") update(@path id: int32, @body body: RoleInput): StaffRoleDetail | UnauthorizedError | ForbiddenError | ApiError;

  @delete @useAuth([AdminSession, XsrfToken]) @route("/{id}") @operationId("adminDeleteRole") delete(@path id: int32): NoContent | UnauthorizedError | ForbiddenError | ApiError;
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
@useAuth(AdminSession)
interface AdminRolePermissions {
  /** The permission matrix for a role. */
  @get @route("/{roleId}") @operationId("adminGetRolePermissions") get(@path roleId: int32): StaffRoleDetail | UnauthorizedError | ForbiddenError | ApiError;

  /** Replace the permission matrix for a role. */
  @put @useAuth([AdminSession, XsrfToken]) @route("/{roleId}") @operationId("adminUpdateRolePermissions") update(@path roleId: int32, @body body: RolePermissionsInput): StaffRoleDetail | UnauthorizedError | ForbiddenError | ApiError;
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
  @get @operationId("adminListStaffMembers") list(...ListQuery): Page<StaffMember> | UnauthorizedError | ForbiddenError | ApiError;

  @get @route("/{id}") @operationId("adminGetStaffMember") get(@path id: int32): StaffMember | UnauthorizedError | ForbiddenError | ApiError;

  @post @useAuth([AdminSession, XsrfToken]) @operationId("adminCreateStaffMember") create(@body body: StaffMemberInput): {
    @statusCode statusCode: 201;
    @body staffMember: StaffMember;
  } | UnauthorizedError | ForbiddenError | ApiError;

  @put @useAuth([AdminSession, XsrfToken]) @route("/{id}") @operationId("adminUpdateStaffMember") update(@path id: int32, @body body: StaffMemberInput): StaffMember | UnauthorizedError | ForbiddenError | ApiError;

  @delete @useAuth([AdminSession, XsrfToken]) @route("/{id}") @operationId("adminDeleteStaffMember") delete(@path id: int32): NoContent | UnauthorizedError | ForbiddenError | ApiError;
}

@route("/admin/management/users")
@tag("Admin: Users")
@useAuth(AdminSession)
interface AdminManagementUsers {
  @get @operationId("adminListManagementUsers") list(...ListQuery): Page<UserCrud> | UnauthorizedError | ForbiddenError | ApiError;

  /** Get a management user. A missing id surfaces as 404 via the central
      ent-error handler. */
  @get @route("/{id}") @operationId("adminGetManagementUser") get(@path id: int32): UserCrud | UnauthorizedError | ForbiddenError | ApiError;

  @put @useAuth([AdminSession, XsrfToken]) @route("/{id}") @operationId("adminUpdateManagementUser") update(@path id: int32, @body body: UserInput): UserCrud | UnauthorizedError | ForbiddenError | ApiError;
}
