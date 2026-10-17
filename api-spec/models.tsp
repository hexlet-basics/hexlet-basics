import "@typespec/http";

using Http;

namespace HexletBasics;

// ---------------------------------------------------------------------------
// Domain read models (API response shapes) and domain enums.
//
// Every model is the DOMAIN projection (`Course`, not `Language`) of a legacy
// typelizer serializer; the legacy serializer name is noted for traceability.
// Fields are camelCase (matches the existing contract); the Go/ogen side maps
// snake_case columns to these. Write models (`*Input`) live next to their
// operations in auth.tsp / admin.tsp — they are deliberately a subset.
// ---------------------------------------------------------------------------

// ---- Enums ----------------------------------------------------------------

/** How a course is meant to be learned. */
enum CourseLearnAs {
  firstLanguage: "first_language",
  secondLanguage: "second_language",
}

/** Publication/readiness state of a course. */
enum CourseReadiness {
  completed: "completed",
  inDevelopment: "in_development",
  draft: "draft",
}

/** Enrollment and lesson-progress lifecycle (legacy `MemberState`). */
enum EnrollmentState {
  started: "started",
  finished: "finished",
}

/** Publication state shared by landing pages. */
enum LandingPageState {
  draft: "draft",
  archived: "archived",
  published: "published",
}

/** Banner display background. */
enum BannerBackground {
  ctaGradient: "cta_gradient",
  dark: "dark",
  blue: "blue",
}

/** Banners only target ru/en today (legacy `BannerLocale`). */
enum BannerLocale {
  ru: "ru",
  en: "en",
}

/** Banner publication state. */
enum BannerState {
  draft: "draft",
  published: "published",
  archived: "archived",
}

/** Blog post publication state. */
enum BlogPostState {
  draft: "draft",
  published: "published",
  archived: "archived",
}

/** Reviews only target ru/en today (legacy `ReviewLocale`). */
enum ReviewLocale {
  ru: "ru",
  en: "en",
}

/** Review publication state. */
enum ReviewState {
  draft: "draft",
  published: "published",
  archived: "archived",
}

// ---- Course core ----------------------------------------------------------

/** A build/version of a course's exercises (legacy `LanguageVersion`). */
model CourseVersion {
  id: int32;

  /**
   * The name the course's spec gives this build (e.g. "JavaScript"). Legacy's
   * lesson page opens its meta description with it, in brackets.
   */
  name: string | null;

  result: string | null;
  state: string | null;
  createdAt: utcDateTime;
}

/**
 * A course (legacy: `Language`). Superset read model: the public catalog embeds
 * it, admin course screens read the full shape.
 */
model Course {
  id: int32;
  slug: string;
  name: string | null;
  learnAs: CourseLearnAs | null;
  readiness: CourseReadiness | null;
  categoryId: int32 | null;
  currentVersionId: int32 | null;
  currentVersion: CourseVersion | null;
  createdAt: utcDateTime;
  enrollmentsCount: int32;
  lessonsCount: int32;
  ratingCount: int32;
  ratingValue: float64;
  repositoryUrl: string | null;
  hexletProgramLandingPage: string | null;
  coverListVariant: string | null;
  coverThumbVariant: string | null;
}

/**
 * A catalog entry as shown on the courses index (legacy:
 * `LanguageLandingPageForLists`). Embeds the full course plus derived fields.
 */
model CourseCatalogItem {
  id: int32;
  slug: string;
  header: string | null;
  name: string | null;
  locale: string | null;
  enrollmentsCount: int32;

  /** Estimated duration in hours (derived from lesson count). */
  duration: int32;

  /** Course cover image URL; null until assets are re-uploaded. */
  coverUrl: string | null;

  course: Course;
}

/** A learner's enrollment in a course (table `language_members`). */
model Enrollment {
  id: int32;
  userId: int32;
  courseId: int32;
  state: EnrollmentState | null;

  /** Completion: the share of the course's current lessons finished (0–100). */
  completion: int32;

  nextLessonName: string | null;

  /** Where the learner stands, and what they may take next. */
  progress: CourseProgress;
}

/**
 * Where a learner stands in one Course, computed against its current Version:
 * which Lessons they have finished, in whatever order they took them.
 */
model CourseProgress {
  state: EnrollmentState | null;

  /** The share of the Course's current Lessons finished (0–100). */
  completion: int32;

  /** The first unfinished Lesson in course order — where the learner resumes. */
  nextLessonSlug: string | null;

  lessons: LessonProgressItem[];
}

/** Per-Lesson state for the course and lesson pages: the checkmarks. */
model LessonProgressItem {
  slug: string;

  /** Position in the current Version, 1..N across the whole Course. */
  position: int32;

  finished: boolean;
}

// ---- Course categories & landing pages ------------------------------------

/** A course category (legacy: `LanguageCategory` / table `language_categories`). */
model CourseCategory {
  id: int32;
  slug: string | null;
  name: string | null;
  header: string | null;
  description: string | null;
  locale: string | null;
  createdAt: utcDateTime;
}

