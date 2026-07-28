import { setupWorker } from "msw/browser";

// Shared MSW worker for Vitest Browser Mode. It starts with NO handlers; each
// test declares what it needs via `worker.use(http.post(...))`, and the setup
// file resets handlers between tests. MSW intercepts the real fetch the generated
// hey-api client makes, so tests never touch the client's transport by hand.
export const worker = setupWorker();
