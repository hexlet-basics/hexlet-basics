// The Publication State (GLOSSARY) of blog posts, banners, reviews and landing
// pages, and the moves between its states. A state is not an update: each move
// is an operation of its own on the server (publish, unpublish, archive).
//
// The table mirrors internal/transition on the server, and only to pick the
// buttons the admin screens offer. The server stays the judge and refuses an
// illegal move with a 409.
export type PublicationState = "draft" | "published" | "archived";
const STATE_TRANSITIONS = ["publish", "unpublish", "archive"] as const;
export type StateTransition = (typeof STATE_TRANSITIONS)[number];

const transitionSources: Record<StateTransition, PublicationState[]> = {
  archive: ["draft", "published"],
  publish: ["draft", "archived"],
  unpublish: ["published"],
};

export function transitionsFrom(state: PublicationState): StateTransition[] {
  return STATE_TRANSITIONS.filter((name) => transitionSources[name].includes(state));
}
