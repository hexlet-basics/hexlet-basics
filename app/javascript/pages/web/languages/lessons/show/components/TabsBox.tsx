import React from "react";
import { Nav, Tab } from "react-bootstrap";
import { useTranslation } from "react-i18next";
import slice, { type GeneralState } from "../slices/GeneralSlice.ts";

import Editor from "./EditorTab.tsx";
import Output from "./OutputTab.tsx";
import Solution from "./SolutionTab.tsx";
import TestsBox from "./TestsBox.tsx";
import { useAppDispatch, useAppSelector } from "../slices/index.ts";

function TabsBox() {
  const { t } = useTranslation();
  const dispatch = useAppDispatch();

  const currentTab = useAppSelector((state) => state.currentTab);

  const changeTab = (newTabState: string | null) => {
    dispatch(
      slice.actions.changeTab(newTabState as GeneralState["currentTab"]),
    );
  };

  // TODO: use anchor on load for choosing previuosly selected tab
  // useEffect(() => {
  //   const activeTab = window.location.hash.replace('#', '');
  // }, [null]);
  // console.log(activeTab)

  return (
    <Tab.Container
      id="tabs"
      activeKey={currentTab}
      defaultActiveKey="editor"
      onSelect={changeTab}
    >
      <Nav variant="underline" fill justify className="mb-3 small">
        <Nav.Item>
          <Nav.Link className="link-body-emphasis" eventKey="editor">
            {t("languages.lessons.show.editor")}
          </Nav.Link>
        </Nav.Item>
        <Nav.Item>
          <Nav.Link className="link-body-emphasis" eventKey="output">
            {t("languages.lessons.show.output")}
          </Nav.Link>
        </Nav.Item>
        <Nav.Item>
          <Nav.Link className="link-body-emphasis" eventKey="tests">
            {t("languages.lessons.show.tests")}
          </Nav.Link>
        </Nav.Item>
        <Nav.Item>
          <Nav.Link className="link-body-emphasis" eventKey="solution">
            {t("languages.lessons.show.solution")}
          </Nav.Link>
        </Nav.Item>
      </Nav>

      <Tab.Content className="h-100 overflow-hidden">
        <Tab.Pane className="h-100" eventKey="editor">
          <Editor />
        </Tab.Pane>
        <Tab.Pane eventKey="output" className="overflow-auto h-100">
          <Output />
        </Tab.Pane>
        <Tab.Pane eventKey="tests" className="overflow-auto h-100">
          <TestsBox />
        </Tab.Pane>
        <Tab.Pane eventKey="solution" className="overflow-auto h-100">
          <Solution />
        </Tab.Pane>
      </Tab.Content>
    </Tab.Container>
  );
}

export default TabsBox;
