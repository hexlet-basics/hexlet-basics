// @ts-check

import * as Sentry from '@sentry/browser';
import { BrowserTracing } from '@sentry/tracing';
import railsUjs from '@rails/ujs';

// import 'core-js/stable';
// import 'regenerator-runtime/runtime';
// require('@popperjs/core');

import 'bootstrap';
import hljs from 'highlight.js';

// https://docs.sentry.io/platforms/javascript/
Sentry.init({
  dsn: 'https://9aa813134b4a4cd2a5700ace36241a58@o1090356.ingest.sentry.io/6298540',
  // release: "my-project-name@2.3.12",
  integrations: [new BrowserTracing()],
  tracesSampleRate: 0.02,
});

railsUjs.start();
hljs.highlightAll();
