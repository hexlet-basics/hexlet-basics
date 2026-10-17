import { test } from "../support/test";

test("the catalog lists the fixture courses", async ({ site }) => {
  await site.catalog.open();

  await site.expectHeading("Programming language courses", 1);
  await site.catalog.expectCourseCard("Ruby Course");
});

test("a course page shows its landing copy and syllabus, and a lesson opens from it", async ({
  site,
}) => {
  await site.catalog.openCoursePage("/ru/languages/javascript-ru");

  await site.expectHeading("Курс JavaScript", 1);
  await site.expectHeading("Программа обучения JavaScript", 2);

  await site.catalog.openLesson("Привет, мир!");

  await site.expectPath("/ru/languages/javascript/lessons/hello-world");
  await site.expectHeading("JavaScript: Привет, мир!", 1);
});

// The one test that follows a link the site itself renders: the API resolves a
// course page by its landing page's slug, so a link built from anything else
// (the course's own slug) would be a 500 that no test with a faked API sees.
test("a course opens from its card in the catalog", async ({ site }) => {
  await site.catalog.open();

  await site.catalog.openCourseCard("Ruby Course");

  await site.expectPath("/languages/ruby-en");
  await site.expectHeading("Ruby Course", 1);
});

test("a landing page slug of another locale is not found", async ({ site }) => {
  await site.expectNotFound("/languages/ruby-ru");
});

test("an archived landing page redirects to the course's main page", async ({ site }) => {
  await site.catalog.openCoursePage("/ru/languages/javascript-ru-archived");

  await site.expectPath("/ru/languages/javascript-ru");
  await site.expectHeading("Курс JavaScript", 1);
});
