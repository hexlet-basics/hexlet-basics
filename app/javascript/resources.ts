export const resources = import.meta.glob<ResourceModule>('./images/**/*', {
  eager: true,
});

export function getResourceUrl(name: string): string {
  const mod = resources[`./images/${name}`];
  if (!mod) throw new Error(`No resource available for ${name}`);
  return mod.default;
}
