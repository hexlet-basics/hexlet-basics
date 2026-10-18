import type { Course, PublicUser, User } from "@/client/types.gen";

// An administrator, for tests that open the admin pages: the admin layout's
// guard lets in a signed-in user who `canAccessAdmin`.
export const admin: User = {
  admin: true,
  assistantMessagesCount: 0,
  canAccessAdmin: true,
  createdAt: "2026-01-01T00:00:00Z",
  createdAtAsTimestamp: null,
  email: "ada@example.com",
  firstName: "Ada",
  id: 1,
  lastName: null,
  name: "Ada",
  type: "user",
};

// The account embedded in a record that names its author, as the API reads it.
export const eva: PublicUser = { firstName: "Eva", id: 1, lastName: null, name: "Eva" };

// A course, for the admin screens that embed or list one.
export const ruby: Course = {
  categoryId: null,
  coverListVariant: null,
  coverThumbVariant: null,
  createdAt: "2026-01-01T00:00:00Z",
  currentVersion: null,
  currentVersionId: null,
  enrollmentsCount: 0,
  hexletProgramLandingPage: null,
  id: 7,
  learnAs: "first_language",
  lessonsCount: 3,
  name: null,
  ratingCount: 0,
  ratingValue: 0,
  readiness: "draft",
  repositoryUrl: null,
  slug: "ruby",
};
