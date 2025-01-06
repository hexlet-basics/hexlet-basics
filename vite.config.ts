import path from "node:path";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";
import ViteRails from "vite-plugin-rails";

export default defineConfig({
  plugins: [react(), ViteRails()],
  resolve: {
    alias: {
      // "~bootstrap": path.resolve(__dirname, "node_modules/bootstrap"),
      "@": path.resolve(__dirname, "./app/javascript"),
    },
  },
});
