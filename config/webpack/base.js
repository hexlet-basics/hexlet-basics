const { webpackConfig, merge } = require('@rails/webpacker');
const customConfig = require('./environment.js');

const config = merge(webpackConfig, customConfig);

module.exports = config;
