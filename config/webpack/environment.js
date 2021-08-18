// @ts-check

const path = require('path');
const browserslist = require('browserslist');
const fs = require('fs');

fs.writeFileSync('./browsers.json', JSON.stringify(browserslist()));

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
      vendor: path.resolve(__dirname, '../../app/packs/vendor'),
      lib: path.resolve(__dirname, '../../app/packs/lib'),
    },
  },
};
