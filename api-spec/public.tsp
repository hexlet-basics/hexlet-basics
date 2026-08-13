import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi";

using Http;
using Rest;
using OpenAPI;

namespace HexletBasics;

// ---------------------------------------------------------------------------
// Public read / SEO surface.
//
// Legacy URLs are preserved (ADR-0002): `/languages`, `/language_categories`,
// locale prefixes handled by the SSR router, not the contract. Reads are
// unauthenticated; writes here are the marketing entry points (leads, likes).
// ---------------------------------------------------------------------------

/** Everything the public course landing needs in one payload. */
model CourseView {
  course: Course;
  landingPage: CourseLandingPage | null;
  lessons: CourseLessonListItem[];

  /** The current user's enrollment, when signed in. */
  enrollment: Enrollment | null;

  /**
   * Where the visitor stands in this course. Identical in shape for a
   * signed-in learner, whose position is rows, and for a guest, whose position
   * is the signed progress cookie — which is what keeps the client free of any
   * guest/member branch. Null only for a course with no built version, where
   * there is nothing to stand in.
   */
  progress: CourseProgress | null;
}

/** A category page with the courses it groups. */
model CourseCategoryView {
  category: CourseCategory;
  landingPages: CourseCatalogItem[];
}

/** A static content page (about, authors, privacy, tos, cookie). */
model PageContent {
  slug: string;
  title: string;
  bodyHtml: string;
}

@route("/languages")
@tag("Courses")
interface Courses {
  /** List the published course catalog. */
  @get
  @operationId("listCourses")
  list(): CourseCatalogItem[] | ApiError;

  /** Course landing page by slug. */
  @get
  @route("/{slug}")
  @operationId("getCourse")
  get(@path slug: string): CourseView | NotFoundError | ApiError;
}

/** The lesson page: the lesson itself, and where the visitor stands in it. */
model CourseLessonView {
  lesson: CourseLesson;

  /**
   * The course's landing copy — the same page the course landing read
   * resolves: the one flagged main, otherwise the oldest.
   *
   * The player titles the page with the landing page's name, not the course's
   * own `name`: those are two different strings, and the marketing name is the
   * one a learner recognises. The description of the page is composed from it
   * too, which is why the whole page travels rather than just its name.
   *
   * Null for a course with no landing page, which is a course nobody has
   * written copy for yet — not a reason to fail the read.
   */
  landingPage: CourseLandingPage | null;

  /**
   * Every Lesson of the course's current version, in course order — the same
   * list the course landing read returns.
   *
   * The player's navigation tab is this list joined to `progress.lessons` by
   * slug: names and order from here, checks and locks from there. Both screens
   * render from the same two arrays, which is what keeps the Course page's list
   * and the player's from disagreeing.
   */
  lessons: CourseLessonListItem[];

  /**
   * Where the visitor stands in the course this lesson belongs to — the same
   * payload the course landing read returns, for a signed-in learner and a
   * guest alike. The player reads `lessons[].available` to know whether this
   * one may be taken, and `nextLessonSlug` for where "next" goes.
   *
   * Null only for a course with no built version.
   */
  progress: CourseProgress | null;
}

@route("/languages/{courseSlug}/lessons")
@tag("Courses")
interface CourseLessons {
  /**
   * Lesson player payload (theory, starter code, tests) by slug, plus the
   * visitor's progress.
   *
   * Reading a lesson never starts it: the frontend preloads routes on hover,
   * so a read with that effect would enroll a learner in every lesson they
   * pointed at (ADR-0012). Theory stays public and indexable for everyone,
   * including a lesson the visitor may not take yet.
   */
  @get
  @route("/{slug}")
  @operationId("getCourseLesson")
  get(@path courseSlug: string, @path slug: string): CourseLessonView | NotFoundError | ApiError;
}

@route("/language_categories")
@tag("Course Categories")
interface PublicCourseCategories {
  /** List published categories. */
  @get
  @operationId("listPublicCourseCategories")
  list(): CourseCategory[] | ApiError;

  /** A category and the courses it groups. */
  @get
  @route("/{slug}")
  @operationId("getPublicCourseCategory")
  get(@path slug: string): CourseCategoryView | NotFoundError | ApiError;
}

@route("/blog_posts")
@tag("Blog")
interface Blog {
  /** Paginated published blog posts. */
  @get
  @operationId("listBlogPosts")
  list(...ListQuery): Page<BlogPost> | ApiError;

  /** A single blog post by slug. */
  @get
  @route("/{slug}")
  @operationId("getBlogPost")
  get(@path slug: string): BlogPost | NotFoundError | ApiError;

  /** The next post to read after this one. */
  @get
  @route("/{id}/next")
  @operationId("getNextBlogPost")
  next(@path id: int32): BlogPost | NotFoundError | ApiError;

  /** Like a post (idempotent per visitor). */
  @post
  @useAuth([UserSession, XsrfToken])
  @route("/{id}/likes")
  @operationId("likeBlogPost")
  like(@path id: int32): {
    @statusCode statusCode: 201;
    @body body: BlogPost;
  } | NotFoundError | UnauthorizedError | ApiError;
}

@route("/reviews")
@tag("Reviews")
interface PublicReviews {
  /** Paginated published student reviews. */
  @get
  @operationId("listPublicReviews")
  list(...ListQuery): Page<Review> | ApiError;
}

@route("/pages")
@tag("Pages")
interface Pages {
  /** A static content page by slug. */
  @get
  @route("/{slug}")
  @operationId("getPage")
  get(@path slug: string): PageContent | NotFoundError | ApiError;
}

// ---- Marketing writes -----------------------------------------------------

/** A lead captured from a marketing form (legacy: `LeadCrud`). */
model LeadInput {
  contactMethod: "telegram" | "phone" | "whatsapp";
  @minLength(1) contactValue: string;
  ymClientId: string | null;
}

@route("/leads")
@tag("Leads")
@useAuth([UserSession, XsrfToken])
interface PublicLeads {
  /** Submit a contact request. */
  @post
  @operationId("createLead")
  create(@body body: LeadInput): {
    @statusCode statusCode: 201;
    @body lead: Lead;
  } | ValidationError | UnauthorizedError | ApiError;
}

/** A request to download the free book (ru marketing funnel). */
model BookRequestInput {
  @minLength(1) email: string;
  fullName: string | null;
}

@route("/book")
@tag("Marketing")
@useAuth([UserSession, XsrfToken])
interface Book {
  /** Request the book download link by email. */
  @post
  @route("/create_request")
  @operationId("createBookRequest")
  createRequest(@body body: BookRequestInput): NoContent | ValidationError | UnauthorizedError | ApiError;
}

// ---- Sitemap --------------------------------------------------------------

/** Aggregated data for the XML sitemap / SEO index. */
model Sitemap {
  landingPages: SitemapCourseLandingPage[];
  lessons: SitemapCourseLesson[];
  blogPosts: SitemapBlogPost[];
}

@route("/map")
@tag("Sitemap")
interface SitemapFeed {
  /** Everything the sitemap generator needs. */
  @get
  @operationId("getSitemap")
  get(): Sitemap | ApiError;
}
