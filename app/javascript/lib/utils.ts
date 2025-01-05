export function deviconClass(langName: string): string {
  const mapping: Record<string, string> = {
    css: "css3",
    html: "html5",
    cpp: "cplusplus",
    clang: "c",
    racket: "devicon",
    prolog: "devicon",
    fortran: "devicon",
  };
  const normalizedLangName = mapping[langName] ?? langName;
  return `devicon-${normalizedLangName}-plain`;
}

export function assetPath(filepath: string) {
  return `/vite-dev/images/${filepath}`;
}
