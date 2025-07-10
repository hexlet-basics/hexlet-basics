import type { Locale } from '@/types';
import { BackendEvent } from '@/types/events';

type GonData = {
  [key: string]: any;
  suffix: 'ru' | null;
  locale: Locale;
};

// @ts-expect-error not defined
if (!window.gon) {
  throw new Error('gon is not initialized');
}

// @ts-expect-error not defined
export const gon = window.gon as GonData;
