import type { User } from "@/client/types.gen";

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
