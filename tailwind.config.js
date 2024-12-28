import daisyui from "daisyui";

/** @type {import('tailwindcss').Config} */
export default {
	content: ["./app/javascript/**/*.{js,ts,jsx,tsx}"],
	theme: {
		extend: {},
	},
	plugins: [
    daisyui,
  ],
};
