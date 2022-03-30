// @ts-check

import path from 'path';
import { fileURLToPath } from 'url';
import webpack from 'webpack';
import { ESBuildMinifyPlugin } from 'esbuild-loader';
import MiniCssExtractPlugin from 'mini-css-extract-plugin';
// import { WebpackManifestPlugin } from 'webpack-manifest-plugin';

const __dirname = path.dirname(fileURLToPath(import.meta.url));

const getPath = (filename, type = 'javascript') => path.resolve(__dirname, `app/${type}/${filename}`);

export default {
  entry: {
    application: [getPath('application.js'), getPath('application.scss', 'assets/stylesheets')],
    welcomePage: getPath('welcome-page.js'),
    lessonPage: getPath('lesson-page.js'),
    amp: getPath('amp.css', 'assets/stylesheets'),
  },
  devtool: false,
  externals: {
    jquery: 'jQuery',
    gon: 'gon',
  },
  plugins: [
    // new WebpackManifestPlugin({
    //   useEntryKeys: true,
    //   writeToFileEmit: true,
    // }),
    new webpack.SourceMapDevToolPlugin({
      filename: '[file].map[query]',
      exclude: 'vendors',
    }),
    new MiniCssExtractPlugin({
      filename: '[name].css',
    }),
  ],
  output: {
    filename: '[name].js',
    // chunkFilename: '[name].chunk.js',
    // sourceMapFilename: '[name].js.map',
    path: path.resolve(__dirname, 'app/assets/builds'),
  },
  optimization: {
    splitChunks: {
      cacheGroups: {
        commons: {
          test: /[\\/]node_modules[\\/]/,
          name: 'vendors',
          chunks: 'initial',
        },
      },
    },
    minimizer: [
      new ESBuildMinifyPlugin({
        target: 'es2015',
        css: true,
      }),
    ],
  },
  module: {
    rules: [
      {
        test: /\.css$/i,
        use: [
          MiniCssExtractPlugin.loader,
          'css-loader',
          {
            loader: 'postcss-loader',
            options: {
              postcssOptions: {
                plugins: [
                  ['postcss-preset-env', { }],
                ],
              },
            },
          },
        ],
      },
      {
        test: /\.s[ac]ss$/i,
        use: [
          MiniCssExtractPlugin.loader,
          'css-loader',
          {
            loader: 'postcss-loader',
            options: {
              postcssOptions: {
                plugins: [
                  ['postcss-preset-env', { }],
                ],
              },
            },
          },
          'resolve-url-loader',
          {
            loader: 'sass-loader',
            options: {
              sourceMap: true, // <-- !!IMPORTANT!!
            },
          },
        ],
      },
      {
        test: /\.(png|svg|jpg|jpeg|gif)$/i,
        type: 'asset/resource',
      },
      {
        test: /\.(woff|woff2|eot|ttf|otf)$/i,
        type: 'asset/resource',
      },
      {
        test: /\.jsx?$/,
        loader: 'esbuild-loader',
        options: {
          loader: 'jsx',
          target: 'es2015',
        },
      },
    ],
  },
};
