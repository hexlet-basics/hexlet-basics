// @ts-check

import React, { useContext } from 'react';
import { useTranslation } from 'react-i18next';
import hljs from '../../lib/hljs.js';

import EntityContext from '../EntityContext.js';
import { getLanguage } from '../utils/editorUtils.js';

function Test() {
  const { language, lessonVersion } = useContext(EntityContext);
  const { t } = useTranslation();

  const code = hljs.highlight(lessonVersion.test_code, { language: getLanguage(language) }).value;

  return (
    <div>
      <p className="text-center lead">{t('testInstructions')}</p>
      <pre>
        <code>
          <div className="hexlet-basics-content" dangerouslySetInnerHTML={{ __html: code }} />
        </code>
      </pre>
    </div>
  );
}

export default Test;
