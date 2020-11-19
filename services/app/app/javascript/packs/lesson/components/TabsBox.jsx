import React from 'react';
import { useTranslation } from 'react-i18next';
import cn from 'classnames';
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
  const { currentTab, checkInfo } = useSelector((state) => ({
    ...state.tabsBoxSlice,
    checkInfo: state.checkInfoSlice,
  }));
  const changeTab = (newTabState) => {
    dispatch(actions.changeTab({ newTabState }));
  };

  const { editor, console, solution } = currentTabValues;

  const badgeClassName = cn('badge mb-2 mb-sm-0 p-2', {
    'badge-success': checkInfo.passed,
    'badge-danger': !checkInfo.passed,
  });
  const headline = checkInfo.result ? t(`check.${checkInfo.result}.headline`) : null;

  const vdom = (
    <Tab.Container id="tabs" activeKey={currentTab} onSelect={changeTab}>
      <div className="d-flex flex-column flex-sm-row-reverse flex-shrink-0">
        <div className="my-auto d-none d-sm-block">
          {headline && <span className={badgeClassName}>{headline}</span>}
        </div>
        <div className="mr-auto flex-shrink-0">
          <Nav variant="tabs">
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
        </div>
      </div>
      <Tab.Content bsPrefix="d-flex h-100 tab-content overflow-auto">
        <Tab.Pane eventKey={editor} bsPrefix="tab-pane h-100 w-100 overflow-hidden">
          <Editor />
        </Tab.Pane>
        <Tab.Pane eventKey={console} bsPrefix="tab-pane h-100 w-100">
          <Console />
        </Tab.Pane>
        <Tab.Pane eventKey={solution} bsPrefix="tab-pane h-100 w-100">
          <Solution />
        </Tab.Pane>
      </Tab.Content>
    </Tab.Container>
  );
  return vdom;
};

export default TabsBox;
