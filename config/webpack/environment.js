// @ts-check

const path = require('path');

module.exports = {
  externals: {
    // beacon: 'Beacon',
    // history: 'History',
    // links: 'Links',
    gon: 'gon',
  },
  // https://docs.rollbar.com/docs/unknown-script-error
  output: {
    crossOriginLoading: 'anonymous',
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, '../../app/javascript'),
    },
  },
  module: {
    rules: [
      {
        test: require.resolve('jquery'),
        loader: 'expose-loader',
        options: {
          exposes: 'jQuery',
        },
      },
    ],
  },
};
