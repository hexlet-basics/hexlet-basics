import React from "react";
import { useSelector } from "react-redux";

import { usePage } from "@inertiajs/react";

import { neededPreview } from "@/lib/utils.ts";
import _ from "lodash";
import { useAppSelector } from "../slices/index.ts";
import type { Props } from "../types.ts";
import ControlBox from "./ControlBox.tsx";
import HTMLPreview from "./HTMLPreview.tsx";
import TabsBox from "./TabsBox.tsx";

function App() {
  const { course } = usePage<Props>().props;

  const content = useAppSelector((state) => state.content);
  const currentTab = useAppSelector((state) => state.currentTab);

  const renderHtmlPreview = () => {
    if (currentTab !== "editor") {
      return null;
    }
    if (!neededPreview(course.slug!)) {
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
