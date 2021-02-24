const { webpackConfig, merge } = require('@rails/webpacker');
const CompressionPlugin = require('compression-webpack-plugin');
const WebpackAssetsManifest = require('webpack-assets-manifest');

const customConfig = require('./environment.js');

// https://github.com/rails/webpacker/issues/2864
webpackConfig.plugins.splice(2, 1);
webpackConfig.plugins.push(
  new WebpackAssetsManifest({
    enabled: true,
    entrypoints: true,
    writeToDisk: true,
    output: 'manifest.json',
    entrypointsUseAssets: true,
    space: 2,
    publicPath: true,
  }),
);

// https://github.com/rails/webpacker/issues/2399
webpackConfig.plugins = webpackConfig.plugins.filter((p) => !(p instanceof CompressionPlugin));

module.exports = merge(webpackConfig, customConfig);
