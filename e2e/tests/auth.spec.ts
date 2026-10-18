import { learner, password, uniqueEmail } from "../support/people";
import { test } from "../support/test";

test("a visitor signs up, signs out and signs in again", async ({ site }) => {
  const email = uniqueEmail();
  const secret = "s3cret-pass-e2e";

  await site.session.signUp({ email, firstName: "Ada", password: secret });

  // Signing up also signs in: the header shows the new account.
  await site.session.expectSignedInAs("Ada");

  await site.session.signOut("Ada");
  await site.session.expectSignedOut();

  await site.session.signIn(email, secret);
  await site.session.expectSignedInAs("Ada");
});

test("a wrong password is refused and nobody is signed in", async ({ site }) => {
  await site.session.signIn(learner.email, `${password}-wrong`);

  await site.session.expectSignInRefused();
});

test("a fixture learner signs in with the seeded password", async ({ site }) => {
  await site.session.signInAs(learner);
  await site.expectPath("/");
});