/** A Q&A entry attached to a category or landing page (legacy QnA item). */
model QnaItem {
  id: int32;
  question: string;
  answer: string;
}

/** A marketing landing page for a course (legacy: `LanguageLandingPage`). */
model CourseLandingPage {
  id: int32;
  courseId: int32;
  courseSlug: string;
  createdAt: utcDateTime;
  slug: string;
  name: string;
  main: boolean | null;
  listed: boolean | null;
  state: LandingPageState | null;

  /** Free-form ordering key (legacy stores it as a string). */
  order: string | null;

  /** Footer visibility + link caption. Part of the read model so the admin
      edit form can round-trip the fields it writes. */
  footer: boolean | null;
  footerName: string | null;
  landingPageToRedirectId: int32 | null;

  metaTitle: string;
  metaDescription: string;
  header: string;
  description: string;
  usedInHeader: string | null;
  usedInDescription: string | null;
  outcomesHeader: string | null;
  outcomesDescription: string | null;
  outcomesImage: string | null;

  /** Estimated duration in hours. */
  duration: int32;

  enrollmentsCount: int32;
}

// ---- Lessons --------------------------------------------------------------

/** A lesson definition entry (name + description pair). */
model LessonDefinition {
  name: string;
  description: string;
}

/** A single lesson (legacy: `LanguageLesson`). */
model CourseLesson {
  course: Course;
  id: int32;
  name: string | null;
  slug: string;
  locale: string | null;
  naturalOrder: int32;
  versionId: int32;
  version: int32 | null;
  description: string | null;
  instructions: string | null;
  theory: string | null;
  definitions: LessonDefinition[];
  tips: string[];
  preparedCode: string | null;
  originalCode: string | null;
  testCode: string | null;

  /** File extension of the solution file, without the dot (`rb`, `md`), from
      the course build's spec.yml. Highlighting follows it, not the course slug:
      a course's code and tests need not be in the language its slug names. Null
      when the build records none. */
  exerciseExtension: string | null;

  /** File extension of the test file, without the dot (`js`), taken from the
      build's `exercise_test_filename`. Null when that name has no extension. */
  testExtension: string | null;
  sourceCodeUrl: string | null;
  createdAt: utcDateTime;
}

/** A lesson as shown in admin lists (legacy: `LanguageLessonForLists`). */
model CourseLessonListItem {
  id: int32;
  name: string | null;
  description: string | null;
  slug: string;
}

/** A module of a course version in one locale (legacy:
`Language::Module::Version::Info`, as `Language::ModuleResource` served it). */
model CourseModuleListItem {
  /** The module's stable id, as legacy's resource exposed it. */
  id: int32;

  name: string | null;
  description: string | null;

  /** The module's lessons in this version, in course order, by slug. */
  lessonSlugs: string[];
}

/** A learner's progress through one lesson (legacy: `LanguageLessonMember`,
table `language_lesson_members`). */
model LessonProgress {
  id: int32;
  userId: int32;
  state: EnrollmentState;
  messagesCount: int32 | null;
  createdAt: utcDateTime;
  courseSlug: string;
  courseLessonSlug: string;
  courseLessonName: string;
}

/** An AI-generated lesson review record (legacy: `LanguageLessonReview`). */
model CourseLessonReview {
  id: int32;
  locale: string;
  courseId: int32;
  courseLessonId: int32;
  courseLessonVersionId: int32;
  courseLessonVersionInfoId: int32;
  summary: string;
  slug: string;
  lessonNaturalOrder: int32;
  courseSlug: string;
  createdAt: utcDateTime;
}

// ---- People ---------------------------------------------------------------

/** A user (legacy: `User`). Passwords are never part of a response. */
model User {
  id: int32;
  firstName: string | null;
  lastName: string | null;
  name: string | null;
  email: string | null;
  admin: boolean | null;
  canAccessAdmin: boolean;
  assistantMessagesCount: int32 | null;
  createdAt: utcDateTime;
  createdAtAsTimestamp: int64 | null;
  type: "user";
}

/** The part of a user that is safe to show anonymous visitors: the display
    name only, never the email or the admin flag. Embedded in public review and
    blog responses in place of `User`. */
model PublicUser {
  id: int32;
  firstName: string | null;
  lastName: string | null;
  name: string | null;
}

/** A user row as shown in admin lists/forms (legacy: `UserCrud`). */
model UserCrud {
  id: int32;
  email: string | null;
  firstName: string | null;
  lastName: string | null;
  admin: boolean | null;
}

/** A staff role (legacy: `StaffRole`). */
model StaffRole {
  id: int32;
  name: string;
  description: string | null;
  permissionsCount: int32;
  createdAt: utcDateTime;
}

/** A staff role with the Permissions it holds (legacy: `StaffRoleCrud`). */
model StaffRoleDetail {
  id: int32;
  name: string;
  description: string | null;

  /** Names of the declared Permissions the role holds (see `Permission`). */
  grants: string[];
}

