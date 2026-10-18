import { learner } from "../support/people";
import { test } from "../support/test";

// The player titles the lesson with the course's landing page in the request
// locale; the fixtures give the JavaScript course no English page but "Frontend".
const lesson = "/languages/javascript/lessons/hello-world";

test("a guest reads a lesson, sees the other lessons as links and is invited to sign in", async ({
  site,
}) => {
  await site.lesson.open(lesson);

  await site.expectHeading("Frontend: Hello, World!", 1);
  await site.lesson.expectLessonTab();
  await site.lesson.expectCodeEditor();

  await site.lesson.openNavigation();
  await site.lesson.expectLinkTo("Variables");

  await site.lesson.openAssistant();
  await site.lesson.expectSignInInvitation();
});

test("a signed-in learner opens the same lesson and can use the assistant", async ({ site }) => {
  await site.session.signInAs(learner);
  await site.lesson.open(lesson);

  await site.expectHeading("Frontend: Hello, World!", 1);

  await site.lesson.openAssistant();
  await site.lesson.expectNoSignInInvitation();
  await site.lesson.expectAssistantGreeting();

  await site.lesson.expectCannotSend();
  await site.lesson.typeQuestion("How do I print a line?");
  await site.lesson.expectCanSend();
});
