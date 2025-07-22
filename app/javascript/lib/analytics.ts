// import carrotquest from "@hexlet/analytics-plugin-carrotquest";

import yandexMetrika from '@hexlet/analytics-plugin-yandex-metrika';
import Analytics from 'analytics';
import type { BackendEvent } from '@/types/events';
import { log } from './utils';

// import googleAnalytics from "@analytics/google-analytics"

const plugins = [];

if (!import.meta.env.SSR && !import.meta.env.DEV) {
  console.log('!!!!');
  plugins.push(
    // googleAnalytics({
    //   measurementIds: ["G-XXXXXXXX"],
    // }),
    yandexMetrika({
      counterId: import.meta.env.VITE_YANDEX_METRIKA_COUNTER_ID,
      // enabled,
      mapEvents: {
        signed_up: 'ym-register',
        signed_in: 'ym-login',
        course_started: 'ym-start-course',
        lead_created: 'ym-submit-leadform',
        // purchase: 'ym-purchase',
        // begin_checkout: 'ym-begin-checkout',
        // add_to_wishlist: 'ym-add-to-wishlist',
        // ...добавляй любые кастомные события
      },
    }),
  );
}

/* Initialize analytics */
const analytics = Analytics({
  debug: import.meta.env.DEV,
  app: 'hexlet-basics',
  version: 100,
  plugins,
});

export function processHappendEvents(happendEvents: BackendEvent[]) {
  for (const happendEvent of happendEvents) {
    switch (happendEvent.type) {
      case 'UserSignedInEvent':
        analytics.identify(happendEvent.data.id.toString(), happendEvent.data);
        analytics.track('signed_in', happendEvent.data);
        break;
      case 'UserSignedUpEvent':
        analytics.identify(happendEvent.data.id.toString(), happendEvent.data);
        analytics.track('signed_up', happendEvent.data);
        break;
      case 'CourseStartedEvent':
        analytics.track('course_started', happendEvent.data);
        break;
      case 'LeadCreatedEvent': {
        const data = {
          ...happendEvent.data,
          phone_number: happendEvent.data.phone,
        };
        analytics.identify(happendEvent.data.user_id.toString(), data);
        analytics.track('lead_created', happendEvent.data);
        break;
      }
      case 'LessonStartedEvent':
        analytics.track('lesson_started', happendEvent.data);
        break;
      default:
        log('Unprocessed event: ', happendEvent);
        break;
    }
  }
}

export default analytics;
