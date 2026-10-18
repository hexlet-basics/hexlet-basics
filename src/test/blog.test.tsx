import { HttpResponse, http } from "msw";
import { describe, expect, it } from "vitest";
import { page } from "vitest/browser";
import type {
  BlogPost,
  BlogPostPage,
  BlogPostView,
  Course,
  CourseCatalogItem,
  User,
} from "@/client/types.gen";
import { Route as postRoute } from "@/routes/{-$locale}/blog_posts/$slug";
import { Route as indexRoute } from "@/routes/{-$locale}/blog_posts/index";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";
import type { RenderedRoute } from "@/test/renderRoute";

// The public blog, driven through its real routes with the API faked at the
// HTTP boundary. Asserted is what a reader sees: the posts, the body, the
// promoted courses, and what a like answers.

const author: User = {
  admin: false,
  assistantMessagesCount: null,
  canAccessAdmin: false,
  createdAt: "2026-01-01T00:00:00Z",
  createdAtAsTimestamp: null,
  email: "alice@example.com",
  firstName: "Alice",
  id: 1,
  lastName: null,
  name: "Alice",
  type: "user",
};

function blogPost(overrides: Partial<BlogPost>): BlogPost {
  return {
    coverUrl: null,
    createdAt: "2026-01-01T00:00:00Z",
    creator: author,
    description: "The very first post",
    id: 10,
    likesCount: 2,
    locale: "en",
    name: "Hello world",
    readingTime: 0,
    relatedCourseIds: [1],
    richBodyHtml: "<p>Hello <strong>world</strong> from the blog</p>",
    slug: "hello-world",
    state: "published",
    ...overrides,
  };
}

const course: Course = {
  categoryId: null,
  coverListVariant: null,
  coverThumbVariant: null,
  createdAt: "2026-01-01T00:00:00Z",
  currentVersion: null,
  currentVersionId: null,
  enrollmentsCount: 0,
  hexletProgramLandingPage: null,
  id: 1,
  learnAs: null,
  lessonsCount: 3,
  name: "ruby",
  ratingCount: 0,
  ratingValue: 0,
  readiness: null,
  repositoryUrl: null,
  slug: "ruby",
};

const rubyLanding: CourseCatalogItem = {
  course,
  coverUrl: null,
  duration: 3,
  enrollmentsCount: 42,
  header: "Ruby course",
  id: 100,
  locale: "en",
  name: "Ruby",
  slug: "ruby-en",
};

const post = blogPost({});
const newer = blogPost({ id: 11, likesCount: 0, name: "Learning Ruby", slug: "learning-ruby" });
const older = blogPost({ id: 9, likesCount: 0, name: "Oldest post", slug: "oldest-post" });

function postView(): BlogPostView {
  return { post, recommendedPosts: [newer, older], relatedLandingPages: [rubyLanding] };
}

function renderPost(user: User | null = null): Promise<RenderedRoute> {
  return renderRoute(postRoute, {
    initialPath: "/blog_posts/hello-world",
    path: "/{-$locale}/blog_posts/$slug",
    user,
  });
}

describe("blog", () => {
  it("the blog lists the posts the API returns", async () => {
    expect.hasAssertions();
    const list: BlogPostPage = { items: [newer, post], page: 1, perPage: 20, total: 2 };
    worker.use(http.get("*/api/blog_posts", () => HttpResponse.json(list)));

    const { router } = await renderRoute(indexRoute, {
      initialPath: "/blog_posts/",
      path: "/{-$locale}/blog_posts/",
    });

    await expect.element(page.getByRole("heading", { level: 1, name: "Blog" })).toBeVisible();
    await expect.element(page.getByRole("heading", { name: "Learning Ruby" })).toBeVisible();
    await expect
      .element(page.getByRole("link", { name: /Hello world/u }))
      .toHaveAttribute("href", "/blog_posts/hello-world");

    // The route's head: legacy title behind the site name, and its Twitter card.
    const meta = router.state.matches.at(-1)?.meta;
    expect(meta).toContainEqual({ title: "CodeBasics | Blog" });
    expect(meta).toContainEqual({ content: "@hexlet_io", name: "twitter:site" });
    expect(router.state.matches.at(-1)?.links).toContainEqual({
      href: `${globalThis.location.origin}/blog_posts`,
      rel: "canonical",
    });
  }, 60_000);

  it("a post shows its stored body, its courses and two more posts", async () => {
    expect.hasAssertions();
    worker.use(http.get("*/api/blog_posts/hello-world", () => HttpResponse.json(postView())));

    const { router } = await renderPost();

    await expect
      .element(page.getByRole("heading", { level: 1, name: "Hello world" }))
      .toBeVisible();
    await expect.element(page.getByText("world", { exact: true })).toBeVisible();
    await expect.element(page.getByRole("heading", { name: "Ruby" })).toBeVisible();
    await expect.element(page.getByRole("heading", { name: "Learning Ruby" })).toBeVisible();
    await expect.element(page.getByRole("heading", { name: "Oldest post" })).toBeVisible();

    // The route's head: title, canonical URL and the schema.org Article.
    const match = router.state.matches.at(-1);
    expect(match?.meta).toContainEqual({ title: "CodeBasics | Hello world" });
    expect(match?.links).toContainEqual({
      href: `${globalThis.location.origin}/blog_posts/hello-world`,
      rel: "canonical",
    });
    expect(match?.headScripts?.[0]?.children).toContain('"@type":"Article"');
  }, 60_000);

  it("a visitor's like leads to sign-in", async () => {
    expect.hasAssertions();
    worker.use(http.get("*/api/blog_posts/hello-world", () => HttpResponse.json(postView())));

    await renderPost();

    await expect
      .element(page.getByRole("link", { name: "like" }).first())
      .toHaveAttribute("href", "/session/new");
  }, 60_000);

  it("a repeat like is answered as already counted", async () => {
    expect.hasAssertions();
    let likes = 0;
    worker.use(
      http.get("*/api/blog_posts/hello-world", () => HttpResponse.json(postView())),
      http.post("*/api/blog_posts/10/likes", () => {
        likes += 1;
        // The first like counts; the second changes nothing.
        return HttpResponse.json({ ...post, likesCount: 3 }, { status: 201 });
      }),
    );

    await renderPost(author);

    const like = page.getByRole("button", { name: "like" }).first();
    await like.click();
    await expect.element(page.getByText("Спасибо за лайк!")).toBeVisible();

    await like.click();
    await expect.element(page.getByText("Ваш лайк уже засчитан :)")).toBeVisible();
    expect(likes).toBe(2);
  }, 60_000);

  it("a post that is not published here is not found", async () => {
    expect.hasAssertions();
    worker.use(
      http.get("*/api/blog_posts/hello-world", () =>
        HttpResponse.json({ status: 404, title: "Not Found" }, { status: 404 }),
      ),
    );

    await renderPost();

    await expect.element(page.getByRole("heading", { name: "Page Not Found" })).toBeVisible();
  }, 60_000);
});
