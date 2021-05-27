process.env.NODE_ENV = process.env.NODE_ENV || 'development'

// const webpackConfig = require('./base');
// const { merge } = require('@rails/webpacker');
// const WebpackAssetsManifest = require('webpack-assets-manifest');
const webpackConfig = require('./base.js');

// Remove after upgrade webpacker
delete webpackConfig.devServer.watchOptions
delete webpackConfig.devServer.stats
delete webpackConfig.devServer.overlay
delete webpackConfig.devServer.publicPath
delete webpackConfig.devServer.useLocalIp
delete webpackConfig.devServer.injectClient
delete webpackConfig.devServer.inline
delete webpackConfig.devServer.contentBase
delete webpackConfig.devServer.disableHostCheck
delete webpackConfig.devServer.quiet
delete webpackConfig.devServer.public
delete webpackConfig.devServer.clientLogLevel
webpackConfig.devServer.firewall = false
webpackConfig.devServer.client = {
  webSocketURL: {
    port: 80
  }
};

// https://github.com/rails/webpacker/issues/2864
// webpackConfig.plugins.splice(2, 1);
// webpackConfig.plugins.push(
//   new WebpackAssetsManifest({
//     enabled: true,
//     entrypoints: true,
//     writeToDisk: true,
//     output: 'manifest.json',
//     entrypointsUseAssets: true,
//     space: 2,
//     publicPath: true,
//   }),
// );
console.log(JSON.stringify(webpackConfig, '  ', 2));


module.exports = webpackConfig
