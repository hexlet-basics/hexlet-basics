// @ts-check

import { init } from '@sentry/browser';
import railsUjs from '@rails/ujs';

// import 'core-js/stable';
// import 'regenerator-runtime/runtime';
// require('@popperjs/core');

import 'bootstrap';
import hljs from 'highlight.js';

// TODO: https://github.com/getsentry/sentry-javascript/tree/master/packages/react
init({
  dsn: '__DSN__',
  // ...
});

railsUjs.start();
hljs.highlightAll();
