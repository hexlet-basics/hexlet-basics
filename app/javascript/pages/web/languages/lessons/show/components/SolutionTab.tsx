import { format } from "date-fns";
import React from "react";
import Countdown, { type CountdownRenderProps } from "react-countdown";
import { useTranslation } from "react-i18next";
import { github } from "react-syntax-highlighter/dist/esm/styles/hljs";

import waitingClock from "@/images/waiting_clock.png";
import { getEditorLanguage } from "@/lib/utils.ts";
import { usePage } from "@inertiajs/react";
import { Alert } from "react-bootstrap";
import SyntaxHighlighter from "react-syntax-highlighter";
import slice from "../slices/RootSlice.ts";
import { useAppDispatch, useAppSelector } from "../slices/index.ts";
import type { Props } from "../types.ts";

const waitingTime = 20 * 60 * 1000; // 20 min
// const waitingTime = 3000;

export default function SolutionTab() {
  const { course, lesson } = usePage<Props>().props;
  const content = useAppSelector((state) => state.content);
  const solutionState = useAppSelector((state) => state.solutionState);
  const startTime = useAppSelector((state) => state.startTime);
  const { t: tCommon } = useTranslation("common");
  const dispatch = useAppDispatch();

  const renderUserCode = () => {
    if (content === "") {
      return <Alert className="mt-3">{tCommon("userCodeInstructions")}</Alert>;
    }

    return (
      <>
        <SyntaxHighlighter
          style={github}
          showLineNumbers
          language={getEditorLanguage(course.slug!)}
        >
          {content}
        </SyntaxHighlighter>
      </>
    );
  };

  const renderSolution = () => {
    return (
      <div className="p-lg-3 hexlet-basics-content" id="basics-solution">
        <div className="mb-5">
          <h2 className="h3">{tCommon("teacherSolution")}</h2>
          <SyntaxHighlighter
            style={github}
            showLineNumbers
            language={getEditorLanguage(course.slug!)}
          >
            {lesson.original_code!}
          </SyntaxHighlighter>
        </div>
        <div>
          <h2 className="h3">{tCommon("userCode")}</h2>
          {renderUserCode()}
        </div>
      </div>
    );
  };

  const handleShowSolution = () => {
    dispatch(slice.actions.changeSolutionState("shown"));
  };

  const renderShowButton = () => (
    <>
      <p>{tCommon("solutionNotice")}</p>
      <div className="text-center">
        <button
          type="button"
          className="btn btn-secondary px-4"
          onClick={handleShowSolution}
        >
          {tCommon("showSolution")}
        </button>
      </div>
    </>
  );

  const renderContent = (countdownData: CountdownRenderProps) => {
    const { completed } = countdownData;

    if (completed || solutionState === "canBeShown") {
      return renderShowButton();
    }

    const remainingTime = format(new Date(countdownData.total), "mm:ss");

    return (
      <div className="text-center">
        <p className="lead">{tCommon("solutionInstructions")}</p>
        <div className="display-4">{remainingTime}</div>
        <img
          className="img-fluid px-5"
          src={waitingClock}
          alt="waiting_clock"
        />
      </div>
    );
  };

  return (
    <div>
      {solutionState === "shown" ? (
        renderSolution()
      ) : (
        <Countdown date={startTime + waitingTime} renderer={renderContent} />
      )}
    </div>
  );
}
