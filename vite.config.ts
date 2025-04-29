import path from "node:path";
import { sentryVitePlugin } from "@sentry/vite-plugin";
import react from "@vitejs/plugin-react";
import browserslist from "browserslist";
import { browserslistToTargets } from "lightningcss";
import { defineConfig, loadEnv } from "vite";
import ViteRails from "vite-plugin-rails";

export default defineConfig(({ mode, isSsrBuild }) => {
  const env = loadEnv(mode, process.cwd(), "VITE_");

  return {
    css: {
      lightningcss: {
        targets: browserslistToTargets(browserslist(">= 0.25%")),
      },
    },
    build: {
      sourcemap: "hidden",
      cssMinify: "lightningcss",
    },
    plugins: [
      react(),
      ViteRails({
        compress: false,
      }),
      sentryVitePlugin({
        url: env.VITE_SENTRY_URL,
        org: env.VITE_SENTRY_ORG,
        project: env.VITE_SENTRY_PROJECT,
        authToken: env.VITE_SENTRY_AUTH_TOKEN,
        release: {
          name: env.VITE_RELEASE_VERSION,
        },
        sourcemaps: {
          filesToDeleteAfterUpload: ["**/*.js.map"],
        },
        telemetry: false,
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
        "use-inertia-form",
        "lodash",
        "monaco-editor",
        "primereact",
        "react-timer-hook",
        "@monaco-editor/react",
        "react-use",
        "analytics",
        "@metro-fs/analytics-plugin-posthog",
      ], // Ensure it's handled correctly
    },
    resolve: {
      alias: {
        // "~bootstrap": path.resolve(__dirname, "node_modules/bootstrap"),
        // NOTE: код модуля monacoLoader.ts не должен выполнятся в ssr режиме
        "@/lib/monacoLoader.ts": isSsrBuild
          ? path.resolve(__dirname, "./app/javascript/lib/emptyModule.ts")
          : path.resolve(__dirname, "./app/javascript/lib/monacoLoader.ts"),
        "@": path.resolve(__dirname, "./app/javascript"),
      },
    },
  };
});
