const { webpackConfig, merge } = require('@rails/webpacker')
const customConfig = require('./environment.js');

// https://github.com/rails/webpacker/issues/2844
// needed for bootstrap-icons
webpackConfig.module.rules.forEach((module) => { // lets iterate Rules...
  const type = module.test.toString();
  if (module.test && type.includes('scss')) { // scss and sass are in same bag.
    module.use.splice(-1, 0, {
      loader: require.resolve('resolve-url-loader')
    });
    module.use.forEach((info) => {
      if (info.loader.includes('sass')) {
        info.options.sourceMap = true;
      }
    });
  }
});

const config = merge(webpackConfig, customConfig);
module.exports = config;
