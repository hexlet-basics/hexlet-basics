// @ts-check

import Typed from 'typed.js';
import gon from 'gon';

// eslint-disable-next-line no-unused-vars
const typed = new Typed('.hero-header', {
  strings: gon.languages_for_widget,
  typeSpeed: 50,
  backSpeed: 50,
  loop: true,
  smartBackspace: true, // Default value
});
