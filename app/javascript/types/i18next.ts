// import the original type declarations
import 'i18next';
import type locales from '../locales.json';

const defaultNS = 'web' as const;

declare module 'i18next' {
  interface i18n {
    // TODO: remove exclude es after adding es locale
    language: Exclude<keyof typeof locales, 'es'>;
  }
  // Extend CustomTypeOptions
  interface CustomTypeOptions {
    // custom namespace type, if you changed it
    defaultNS: typeof defaultNS;
    resources: (typeof locales)['ru'];
  }
}
