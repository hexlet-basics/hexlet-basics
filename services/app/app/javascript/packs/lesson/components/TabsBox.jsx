import React from 'react';
import { useTranslation } from 'react-i18next';
import { Tab, Nav } from 'react-bootstrap';
import { useSelector, useDispatch } from 'react-redux';
import { actions } from '../slices/index.js';
import { currentTabValues } from '../utils/stateMachines.js';

import Editor from './Editor.jsx';
import Console from './Console.jsx';
import Solution from './Solution.jsx';

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

  const { editor, console, solution } = currentTabValues;

  return (
    <Tab.Container id="tabs" activeKey={currentTab} onSelect={changeTab}>
      <Nav variant="tabs" className="justify-content-center">
        <Nav.Item>
          <Nav.Link eventKey={editor} title={t(editor)}>{t(editor)}</Nav.Link>
        </Nav.Item>
        <Nav.Item>
          <Nav.Link eventKey={console} title={t(console)}>{t(console)}</Nav.Link>
        </Nav.Item>
        <Nav.Item>
          <Nav.Link eventKey={solution} title={t(solution)}>{t(solution)}</Nav.Link>
        </Nav.Item>
      </Nav>
      <Tab.Content bsPrefix="d-flex h-100 tab-content overflow-auto">
        <Tab.Pane eventKey={editor} bsPrefix="tab-pane h-100 pr-3 border-bottom w-100 overflow-hidden">
          <Editor />
        </Tab.Pane>
        <Tab.Pane eventKey={console} bsPrefix="tab-pane h-100 p-3 w-100">
          <Console />
        </Tab.Pane>
        <Tab.Pane eventKey={solution} bsPrefix="tab-pane h-100 p-3 w-100">
          <Solution />
        </Tab.Pane>
      </Tab.Content>
    </Tab.Container>
  );
};

export default TabsBox;
