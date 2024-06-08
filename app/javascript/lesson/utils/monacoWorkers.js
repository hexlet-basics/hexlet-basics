// @ts-check

// eslint-disable-next-line import/no-unresolved
import gon from 'gon';
// eslint-disable-next-line import/no-extraneous-dependencies, no-unused-vars
import * as monaco from 'monaco-editor';

const getWorkerPath = (workerFileName) => {
  const { monaco_worker_assets: monacoWorkerAssets = {}, assets_prefix: assetsPrefix = '.' } = gon;
  const workerFileNameWithHash = monacoWorkerAssets[workerFileName] || workerFileName;
  return `${assetsPrefix}/${workerFileNameWithHash}`;
};

window.MonacoEnvironment = {
  getWorkerUrl(_, label) {
    switch (label) {
      case 'json':
        return getWorkerPath('json.worker.js');
      case 'css':
      case 'scss':
      case 'less':
        return getWorkerPath('css.worker.js');
      case 'html':
      case 'handlebars':
      case 'razor':
        return getWorkerPath('html.worker.js');
      case 'typescript':
      case 'javascript':
        return getWorkerPath('ts.worker.js');
      default:
        return getWorkerPath('editor.worker.js');
    }
  },
};
