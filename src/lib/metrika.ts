// The Yandex Metrika client id. Legacy's layout asked the counter for it
// (`ym(counter, "getClientID")`) and left it on `window.ymClientId`, where the
// lead form picked it up. Null until the counter has answered, or wherever it
// is not loaded.
declare global {
  // `var` is how a global is declared, so `globalThis.ymClientId` is typed too.
  // eslint-disable-next-line no-var
  var ymClientId: string | undefined;
}

export function metrikaClientId(): string | null {
  if (!("window" in globalThis)) {
    return null;
  }
  return typeof globalThis.ymClientId === "string" && globalThis.ymClientId !== ""
    ? globalThis.ymClientId
    : null;
}
