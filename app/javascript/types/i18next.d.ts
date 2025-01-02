// import the original type declarations
import "i18next";
import type { defaultNS, resources } from "../i18n.ts";

declare module "i18next" {
	interface i18n {
		language: "en" | "ru";
	}
	// Extend CustomTypeOptions
	interface CustomTypeOptions {
		// custom namespace type, if you changed it
		defaultNS: typeof defaultNS;
		resources: (typeof resources)["ru"];
	}
}
