import { setupWorker } from "msw/browser";

// Shared MSW worker for Vitest Browser Mode. It starts with NO handlers; each
// Test declares what it needs via `worker.use(http.post(...))`, and the setup
// File resets handlers between tests. MSW intercepts the real fetch the generated
// Hey-api client makes, so tests never touch the client's transport by hand.
export const worker = setupWorker();
