import path from "node:path";
import react from "@vitejs/plugin-react";
import { defineConfig, loadEnv } from "vite";
import ViteRails from "vite-plugin-rails";

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), "VITE_");

  return {
    plugins: [
      react(),
      ViteRails({
        compress: mode !== "test",
      }),
    ],
    server: {
      allowedHosts: [env.VITE_APP_HOST],
      hmr: {
        clientPort: 443,
      },
    },
    ssr: {
      noExternal: [
        "lowlight",
        "highlight.js",
        "react-syntax-highlighter",
        "lodash",
        "use-inertia-form",
        "monaco-editor",
        "primereact",
        "react-timer-hook",
        "@monaco-editor/react",
      ], // Ensure it's handled correctly
    },
    resolve: {
      alias: {
        // "~bootstrap": path.resolve(__dirname, "node_modules/bootstrap"),
        "@": path.resolve(__dirname, "./app/javascript"),
      },
    },
  };
});
