// @ts-check

import React from "react";
import { useSelector } from "react-redux";

import EntityContext from "../EntityContext.js";
import { neededPreview } from "../utils/languagesUtils.js";
import { currentTabValues } from "../utils/maps.js";
import ControlBox from "./ControlBox.jsx";
import HTMLPreview from "./HTMLPreview.jsx";
import TabsBox from "./TabsBox.jsx";

function App() {
  const { currentTab, content } = useSelector((state) => ({
    ...state.tabsBoxSlice,
    ...state.editorSlice,
  }));
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
    </div>
  );
}

export default App;
