// @ts-check

import React, { useContext } from 'react';
import { useTranslation } from 'react-i18next';
import { Highlight } from 'react-fast-highlight';

import EntityContext from '../EntityContext.js';
import { getLanguage } from '../utils/editorUtils.js';

const Test = () => {
  const { language, lessonVersion } = useContext(EntityContext);
  const { t } = useTranslation();

  return (
    <>
      <div className="p-lg-3 hexlet-basics-content">
        <h2 className="h3">{t('testInstructions')}</h2>
        <Highlight languages={[getLanguage(language)]}>
          {lessonVersion.test_code}
        </Highlight>
      </div>
    </>
  );
};

export default Test;
