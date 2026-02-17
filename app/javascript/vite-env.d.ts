/// <reference types="vite/client" />

type ResourceModule = {
  default: string;
};

// type ViteTypeOptions = {};

// interface ImportMetaEnv {
//   readonly VITE_APP_HOST: string;
//   // more env variables...
// }

// interface ImportMeta {
//   readonly env: ImportMetaEnv;
// }

declare module "*?format=webp" {
  const src: string;
  export default src;
}
