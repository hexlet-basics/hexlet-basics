// A placeholder origin to resolve a redirect against. Search validation runs
// On the server as well as in the browser, so the page's own origin is not
// Always at hand — and it is not needed: what matters is only whether the value
// Stays on whatever origin it is resolved against.
const BASE = new URL("http://redirect.invalid");

// The same-site path a `redirect` search parameter may send a visitor to, or
// Undefined for anything that would leave the site.
//
// The value is parsed the way the browser will parse it rather than matched by
// Prefix: `/\evil.com`, `//evil.com` and `https://evil.com` all start
// Innocently enough for a string check to miss one, and the URL parser — the
// One navigation uses — resolves each to another host.
export function safeRedirectPath(value: string): string | undefined {
  if (!URL.canParse(value, BASE)) {
    return undefined;
  }
  const url = new URL(value, BASE);
  if (url.origin !== BASE.origin) {
    return undefined;
  }
  return `${url.pathname}${url.search}${url.hash}`;
}
