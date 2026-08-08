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

/** Admin resources a staff role can be granted permissions on. */
enum PermissionResource {
  blogPosts: "blog_posts",
  banners: "banners",
  reviews: "reviews",
  leads: "leads",
  messages: "messages",
  courseCategories: "language_categories",
  courseLessons: "language_lessons",
  courseLessonReviews: "language_lesson_reviews",
  lessonProgress: "language_lesson_members",
  courses: "languages",
  courseLandingPages: "language_landing_pages",
}

// ---- Course core ----------------------------------------------------------

/** A build/version of a course's exercises (legacy `LanguageVersion`). */
model CourseVersion {
  id: int32;
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
  progress: CourseReadiness | null;
  categoryId: int32 | null;
  currentVersionId: int32 | null;
  currentVersion: CourseVersion | null;
  createdAt: utcDateTime;
  membersCount: int32;
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
  membersCount: int32;

  /** Estimated duration in hours (derived from lesson count). */
  duration: int32;

  /** Course cover image URL; null until assets are re-uploaded. */
  coverUrl: string | null;

  course: Course;
}

/** A learner's enrollment in a course (legacy `LanguageMember`, table
`language_members`). */
model Enrollment {
  id: int32;
  userId: int32;
  courseId: int32;
  state: EnrollmentState | null;

  /** Completion percentage (0–100). */
  progress: int32;

  nextLessonName: string | null;
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

  membersCount: int32;
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

/** A user row as shown in admin lists/forms (legacy: `UserCrud`). */
model UserCrud {
  id: int32;
  email: string | null;
  firstName: string | null;
  lastName: string | null;
  admin: boolean | null;
}

/** A registered passkey/WebAuthn credential (legacy: `UserCredential`). */
model UserCredential {
  id: int32;
  nickname: string | null;
  createdAt: utcDateTime;
}

/** A staff role (legacy: `StaffRole`). */
model StaffRole {
  id: int32;
  name: string;
  description: string | null;
  permissionsCount: int32;
  createdAt: utcDateTime;
}

/** A staff role with its full permission matrix (legacy: `StaffRoleCrud`). */
model StaffRoleDetail {
  id: int32;
  name: string;
  description: string | null;
  permissions: StaffRolePermission[];
}

/** One resource's permission bits for a role (legacy: `StaffRolePermission`). */
model StaffRolePermission {
  id: int32;
  roleId: int32;
  resource: PermissionResource;
  canIndex: boolean;
  canCreate: boolean;
  canUpdate: boolean;
  canDestroy: boolean;
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
  creator: User;
  name: string | null;
  slug: string | null;
  description: string | null;
  state: BlogPostState | null;
  locale: string | null;
  url: string;
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
  user: User;
  course: Course;
  userId: int32;
  courseId: int32;
  body: string | null;
  firstName: string | null;
  lastName: string | null;
  fullName: string | null;
  state: ReviewState | null;
  locale: ReviewLocale | null;

  /** Pinned reviews sort to the top of the public list. Part of the read model
      so the admin edit form can round-trip the flag it writes. */
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

// ---- Sitemap projections --------------------------------------------------

/** Minimal landing-page row for the sitemap (legacy: `LanguageSitemapLandingPage`). */
model SitemapCourseLandingPage {
  id: int32;
  courseId: int32;
  slug: string;
  header: string;
  locale: Locale;
}

/** Minimal lesson row for the sitemap (legacy: `LanguageSitemapLesson`). */
model SitemapCourseLesson {
  id: int32;
  courseId: int32;
  name: string;
  slug: string;
  naturalOrder: int32;
  locale: Locale;
}

/** Minimal blog-post row for the sitemap (legacy: `SitemapBlogPost`). */
model SitemapBlogPost {
  id: int32;
  name: string;
  slug: string;
  locale: Locale;
}
