import { visualizer } from "rollup-plugin-visualizer";
import { ViteImageOptimizer } from 'vite-plugin-image-optimizer';
import legacy from '@vitejs/plugin-legacy'
import path from "node:path";
import { sentryVitePlugin } from "@sentry/vite-plugin";
import react from "@vitejs/plugin-react";
import browserslist from "browserslist";
import { browserslistToTargets } from "lightningcss";
import { defineConfig, loadEnv, PluginOption } from "vite";
import ViteRails from "vite-plugin-rails";
import { patchCssModules } from 'vite-css-modules'

export default defineConfig(({ mode, isSsrBuild }) => {
  const env = loadEnv(mode, process.cwd(), "VITE_");

  return {
    css: {
      lightningcss: {
        targets: browserslistToTargets(browserslist(">= 0.25%")),
      },
    },
    build: {
      // sourcemap: "hidden",
      // sourcemap: false,
      cssMinify: "lightningcss",
    },
    plugins: [
      visualizer() as PluginOption,
      ViteImageOptimizer({
        /* pass your config */
      }),
      legacy({
        targets: ['defaults', 'not IE 11'],
      }),
      patchCssModules({
        generateSourceTypes: true
      }),
      react(),
      ViteRails({
        compress: false,
      }),
      sentryVitePlugin({
        disable: env.NODE_ENV != "production",
        applicationKey: env.VITE_APP_HOST,
        url: env.VITE_SENTRY_URL,
        org: env.VITE_SENTRY_ORG,
        project: env.VITE_SENTRY_PROJECT,
        authToken: env.VITE_SENTRY_AUTH_TOKEN,
        release: {
          name: env.VITE_RELEASE_VERSION,
          setCommits: {
            auto: true,
            ignoreEmpty: true,
          },
        },
        sourcemaps: {
          filesToDeleteAfterUpload: ["**/*.js.map"],
          // rewriteSources: (source) => source.replace(/^(\.\.\/|\.\/)+/, ''),
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
        "monaco-editor",
        "react-timer-hook",
        "@monaco-editor/react",
        "analytics",
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
