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
  ],
});

export default analytics;
