// @ts-check

import path from 'path';
import { fileURLToPath } from 'url';
import MiniCssExtractPlugin from 'mini-css-extract-plugin';
import FontminPlugin from 'fontmin-webpack';
import { EsbuildPlugin } from 'esbuild-loader';
import { WebpackSweetEntry } from '@sect/webpack-sweet-entry';
// import { BundleAnalyzerPlugin } from 'webpack-bundle-analyzer';
import MonacoWebpackPlugin from 'monaco-editor-webpack-plugin';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const sourcePath = path.join(__dirname, 'app/javascript');
const monacoPath = path.join(__dirname, 'node_modules/monaco-editor');

export default {
  entry: WebpackSweetEntry(path.resolve(sourcePath, 'entrypoints/**/*.js'), 'js', 'entrypoints'),
  // stats: {
  //   // errorDetails: false,
  //   colors: true,
  //   errorStack: false,
  //   loggingTrace: false,
  // },
  devtool: 'cheap-source-map',
  externals: {
    jquery: 'jQuery',
    gon: 'gon',
    bootstrap: 'bootstrap',
  },
  resolve: {
    modules: ['node_modules', 'node_modules/devicon'],
    alias: {
      images: path.resolve(sourcePath, 'images'),
    },
  },
  output: {
    clean: {
      keep: /.keep/,
    },
    // eslint-disable-next-line max-len
    // https://srivishnu.totakura.in/2022/01/19/overcoming-challenges-with-jsbundling-rails-with-webpack-5-on-production.html#chunks-sourcemaps-and-asset-pipeline-sprockets
    filename: '[name].js',
    chunkFilename: '[name]-[contenthash].digested.js',
    sourceMapFilename: '[file]-[fullhash].digested.map',
    assetModuleFilename: 'images/[name]-[contenthash].digested[ext]',
    path: path.resolve(__dirname, 'app/assets/builds'),
    hashFunction: 'sha256',
    hashDigestLength: 64,
  },
  plugins: [
    new MiniCssExtractPlugin({}),
    // new BundleAnalyzerPlugin({ analyzerHost: '0.0.0.0' }),
    new FontminPlugin({
      autodetect: true, // automatically pull unicode characters from CSS
      // note: these settings are mutually exclusive and allowedFilesRegex has priority over skippedFilesRegex
      allowedFilesRegex: /icon/, // RegExp to only target specific fonts by their names
      skippedFilesRegex: null, // RegExp to skip specific fonts by their names
      textRegex: /\.css$/, // RegExp for searching text reference
    }),
    new MonacoWebpackPlugin({}),
  ],
  optimization: {
    minimizer: [
      new EsbuildPlugin({
        target: 'es2015',
        css: true,
      }),
    ],
    splitChunks: {
      // https://nozzlegear.com/blog/webpack-app-breaks-when-entrypoints-share-functions-or-components
      cacheGroups: {
        vendors: {
          test: /node_modules[\\/](highlight.js|axios|i18next|ansi_up|date-fns)[\\/]/,
          chunks: 'all',
          name: 'vendors',
        },
        routes: {
          test: /[\\/]routes.js/,
          name: 'routes',
          chunks: 'all',
        },
        reactEssentials: {
          test: /node_modules[\\/](react-dom|react)[\\/]/,
          name: 'reactEssentials',
          chunks: 'all',
        },
      },
    },
  },
  module: {
    rules: [
      {
        test: /\.css$/i,
        exclude: monacoPath,
        use: [
          MiniCssExtractPlugin.loader,
          'css-loader',
          'resolve-url-loader',
          {
            loader: 'postcss-loader',
            options: {
              postcssOptions: {
                plugins: [
                  'postcss-import', // NOTE: т.к. не работают нативные импорты
                  ['postcss-preset-env', {
                    features: {
                      'logical-properties-and-values': false, // NOTE: мы пока не используем rtl/ltr,
                      // а эта штука ломает вёрстку
                    },
                  }],
                ],
              },
            },
          },
        ],
      },
      {
        test: /\.css$/i,
        include: monacoPath,
        use: [
          'style-loader',
          'css-loader',
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
                  'postcss-import', // NOTE: т.к. не работают нативные импорты
                  ['postcss-preset-env', {
                    features: {
                      'logical-properties-and-values': false, // NOTE: мы пока не используем rtl/ltr,
                      // а эта штука ломает вёрстку
                    },
                  }],
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
