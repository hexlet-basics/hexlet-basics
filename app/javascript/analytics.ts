import carrotquest from "@hexlet/analytics-plugin-carrotquest";
import postHog from "@metro-fs/analytics-plugin-posthog";
import Analytics from "analytics";
import { mapValues } from "es-toolkit";
// import googleAnalytics from "@analytics/google-analytics"

const carrotQuestEventsMapping = {
  signed_in: "$authorized",
  signed_up: "$registered",
  course_started: "Начал курс",
  course_finished: "Завершил курс",
  lesson_started: "Начал урок",
  lesson_finished: "Завершил урок",
} as const;

export const carrotQuestEvents = mapValues(
  carrotQuestEventsMapping,
  (_, k) => k,
);

const carrotQuestEventPropsMapping = {
  id: "id",
  email: "email",
  locale: "locale",
  course_slug: "course_slug",
  lesson_slug: "lesson_slug",
  slug: "slug",
};

/* Initialize analytics */
const analytics = Analytics({
  debug: true,
  app: "hexlet-basics",
  version: 100,
  plugins: [
    // googleAnalytics({
    //   measurementIds: ["G-XXXXXXXX"],
    // }),
    postHog({
      token: import.meta.env.VITE_REACT_APP_PUBLIC_POSTHOG_KEY,
      enabled: !import.meta.env.SSR,
      options: {
        api_host: import.meta.env.VITE_REACT_APP_PUBLIC_POSTHOG_HOST,
        debug: import.meta.env.DEV,
        disable_session_recording: true,
        autocapture: false,
        capture_pageleave: true,
        disable_surveys: true,
        custom_campaign_params: [
          "promo_name",
          "promo_position",
          "promo_type",
          "promo_creative",
          "promo_start",
        ],
      },
    }),
    carrotquest({
      apiKey: import.meta.env.VITE_CARROTQUEST_API_KEY,
      eventsMapping: carrotQuestEventsMapping,
      // Возможно, стоит сохранять все свойства
      eventPropsMapping: carrotQuestEventPropsMapping,
      userPropsMapping: {},
    }),
  ],
});

export default analytics;
