// The Publication State (GLOSSARY) of blog posts, banners, reviews and landing
// Pages, and the moves between its states. A state is not an update: each move
// Is an operation of its own on the server (publish, unpublish, archive).
//
// The table mirrors internal/transition on the server, and only to pick the
// Buttons the admin screens offer. The server stays the judge and refuses an
// Illegal move with a 409.
export type PublicationState = "draft" | "published" | "archived";
export type StateTransition = "publish" | "unpublish" | "archive";

const transitionSources: Record<StateTransition, PublicationState[]> = {
  publish: ["draft", "archived"],
  unpublish: ["published"],
  archive: ["draft", "published"],
};

export function transitionsFrom(state: PublicationState): StateTransition[] {
  return (Object.keys(transitionSources) as StateTransition[]).filter((name) =>
    transitionSources[name].includes(state),
  );
}
