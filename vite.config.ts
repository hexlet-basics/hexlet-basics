import path from "node:path";
import { sentryVitePlugin } from "@sentry/vite-plugin";
// import legacy from '@vitejs/plugin-legacy';
import react from "@vitejs/plugin-react";
import { visualizer } from "rollup-plugin-visualizer";
import { defineConfig, loadEnv, type PluginOption } from "vite";
import { patchCssModules } from "vite-css-modules";
import { imagetools } from "vite-imagetools";
import { beasties } from "vite-plugin-beasties";
import { ViteImageOptimizer } from "vite-plugin-image-optimizer";
import ViteRails from "vite-plugin-rails";

export default defineConfig(({ mode, isSsrBuild }) => {
  const env = loadEnv(mode, process.cwd(), "VITE_");

  return {
    // css: {
    //   transformer: 'lightningcss'
    // },
    build: {
      // sourcemap: "hidden",
      // sourcemap: false,
      // cssMinify: 'lightningcss',
      // rollupOptions: {
      //   output: {
      //     manualChunks: {
      //       react: ['react', 'react-dom', 'scheduler'],
      //       sentry: ['@sentry/react', '@sentry/tracing'],
      //       monaco: ['monaco-editor', '@monaco-editor/react'],
      //       inertia: ['@inertiajs/inertia', '@inertiajs/react', 'dayjs', 'i18next'],
      //       mantine: ['@mantine/core', '@mantine/hooks', '@mantine/modals'],
      //       // inertia: ['@inertiajs/inertia', '@inertiajs/inertia-react'],
      //       // vendor: ['lodash', 'axios'],
      //     },
      //   },
      // },
    },
    plugins: [
      imagetools(),
      visualizer() as PluginOption,
      beasties({
        // Plugin options
        options: {
          // Beasties library options
          preload: "swap",
          pruneSource: true, // Enable pruning CSS files
          inlineThreshold: 4000, // Inline stylesheets smaller than 4kb
        },
        // Filter to apply beasties only to specific HTML files
        filter: (path) => path.endsWith(".html"),
      }),
      ViteImageOptimizer({
        /* pass your config */
      }),
      // legacy({
      //   targets: ['defaults', 'not IE 11'],
      // }),
      patchCssModules({
        generateSourceTypes: true,
      }),
      react(),
      ViteRails({
        compress: false,
      }),
      sentryVitePlugin({
        disable: env.VITE_NODE_ENV !== "production",
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
