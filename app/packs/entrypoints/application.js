// @ts-check

// This file is automatically compiled by Webpack, along with any other files
// present in this directory. You're encouraged to place your actual application logic in
// a relevant structure within app/javascript and only use these pack files to reference
// that code so it'll be compiled.

import '../stylesheets/application.scss';
import 'bootstrap-icons/font/bootstrap-icons.css';

// eslint-disable-next-line import/no-extraneous-dependencies
import 'programming-languages-logos/src/javascript/javascript.svg';
// eslint-disable-next-line import/no-extraneous-dependencies
import 'programming-languages-logos/src/java/java.svg';
// eslint-disable-next-line import/no-extraneous-dependencies
import 'programming-languages-logos/src/php/php.svg';
// eslint-disable-next-line import/no-extraneous-dependencies
import 'programming-languages-logos/src/python/python.svg';
// eslint-disable-next-line import/no-extraneous-dependencies
import 'programming-languages-logos/src/ruby/ruby.svg';
// eslint-disable-next-line import/no-extraneous-dependencies
import 'programming-languages-logos/src/html/html.svg';
// eslint-disable-next-line import/no-extraneous-dependencies
import 'programming-languages-logos/src/css/css.svg';

import 'core-js/stable';
import 'regenerator-runtime/runtime';

import hljs from 'highlight.js';

require('@rails/ujs').start();
// require('@rails/activestorage').start();
// require('channels');
require('@popperjs/core');
require('bootstrap');

hljs.highlightAll();

const test = '';
console.log(test);

// Uncomment to copy all static images under ../images to the output folder and reference
// them with the image_pack_tag helper in views (e.g <%= image_pack_tag 'rails.png' %>)
// or the `imagePath` JavaScript helper below.
//
// const images = require.context('../images', true);
// const imagePath = (name) => images(name, true);
// console.log(imagePath);
