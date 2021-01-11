// This file is automatically compiled by Webpack, along with any other files
// present in this directory. You're encouraged to place your actual application logic in
// a relevant structure within app/javascript and only use these pack files to reference
// that code so it'll be compiled.
import '../stylesheets/application.scss'
import 'core-js/stable';
import 'regenerator-runtime/runtime';

import '@fortawesome/fontawesome-free/js/all.js';
import hljs from 'highlight.js';

require('jquery');
require('@rails/ujs').start();
// require('@rails/activestorage').start();
// require('channels');
require('popper.js');
require('bootstrap');

hljs.initHighlightingOnLoad();

// Uncomment to copy all static images under ../images to the output folder and reference
// them with the image_pack_tag helper in views (e.g <%= image_pack_tag 'rails.png' %>)
// or the `imagePath` JavaScript helper below.
//
// const images = require.context('../images', true)
// const imagePath = (name) => images(name, true)
