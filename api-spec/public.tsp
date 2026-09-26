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

  /**
   * The modules of the course's current version in the request locale, in
   * course order — the learning program's accordion. Each names its lessons by
   * slug, joined to `lessons` for names, so the flat list stays the one source
   * of lesson rows the page and the player share.
   */
  modules: CourseModuleListItem[];

  /** The landing page's questions and answers, oldest first; empty without a
      landing page. */
  qnaItems: QnaItem[];

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

  /** The category's listed, published landing pages in the request locale
      whose course is completed, in course order. */
  landingPages: CourseCatalogItem[];

  /** The category's questions and answers, oldest first. */
  qnaItems: QnaItem[];
}

/** A blog post page: the post plus what the page shows under it. */
model BlogPostView {
  post: BlogPost;

  /** Two other published posts in the request locale, newest first. */
  recommendedPosts: BlogPost[];

  /** The main landing page, in the request locale, of each course the post
      promotes, in the post's display order. */
  relatedLandingPages: CourseCatalogItem[];
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
  /** List the categories of the request locale. */
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

  /** A published post in the request locale, by slug, with its page data. */
  @get
  @route("/{slug}")
  @operationId("getBlogPost")
  get(@path slug: string): BlogPostView | NotFoundError | ApiError;

  /** The next post to read after this one: the next older published post in
      the request locale. */
  @get
  @route("/{id}/next")
  @operationId("getNextBlogPost")
  next(@path id: int32): BlogPost | NotFoundError | ApiError;

  /** Like a post once per signed-in user; a repeat like changes nothing. */
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

// ---- Marketing writes -----------------------------------------------------

/** A lead captured from a marketing form (legacy: `LeadCrud`). */
model LeadInput {
  contactMethod: "telegram" | "phone" | "whatsapp";
  @minLength(1) contactValue: string;
  ymClientId: string | null;

  /**
   * Where the visitor first came from. Legacy read this off ahoy's visit; ahoy
   * is not ported (ADR-0015), so the frontend keeps the first visit in a cookie
   * and sends it with the form. Absent when the browser kept no cookie (blocked or cleared).
   */
  firstVisit?: FirstVisit | null;
}

/** The first visit's traffic source, as recorded by the browser. */
model FirstVisit {
  utmSource: string | null;
  utmMedium: string | null;
  utmCampaign: string | null;
  utmContent: string | null;
  utmTerm: string | null;
  /** The full URL the visit landed on, query string included. */
  landingPage: string | null;
  referrer: string | null;
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

/** The free book page (ru marketing funnel, legacy `books#show`). */
model BookView {
  /** Whether the signed-in visitor has requested the book, in any state; false for a visitor. */
  requested: boolean;
}

@route("/book")
@tag("Marketing")
interface Book {
  /** The book page's state for whoever is visiting; a visitor is answered too. */
  @get
  @operationId("getBook")
  get(): BookView | ApiError;

  /** Request the book for the signed-in user. A repeat request changes nothing. */
  @post
  @useAuth([UserSession, XsrfToken])
  @route("/create_request")
  @operationId("createBookRequest")
  createRequest(): NoContent | UnauthorizedError | ApiError;

  /**
   * Download the book: mark the request downloaded and redirect to the PDF in
   * blob storage, or back to the book page when there is no request yet. The
   * browser follows it as a link, so only the session cookie is required — a
   * navigation cannot carry the XSRF header, and a GET changes nothing unsafe.
   */
  @get
  @useAuth(UserSession)
  @route("/download")
  @operationId("downloadBook")
  download(): {
    @statusCode statusCode: 302;
    @header("Location") location: string;
  } | UnauthorizedError | ApiError;
}

// ---- Sitemap --------------------------------------------------------------

/**
 * The rows behind the HTML sitemap page (legacy `HomeController#sitemap`,
 * `/ru/map`). Flat lists carrying their own locale; the page groups them.
 *
 * Landing pages come already narrowed to those whose course has lessons in
 * the page's locale: legacy sent the lessons only for that check (their list
 * on the page was commented out), so the filter runs in SQL instead of
 * shipping every lesson.
 */
model Sitemap {
  landingPages: SitemapCourseLandingPage[];
  blogPosts: SitemapBlogPost[];
  categories: CourseCategory[];
}

@route("/map")
@tag("Sitemap")
interface SitemapFeed {
  /** Everything the sitemap page lists, across the ru and en locales. */
  @get
  @operationId("getSitemap")
  get(): Sitemap | ApiError;
}

// ---- Feeds ----------------------------------------------------------------

/**
 * The Yandex course catalogue: a YML document (Yandex's own XML dialect) of
 * the ru courses, which Yandex polls to list them in its course results.
 *
 * `bytes` rather than a modelled XML body: ogen has no XML codec and generates
 * a non-JSON response only as a stream (`format: binary`), so the document is
 * built with `encoding/xml` in `internal/feeds` and handed over as a reader.
 */
model YandexCoursesFeed {
  @header contentType: "application/xml";
  @body body: bytes;
}

/**
 * Machine-read feeds, at the addresses their consumers are registered with —
 * the URL is the contract here, so it keeps the legacy path exactly
 * (ADR-0015).
 */
@route("/feeds")
@tag("Feeds")
interface Feeds {
  /**
   * The feed at the address production actually answers on. Legacy routes
   * `/api` with a JSON default format and the action only knows XML, so the
   * bare path has always answered 406 and Yandex reads this one.
   */
  @get
  @route("/yandex_courses.xml")
  @operationId("getYandexCoursesFeedXml")
  getXml(): YandexCoursesFeed | ApiError;

  /** The same feed at the path the legacy route declares. */
  @get
  @route("/yandex_courses")
  @operationId("getYandexCoursesFeed")
  get(): YandexCoursesFeed | ApiError;
}
