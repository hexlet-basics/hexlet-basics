declare global {
  interface Window {
    gon: {
      [key: string]: any;
      suffix: "ru" | null;
      locale: Locale;
    };
    dataLayer: Record<string, unknown>[];
  }
}

