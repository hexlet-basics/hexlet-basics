import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";
import ViteRails from "vite-plugin-rails";

export default defineConfig({
	plugins: [
		react(),
		ViteRails(),
	],
});
