// @ts-check

import React from 'react';
import { useTranslation } from 'react-i18next';
import {
  Tab, Nav, ToggleButtonGroup, ToggleButton,
} from 'react-bootstrap';
import { useSelector, useDispatch } from 'react-redux';

import { actions } from '../slices/index.js';
import { currentTabValues } from '../utils/maps.js';

import EditorContainer from './EditorContainer';
import Output from './Output.jsx';
import Solution from './Solution.jsx';
import TestForExercise from './TestForExercise.jsx';
import { EDITORS } from '../slices/editorSlice.js';

const TabsBox = () => {
  const { t } = useTranslation();
  const dispatch = useDispatch();

  const currentTab = useSelector((state) => state.tabsBoxSlice.currentTab);
  const editorType = useSelector((state) => state.editorSlice.editorType);

  const changeTab = (newTabState) => {
    dispatch(actions.changeTab({ newTabState }));
  };

  const changeEditor = (editor) => {
    dispatch(actions.changeEditor({ editor }));
  };

  const {
    editor, output, solution, testForExercise,
  } = currentTabValues;

  // TODO: use anchor on load for choosing previuosly selected tab
  // useEffect(() => {
  //   const activeTab = window.location.hash.replace('#', '');
  // }, [null]);
  // console.log(activeTab)

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
        {currentTab === 'editor' && (
          <ToggleButtonGroup
            type="radio"
            name="options"
            value={editorType}
            onChange={changeEditor}
            size="sm"
            className="align-self-center"
          >
            <ToggleButton id="monaco" value={EDITORS.monaco} title={t('monaco')} variant="outline-primary">
              {t('monacoShort')}
            </ToggleButton>
            <ToggleButton id="codemirror" value={EDITORS.codemirror} title={t('codemirror')} variant="outline-primary">
              {t('codemirrorShort')}
            </ToggleButton>
          </ToggleButtonGroup>
        )}
      </Nav>
      <Tab.Content bsPrefix="d-flex h-100 tab-content overflow-auto">
        <Tab.Pane eventKey={editor} bsPrefix="tab-pane h-100 pe-3 w-100 overflow-hidden">
          <EditorContainer />
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

    </Tab.Container>
  );
};

export default TabsBox;
