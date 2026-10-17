/**
 * The public password scripts/e2e_db.sh gives the fixture users below. The
 * fixtures themselves carry no password digests.
 */
export const password = "e2e-password";

/** A fixture user who is an admin. */
export const admin = { email: "admin-learner@example.com", name: "Admin Learner" };

/** A fixture user who is not (an editor on the staff, so the back office is open to them). */
export const learner = { email: "staff-editor@example.com", name: "Staff Editor" };

/** A fixture user with neither the admin flag nor a staff record. */
export const outsider = { email: "blank-user@example.com", name: "Blank User" };

/** An email nobody has registered, so journeys that sign up never collide. */
export function uniqueEmail(): string {
  return `e2e-${Date.now()}-${Math.random().toString(36).slice(2, 8)}@example.com`;
}
