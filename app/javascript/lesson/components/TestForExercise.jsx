// @ts-check

import React, { useContext } from 'react';
import { useTranslation } from 'react-i18next';
import SyntaxHighlighter from 'react-syntax-highlighter';

import EntityContext from '../EntityContext.js';
import { getLanguage } from '../utils/editorUtils.js';

function Test() {
  const { language, lessonVersion } = useContext(EntityContext);
  const { t } = useTranslation();

  return (
    <div>
      <p className="text-center lead">{t('testInstructions')}</p>
      <div className="hexlet-basics-content">
        <SyntaxHighlighter showLineNumbers useInlineStyles={false} language={[getLanguage(language)]}>
          {lessonVersion.test_code}
        </SyntaxHighlighter>
      </div>
    </div>
  );
}

export default Test;
