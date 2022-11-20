import * as Sentry from '@sentry/browser';
import { BrowserTracing } from '@sentry/tracing';

// https://docs.sentry.io/platforms/javascript/
Sentry.init({
  dsn: 'https://9aa813134b4a4cd2a5700ace36241a58@o1090356.ingest.sentry.io/6298540',
  // release: "my-project-name@2.3.12",
  integrations: [new BrowserTracing()],
  tracesSampleRate: 0.02,
  ignoreErrors: [
    'top.GLOBALS', // Random plugins/extensions
    'VK is',
    'VK.Retargeting is ',
    "Can't find variable: VK",
    'pktAnnotationHighlighter',
    'Unexpected keyword',
    'illegal character',
    'Unexpected identifier',
    'Illegal invocation',
    'missing = in const declaration',
  ],
  allowUrls: [/https?:\/\/((cdn|cdn2|ru)\.)?#{configus.host}/],
  beforeSend(event, hint) {
    const stack = hint?.originalException?.stack || '';
    const errorInitiator = stack
      .split('\n')
      .map((line) => line.trim())
      .find((line) => line.startsWith('at'));
    const causedByConsole = errorInitiator
      ? errorInitiator.includes('<anonymous>:')
      : false;
    return causedByConsole ? null : event;
  },
});
