// @ts-check

import React from 'react';
import { useSelector } from 'react-redux';

import TabsBox from './TabsBox.jsx';
import ControlBox from './ControlBox.jsx';
import GuideWidget from './GuideWidget.jsx';
import HTMLPreview from './HTMLPreview.jsx';
import { currentTabValues } from '../utils/maps.js';
import { neededPreview } from '../utils/languagesUtils.js';
import EntityContext from '../EntityContext.js';

const App = () => {
  const { currentTab, content } = useSelector((state) => ({ ...state.tabsBoxSlice, ...state.editorSlice }));
  const { language } = React.useContext(EntityContext);

  const renderHtmlPreview = () => {
    if (currentTab !== currentTabValues.editor) {
      return null;
    }
    if (!neededPreview(language)) {
      return null;
    }

    return <HTMLPreview html={content} />;
  };

  return (
    <div className="card vh-100 x-h-md-100">
      <TabsBox />
      {renderHtmlPreview()}
      <ControlBox />
      <GuideWidget />
    </div>
  );
};

export default App;
