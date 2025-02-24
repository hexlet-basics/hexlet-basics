import postHog from "@metro-fs/analytics-plugin-posthog";
import Analytics from "analytics";
// import googleAnalytics from '@analytics/google-analytics'

/* Initialize analytics */
const analytics = Analytics({
  debug: true,
  app: "hexlet-basics",
  version: 100,
  plugins: [
    // googleAnalytics({
    //   measurementIds: ['G-XXXXXXXX'],
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
  ],
});

export default analytics;
