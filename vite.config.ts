import path from "node:path";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";
import ViteRails from "vite-plugin-rails";

export default defineConfig({
  plugins: [react(), ViteRails()],
  ssr: {
    noExternal: [
      "lowlight",
      "highlight.js",
      "react-syntax-highlighter",
      "lodash",
      "use-inertia-form",
      "monaco-editor",
      "primereact",
    ], // Ensure it's handled correctly
  },
  resolve: {
    alias: {
      // "~bootstrap": path.resolve(__dirname, "node_modules/bootstrap"),
      "@": path.resolve(__dirname, "./app/javascript"),
    },
  },
});
