import { Nav, Tab } from "react-bootstrap";
import { useTranslation } from "react-i18next";
import slice from "../slices/RootSlice.ts";

import {
  type AppState,
  useAppDispatch,
  useAppSelector,
} from "../slices/index.ts";
import EditorTab from "./EditorTab.tsx";
import OutputTab from "./OutputTab.tsx";
import SolutionTab from "./SolutionTab.tsx";
import TestsTab from "./TestsTab.tsx";

export default function TabsBox() {
  const { t } = useTranslation();
  const dispatch = useAppDispatch();

  const currentTab = useAppSelector((state) => state.currentTab);

  const changeTab = (newTabState: string | null) => {
    dispatch(slice.actions.changeTab(newTabState as AppState["currentTab"]));
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

      <Tab.Content className="x-h-md-100 overflow-hidden">
        <Tab.Pane className="x-h-md-100 vh-100" eventKey="editor">
          <EditorTab />
        </Tab.Pane>
        <Tab.Pane eventKey="output" className="overflow-auto h-100">
          <OutputTab />
        </Tab.Pane>
        <Tab.Pane eventKey="tests" className="overflow-auto h-100">
          <TestsTab />
        </Tab.Pane>
        <Tab.Pane eventKey="solution" className="overflow-auto h-100">
          <SolutionTab />
        </Tab.Pane>
      </Tab.Content>
    </Tab.Container>
  );
}
