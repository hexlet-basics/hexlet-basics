import React from 'react';
import { useSelector } from 'react-redux';
import TabsBox from './TabsBox.jsx';
import ControlBox from './ControlBox.jsx';
import HTMLPreview from './HTMLPreview';
import { tabsBoxSliceName, editorSliceName } from '../slices/index.js';
import { currentTabValues } from '../utils/stateMachines.js';
import { neededPreview } from '../utils/languagesUtils.js';
import EntityContext from '../EntityContext.js';

const App = () => {
  const { currentTab, content } = useSelector((state) => ({ ...state[tabsBoxSliceName], ...state[editorSliceName] }));
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

  const vdom = (
    <>
      <TabsBox />
      {renderHtmlPreview()}
      <ControlBox />
    </>
  );
  return vdom;
};

export default App;
