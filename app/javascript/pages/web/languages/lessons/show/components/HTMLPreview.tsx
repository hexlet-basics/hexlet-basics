// @ts-check

import React from "react";
import Frame from "react-frame-component";

function HTMLPreview({ html }) {
  return (
    <div className="pt-3 px-2 pb-2 h-50 border-top">
      <Frame className="border-0 h-100 w-100">
        <div
          className="text-black"
          dangerouslySetInnerHTML={{ __html: html }}
        />
      </Frame>
    </div>
  );
}

export default HTMLPreview;
