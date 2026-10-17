import { admin, learner, outsider } from "../support/people";
import { test } from "../support/test";

test("a guest is sent to sign in and returns to the admin afterwards", async ({ site }) => {
  await site.admin.open();

  await site.expectPath("/session/new?redirect=%2Fadmin");
  await site.expectHeading("Sign In", 1);
});

test("a signed-in user who is neither an admin nor staff is refused", async ({ site }) => {
  await site.session.signInAs(outsider);

  await site.admin.open();

  await site.expectPath("/");
  await site.session.expectSignedInAs(outsider.name);
});

test("an admin sees the back office and its lists", async ({ site }) => {
  await site.session.signInAs(admin);

  await site.admin.open();
  await site.expectPath("/admin/courses");
  await site.admin.expectCourseList();

  await site.admin.openUsers();
  await site.expectPath("/admin/users");
  await site.expectHeading("Users");
  await site.admin.expectUserListed(admin.email);
  await site.admin.expectUserListed(learner.email);
});
