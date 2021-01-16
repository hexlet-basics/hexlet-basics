const { webpackConfig, merge } = require('@rails/webpacker');
const CompressionPlugin = require('compression-webpack-plugin');

const customConfig = require('./environment.js');

// https://github.com/rails/webpacker/issues/2399
webpackConfig.plugins = webpackConfig.plugins.filter((p) => !(p instanceof CompressionPlugin));

module.exports = merge(webpackConfig, customConfig);
