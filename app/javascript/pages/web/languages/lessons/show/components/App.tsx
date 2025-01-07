// @ts-check

import React from "react";
import { useSelector } from "react-redux";

import { neededPreview } from "@/lib/utils.js";
import { currentTabValues } from "@/lib/utils.js";
import ControlBox from "./ControlBox.js";
import HTMLPreview from "./HTMLPreview.js";
import TabsBox from "./TabsBox.js";
import { usePage } from "@inertiajs/react";
import type { Props } from "../types.js";

function App() {
  const {
    course,
  } = usePage<Props>().props;

  const { currentTab, content } = useSelector((state) => ({
    ...state.tabsBoxSlice,
    ...state.editorSlice,
  }));
  const renderHtmlPreview = () => {
    if (currentTab !== currentTabValues.editor) {
      return null;
    }
    if (!neededPreview(course)) {
      return null;
    }

    return <HTMLPreview html={content} />;
  };

  return (
    <>
      <TabsBox />
      {renderHtmlPreview()}
      <ControlBox />
    </>
  );
}

export default App;
