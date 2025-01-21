export const images = import.meta.glob<ImageModule>("./images/**/*", {
  eager: true,
});

export function getImageUrl(name: string): string {
  const mod = images[`./images/${name}`];
  if (!mod) throw new Error(`No image available for ${name}`);
  return mod.default;
}
