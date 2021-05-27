// @ts-check

import React, { useContext } from 'react';
import { useTranslation } from 'react-i18next';
import Highlight from 'react-highlight';

import EntityContext from '../EntityContext.js';
import { getLanguage } from '../utils/editorUtils.js';

const Test = () => {
  const { language, lessonVersion } = useContext(EntityContext);
  const { t } = useTranslation();

  return (
    <div>
      <p className="text-center lead">{t('testInstructions')}</p>
      <div className="hexlet-basics-content">
        <Highlight className={getLanguage(language)}>
          {lessonVersion.test_code}
        </Highlight>
      </div>
    </div>
  );
};

export default Test;
