import type { GonData } from "@/lib/gon";
import type {
  FlashKey,
  LanguageCategory,
  LanguageLandingPageForLists,
  LeadCrud,
  Locale,
  User,
} from ".";
import type { DataLayer } from "./dataLayer";
import type { ApplicationEvent } from "./events";

declare global {
  interface Window {
    ymClientId?: string;
    dataLayer: DataLayer.Event[];
    gon: GonData;
  }
}

declare module "@inertiajs/core" {
  export interface InertiaConfig {
    flashDataType: {
      events?: ApplicationEvent[];
      flash: Partial<Record<FlashKey, string | null>>;
    };
  }
  export interface PageProps {
    auth: {
      user: User;
    };
    colorScheme: "dark" | "light";
    lead?: LeadCrud;
    host: string;
    locale: Locale;
    httpsHost: string;
    cloudflareTurnstileSiteKey: string;
    vkIdClientId: number;
    yandexIdClientId: string;
    formAuthToken: string;
    shouldAddContactMethod: boolean;
    railsDirectUploadsUrl: string;
    courseCategories: LanguageCategory[];
    suffix: "ru" | null;
    flash: Partial<Record<FlashKey, string | null>>;
    landingPagesForLists: LanguageLandingPageForLists[];
    landingPagesForFooter: LanguageLandingPageForLists[];
    mobileBrowser: boolean;
    metaTagsHTMLString: string;
  }
}
