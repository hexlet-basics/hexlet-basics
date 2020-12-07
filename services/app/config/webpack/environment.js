const { environment } = require('@rails/webpacker');
const webpack = require('webpack');
const MonacoWebpackPlugin = require('monaco-editor-webpack-plugin');

environment.plugins.prepend('Provide',
  new webpack.ProvidePlugin({
    $: 'jquery',
    jQuery: 'jquery',
    Popper: ['popper.js', 'default'],
  }));

environment.plugins.append('Monaco', new MonacoWebpackPlugin({
  languages: ['javascript', 'html', 'php', 'python', 'java', 'scheme', 'ruby', 'go'],
}));

environment.config.merge({
  externals: {
    gon: 'gon',
  },
});

module.exports = environment;
