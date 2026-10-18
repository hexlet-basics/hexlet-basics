import type { CourseReadiness } from "@/client/types.gen";

// The Readiness (GLOSSARY, Course Readiness) of a Course and the moves between
// its values. A state is not an update: each move is an operation of its own on
// the server (mark as draft, in development, completed).
//
// Readiness describes how finished the Course is, not a workflow, so any change
// between two distinct values is legal. The table mirrors internal/transition
// on the server only to pick the buttons the admin screen offers; the server
// stays the judge and refuses repeating the current value with a 409.
const READINESS_TRANSITIONS = ["markDraft", "markInDevelopment", "markCompleted"] as const;
export type ReadinessTransition = (typeof READINESS_TRANSITIONS)[number];

export const readinessTarget: Record<ReadinessTransition, CourseReadiness> = {
  markCompleted: "completed",
  markDraft: "draft",
  markInDevelopment: "in_development",
};

// A Course with no readiness (nullable in the legacy schema) counts as a draft,
// as it does on the server.
export function transitionsFromReadiness(current: CourseReadiness | null): ReadinessTransition[] {
  return READINESS_TRANSITIONS.filter((name) => readinessTarget[name] !== (current ?? "draft"));
}
