// import carrotquest from "@hexlet/analytics-plugin-carrotquest";
import yandexMetrika from "@hexlet/analytics-plugin-yandex-metrika";
import Analytics from "analytics";
// import googleAnalytics from "@analytics/google-analytics"

const enabled = !import.meta.env.SSR;

const plugins = [
  // googleAnalytics({
  //   measurementIds: ["G-XXXXXXXX"],
  // }),
  yandexMetrika({
    counterId: import.meta.env.VITE_YANDEX_METRIKA_COUNTER_ID,
    enabled,
    mapEvents: {
      signed_up: 'ym-register',
      signed_in: 'ym-login',
      course_started: 'ym-start-course',
      lead_created: 'ym-submit-leadform',
      // purchase: 'ym-purchase',
      // begin_checkout: 'ym-begin-checkout',
      // add_to_wishlist: 'ym-add-to-wishlist',
      // ...добавляй любые кастомные события
    }
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
