process.env.NODE_ENV = process.env.NODE_ENV || 'production';

const environment = require('./environment');

environment.plugins.delete('Compression Brotli');
environment.plugins.delete('Compression');

module.exports = environment.toWebpackConfig();
