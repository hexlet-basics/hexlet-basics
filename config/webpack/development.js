const { merge } = require('@rails/webpacker');

process.env.NODE_ENV = process.env.NODE_ENV || 'development';

const webpackConfig = require('./base');

const devConfig = {
  target: 'web',
};

module.exports = merge(webpackConfig, devConfig);
