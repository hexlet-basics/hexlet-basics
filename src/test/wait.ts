// Resolves after `ms` milliseconds: for a test that has to let time pass, as when
// it asserts that a request which would have gone out, did not.
export function sleep(ms: number): Promise<null> {
  const { promise, resolve } = Promise.withResolvers<null>();
  globalThis.setTimeout(() => resolve(null), ms);
  return promise;
}
