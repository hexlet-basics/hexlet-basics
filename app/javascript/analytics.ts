import Analytics from "analytics";
import postHog from '@metro-fs/analytics-plugin-posthog';
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
      enabled: true,
      options: {
        api_host: import.meta.env.VITE_REACT_APP_PUBLIC_POSTHOG_HOST,
        debug: process.env.NODE_ENV === 'development',
        disable_session_recording: true,
        autocapture: false,
        capture_pageleave: false,
      },
    }),
  ],
});

export default analytics;
