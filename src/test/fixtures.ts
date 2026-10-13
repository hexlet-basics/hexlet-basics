import type { User } from "@/client/types.gen";

// An administrator, for tests that open the admin pages: the admin layout's
// guard lets in a signed-in user who `canAccessAdmin`.
export const admin: User = {
  id: 1,
  firstName: "Ada",
  lastName: null,
  name: "Ada",
  email: "ada@example.com",
  admin: true,
  canAccessAdmin: true,
  assistantMessagesCount: 0,
  createdAt: "2026-01-01T00:00:00Z",
  createdAtAsTimestamp: null,
  type: "user",
};
