import carrotquest from "@hexlet/analytics-plugin-carrotquest";
import yandexMetrika from "@hexlet/analytics-plugin-yandex-metrika";
import postHog from "@metro-fs/analytics-plugin-posthog";
import Analytics from "analytics";
// import googleAnalytics from "@analytics/google-analytics"

const enabled = !import.meta.env.SSR && import.meta.env.PROD;

const plugins = [
  // googleAnalytics({
  //   measurementIds: ["G-XXXXXXXX"],
  // }),
  yandexMetrika({
    counterId: import.meta.env.VITE_YANDEX_METRIKA_COUNTER_ID,
    enabled,
  }),
  postHog({
    token: import.meta.env.VITE_REACT_APP_PUBLIC_POSTHOG_KEY,
    enabled,
    options: {
      api_host: import.meta.env.VITE_REACT_APP_PUBLIC_POSTHOG_HOST,
      debug: import.meta.env.DEV,
      disable_session_recording: true,
      autocapture: false,
      capture_pageview: false,
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
    enabled,
    eventsMapping: {
      signed_in: "$authorized",
      signed_up: "$registered",
      course_started: "Начал курс",
      course_finished: "Завершил курс",
      lesson_started: "Начал урок",
      lesson_finished: "Завершил урок",
    },
    // Возможно, стоит сохранять все свойства
    eventPropsMapping: {
      id: "id",
      email: "email",
      locale: "locale",
      course_slug: "course_slug",
      lesson_slug: "lesson_slug",
      slug: "slug",
    },
    userPropsMapping: {},
  }),
];

/* Initialize analytics */
const analytics = Analytics({
  debug: import.meta.env.DEV,
  app: "hexlet-basics",
  version: 100,
  plugins,
});

export default analytics;
