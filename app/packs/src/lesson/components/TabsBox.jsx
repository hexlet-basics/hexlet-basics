// @ts-check

import React from 'react';
import { useTranslation } from 'react-i18next';
import { Tab, Nav, OverlayTrigger, Popover } from 'react-bootstrap';
import { useSelector, useDispatch } from 'react-redux';
import { actions } from '../slices/index.js';
import { currentTabValues } from '../utils/maps.js';

import Editor from './Editor.jsx';
import Output from './Output.jsx';
import Solution from './Solution.jsx';
import TestForExercise from './TestForExercise.jsx';

const TabsBox = () => {
  const { t } = useTranslation();
  const dispatch = useDispatch();

  const { currentTab } = useSelector((state) => ({
    ...state.tabsBoxSlice,
    checkInfo: state.checkInfoSlice,
  }));

  const changeTab = (newTabState) => {
    dispatch(actions.changeTab({ newTabState }));
  };

  const {
    editor, output, solution, testForExercise,
  } = currentTabValues;

  // TODO: use anchor on load for choosing previuosly selected tab
  // useEffect(() => {
  //   const activeTab = window.location.hash.replace('#', '');
  // }, [null]);
  // console.log(activeTab)

  const popover = (
    <Popover id="popover-basic">
      <Popover.Header as="h3">Popover right</Popover.Header>
      <Popover.Body>
        right?
      </Popover.Body>
    </Popover>
  );

  return (
    <Tab.Container id="tabs" activeKey={currentTab} onSelect={changeTab}>
      <Nav variant="tabs" className="justify-content-center small">
        {Object.values(currentTabValues).map((tabName) => (
          <Nav.Item key={tabName}>
            <Nav.Link
              href={`#${tabName}`}
              className="border-top-0 text-muted rounded-0"
              eventKey={tabName}
            >
              {t(tabName)}

            </Nav.Link>
          </Nav.Item>
        ))}
      </Nav>
      <Tab.Content bsPrefix="d-flex h-100 tab-content overflow-auto">
        <Tab.Pane eventKey={editor} bsPrefix="tab-pane h-100 pe-3 w-100 overflow-hidden">
          <Editor />
        </Tab.Pane>
        <Tab.Pane eventKey={output} bsPrefix="tab-pane h-100 p-3 w-100">
          <Output />
        </Tab.Pane>
        <Tab.Pane eventKey={testForExercise} bsPrefix="tab-pane h-100 p-3 w-100">
          <TestForExercise />
        </Tab.Pane>
        <Tab.Pane eventKey={solution} bsPrefix="tab-pane h-100 p-3 w-100">
          <Solution />
        </Tab.Pane>
      </Tab.Content>

      <OverlayTrigger trigger="click" placement="right" overlay={popover}>
        <i className="bi bi-question" />
      </OverlayTrigger>
    </Tab.Container>
  );
};

export default TabsBox;
