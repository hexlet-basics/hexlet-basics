const { environment } = require('@rails/webpacker');
const webpack = require('webpack');
const MonacoWebpackPlugin = require('monaco-editor-webpack-plugin');

environment.plugins.prepend('Provide',
  new webpack.ProvidePlugin({
    $: 'jquery',
    jQuery: 'jquery',
    Popper: ['popper.js', 'default'],
  }));

environment.plugins.prepend('Provide', new MonacoWebpackPlugin({
  languages: ['javascript', 'css', 'html', 'php', 'python', 'java', 'scheme', 'ruby', 'go'],
}));

environment.config.merge({
  externals: {
    gon: 'gon',
  },
});

module.exports = environment;
