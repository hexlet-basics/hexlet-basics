import { http, HttpResponse } from "msw";
import { expect, test } from "vitest";
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

// The public blog, driven through its real routes with the API faked at the
// HTTP boundary. Asserted is what a reader sees: the posts, the body, the
// promoted courses, and what a like answers.

const author: User = {
  id: 1,
  firstName: "Alice",
  lastName: null,
  name: "Alice",
  email: "alice@example.com",
  admin: false,
  canAccessAdmin: false,
  assistantMessagesCount: null,
  createdAt: "2026-01-01T00:00:00Z",
  createdAtAsTimestamp: null,
  type: "user",
};

function blogPost(overrides: Partial<BlogPost>): BlogPost {
  return {
    id: 10,
    creator: author,
    name: "Hello world",
    slug: "hello-world",
    description: "The very first post",
    state: "published",
    locale: "en",
    url: "https://code-basics.com/blog_posts/hello-world",
    richBodyHtml: "<p>Hello <strong>world</strong> from the blog</p>",
    readingTime: 0,
    likesCount: 2,
    relatedCourseItemsCount: 1,
    relatedCourseIds: [1],
    coverThumbVariant: null,
    coverListVariant: null,
    coverMainVariant: null,
    createdAt: "2026-01-01T00:00:00Z",
    ...overrides,
  };
}

const course: Course = {
  id: 1,
  slug: "ruby",
  name: "ruby",
  learnAs: null,
  readiness: null,
  categoryId: null,
  currentVersionId: null,
  currentVersion: null,
  createdAt: "2026-01-01T00:00:00Z",
  enrollmentsCount: 0,
  lessonsCount: 3,
  ratingCount: 0,
  ratingValue: 0,
  repositoryUrl: null,
  hexletProgramLandingPage: null,
  coverListVariant: null,
  coverThumbVariant: null,
};

const rubyLanding: CourseCatalogItem = {
  id: 100,
  slug: "ruby-en",
  header: "Ruby course",
  name: "Ruby",
  locale: "en",
  enrollmentsCount: 42,
  duration: 3,
  coverUrl: null,
  course,
};

const post = blogPost({});
const newer = blogPost({ id: 11, name: "Learning Ruby", slug: "learning-ruby", likesCount: 0 });
const older = blogPost({ id: 9, name: "Oldest post", slug: "oldest-post", likesCount: 0 });

function postView(): BlogPostView {
  return { post, recommendedPosts: [newer, older], relatedLandingPages: [rubyLanding] };
}

function renderPost(user: User | null = null) {
  return renderRoute(postRoute, {
    path: "/{-$locale}/blog_posts/$slug",
    initialPath: "/blog_posts/hello-world",
    user,
  });
}

test("the blog lists the posts the API returns", async () => {
  const list: BlogPostPage = { items: [newer, post], total: 2, page: 1, perPage: 20 };
  worker.use(http.get("*/api/blog_posts", () => HttpResponse.json(list)));

  await renderRoute(indexRoute, { path: "/{-$locale}/blog_posts/", initialPath: "/blog_posts/" });

  await expect.element(page.getByRole("heading", { name: "Blog", level: 1 })).toBeVisible();
  await expect.element(page.getByRole("heading", { name: "Learning Ruby" })).toBeVisible();
  await expect
    .element(page.getByRole("link", { name: /Hello world/ }))
    .toHaveAttribute("href", "/blog_posts/hello-world");
});

test("a post shows its stored body, its courses and two more posts", async () => {
  worker.use(http.get("*/api/blog_posts/hello-world", () => HttpResponse.json(postView())));

  await renderPost();

  await expect.element(page.getByRole("heading", { name: "Hello world", level: 1 })).toBeVisible();
  await expect.element(page.getByText("world", { exact: true })).toBeVisible();
  await expect.element(page.getByRole("heading", { name: "Ruby" })).toBeVisible();
  await expect.element(page.getByRole("heading", { name: "Learning Ruby" })).toBeVisible();
  await expect.element(page.getByRole("heading", { name: "Oldest post" })).toBeVisible();
});

test("a visitor's like leads to sign-in", async () => {
  worker.use(http.get("*/api/blog_posts/hello-world", () => HttpResponse.json(postView())));

  await renderPost();

  await expect
    .element(page.getByRole("link", { name: "like" }).first())
    .toHaveAttribute("href", "/session/new");
});

test("a repeat like is answered as already counted", async () => {
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
});

test("a post that is not published here is not found", async () => {
  worker.use(
    http.get("*/api/blog_posts/hello-world", () =>
      HttpResponse.json({ status: 404, title: "Not Found" }, { status: 404 }),
    ),
  );

  await renderPost();

  await expect.element(page.getByText("Not Found")).toBeVisible();
});
