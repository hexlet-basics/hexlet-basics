import type { PageProps } from '@inertiajs/core';
import type { router } from '@inertiajs/react';
import type { i18n } from 'i18next';
import type { SetupOptions } from 'node_modules/@inertiajs/react/types/createInertiaApp';
import type { ReactNode } from 'react';
import type { BackendEvent } from './events';
import type {
  LanguageCategory,
  LanguageLandingPageForLists,
  User,
} from './serializers';

export * from './serializers';

export type BreadcrumbItem = {
  name: string;
  url: string;
};

export type FlashKey = 'success' | 'notice' | 'error' | 'warning';

export type FlashVariant =
  | 'blue'
  | 'cyan'
  | 'dark'
  | 'grape'
  | 'gray'
  | 'green'
  | 'indigo'
  | 'lime'
  | 'orange'
  | 'pink'
  | 'red'
  | 'teal'
  | 'violet'
  | 'yellow';

export type FlashVariants = Record<FlashKey, FlashVariant>;

export type Locale = i18n['language'];

export interface SharedProps extends PageProps {
  auth: {
    user: User;
  };
  shouldAddContactMethod: boolean;
  railsDirectUploadsUrl: string;
  courseCategories: LanguageCategory[];
  suffix: 'ru' | null;
  locale: Locale;
  happendEvents: BackendEvent[] | null;
  // eventNames: BackendEvent[]
  flash: Partial<Record<FlashKey, string | null>>;
  landingPagesForLists: LanguageLandingPageForLists[];
  landingPagesForFooter: LanguageLandingPageForLists[];
  mobileBrowser: boolean;
  carrotQuestUserHash: string | null;
  metaTagsHTMLString: string;
}

export type RootProps = SetupOptions<HTMLElement, SharedProps>['props'];

// Temporary type definition, until @inertiajs/react provides one
export type ResolvedComponent = {
  default: ReactNode & {
    layout?: (page: ReactNode) => ReactNode;
  };
  layout?: (page: ReactNode) => ReactNode;
};

export type HttpRouterMethod = Extract<
  keyof typeof router,
  'get' | 'post' | 'put' | 'patch' | 'delete'
>;
