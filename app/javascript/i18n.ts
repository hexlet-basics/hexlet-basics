import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import locales from "./locales.json";
import jsLocales from "./locales/ru.js.json";

export const resources = Object.assign(locales, jsLocales);
export const defaultNS = "web";

i18n
	.use(initReactI18next) // passes i18n down to react-i18next
	.init({
		resources,
		defaultNS,
		ns: Object.keys(resources.ru),
		lng: "ru",
		interpolation: {
			escapeValue: false, // react already safes from xss => https://www.i18next.com/translation-function/interpolation#unescape
		},
	});
