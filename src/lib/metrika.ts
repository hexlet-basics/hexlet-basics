// The Yandex Metrika client id. Legacy's layout asked the counter for it
// (`ym(counter, "getClientID")`) and left it on `window.ymClientId`, where the
// Lead form picked it up. Null until the counter has answered, or wherever it
// Is not loaded.
declare global {
  interface Window {
    ymClientId?: string;
  }
}

export function metrikaClientId(): string | null {
  if (typeof window === "undefined") {
    return null;
  }
  return typeof window.ymClientId === "string" && window.ymClientId !== ""
    ? window.ymClientId
    : null;
}