/** Whether an operation only reads or can change data (ADR-0017). */
enum PermissionKind {
  read: "read",
  write: "write",
}

/**
 * A Permission an admin operation can require, as declared on the operations in
 * this contract with `x-permission`. The catalog is read from those
 * declarations, never kept by hand, plus the capability Permissions the contract
 * lists under `x-capabilities`.
 */
model Permission {
  name: string;
  kind: PermissionKind;

  /**
   * True for a capability Permission (MCP, audit): the administrator flag does
   * not imply it, so a role must be given it explicitly.
   */
  capability: boolean;
}

/** A staff member: a user bound to a role (legacy: `StaffMember`). */
model StaffMember {
  id: int32;
  userId: int32;
  roleId: int32;
  user: User;
  role: StaffRole;
  allowedLocales: string[];
}

// ---- Content: banners, blog, reviews, leads -------------------------------

/** A site banner (legacy: `Banner`). */
model Banner {
  id: int32;
  locale: BannerLocale;
  body: string;
  url: string | null;
  background: BannerBackground;
  state: BannerState;
  startsAt: utcDateTime | null;
  finishesAt: utcDateTime | null;
  createdAt: utcDateTime;
}

/** A blog post (legacy: `BlogPost`). */
model BlogPost {
  id: int32;
  creator: PublicUser;
  name: string | null;
  slug: string | null;
  description: string | null;
  state: BlogPostState | null;
  locale: string | null;
  richBodyHtml: string;
  readingTime: int32;
  likesCount: int32;
  relatedCourseItemsCount: int32;

  /** Promoted-course ids in display order, so the admin related-courses
      editor can round-trip the set it writes. */
  relatedCourseIds: int32[];

  coverThumbVariant: string | null;
  coverListVariant: string | null;
  coverMainVariant: string | null;
  createdAt: utcDateTime;
}

/** A student review/testimonial (legacy: `Review`). */
model Review {
  id: int32;
  user: PublicUser;
  course: Course;
  userId: int32;
  courseId: int32;
  body: string | null;
  firstName: string | null;
  lastName: string | null;
  fullName: string | null;
  state: ReviewState | null;
  locale: ReviewLocale | null;

  /** Stored for the admin edit form to round-trip; it does not affect the
      public list, which legacy ordered newest first only. */
  pinned: boolean | null;

  createdAt: utcDateTime;
}

/** A sales lead (legacy: `Lead`). */
model Lead {
  id: int32;
  userId: int32;
  fullName: string | null;
  email: string | null;
  phone: string | null;
  whatsapp: string | null;
  telegram: string | null;
  surveyAnswersData: string | null;
  coursesData: string | null;
  createdAt: utcDateTime;
}

// ---- Audit trail (ADR-0017) -------------------------------------------------

/**
 * The old and new value of one field an operation changed. `old` is absent on a
 * create and `new` on a delete. A `secret` field (password material, a token)
 * carries its name only: no value was ever stored. The values are JSON of any
 * shape, as the field held them.
 */
model AuditFieldChange {
  name: string;
  old?: unknown;
  new?: unknown;
  secret?: boolean;
}

/** One mutation of one entity inside an operation. A bulk change carries no id. */
model AuditChange {
  /** The ent entity type, e.g. `BlogPost`. */
  entity: string;
  id?: string;

  /** create, update or delete. */
  `op`: string;
  fields?: AuditFieldChange[];
}

/**
 * What a read of personal data touched: how many rows and, up to a cap, their
 * ids. Identifiers and counts only, never a field value.
 */
model AuditReads {
  count: int32;
  ids?: string[];
}

/**
 * One Audit Record: who did what, through which channel, to which entity, and
 * how it ended (legacy: none, the trail is new). A change carries `changes`, a
 * read of personal data carries `reads`, a refusal carries neither.
 */
model AuditRecord {
  id: int64;
  occurredAt: utcDateTime;

  /** Null for a refusal of a caller who was not signed in. */
  actorUserId: int32 | null;
  channel: "web" | "mcp";
  connectedAppId: int32 | null;
  tokenId: int32 | null;
  requestId: string | null;
  address: string | null;
  userAgent: string | null;

  /** The operation name, e.g. `AdminUpdateBlogPost`. */
  operation: string;
  entityType: string | null;
  entityId: string | null;
  locale: string | null;
  outcome: "succeeded" | "forbidden" | "failed";
  changes: AuditChange[];
  reads: AuditReads | null;
}

// ---- Sitemap projections --------------------------------------------------

/** Minimal landing-page row for the sitemap (legacy: `LanguageSitemapLandingPage`). */
model SitemapCourseLandingPage {
  id: int32;
  courseId: int32;
  slug: string;
  header: string;
  locale: Locale;
}

/** Minimal blog-post row for the sitemap (legacy: `SitemapBlogPost`). */
model SitemapBlogPost {
  id: int32;
  name: string;
  slug: string;
  locale: Locale;
}
