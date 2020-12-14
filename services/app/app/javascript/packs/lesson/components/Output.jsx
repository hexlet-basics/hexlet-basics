import React from 'react';
import { useSelector } from 'react-redux';
import { useTranslation } from 'react-i18next';
import cn from 'classnames';
import AnsiUp from 'ansi_up';

import { checkInfoStates } from '../utils/stateMachines.js';

const ansi = new AnsiUp();

const Output = () => {
  const checkInfo = useSelector((state) => state.checkInfoSlice);
  const { t } = useTranslation();

  if (checkInfoStates.checked !== checkInfo.processState) {
    return null;
  }

  const message = t(`check.${checkInfo.result}.message`);
  const alertClassName = cn('mt-auto text-center alert', {
    'alert-success': checkInfo.passed,
    'alert-warning': !checkInfo.passed,
  });
  return (
    <div className="d-flex flex-column h-100">
      <pre>
        <code className="nohighlight" dangerouslySetInnerHTML={{ __html: ansi.ansi_to_html(checkInfo.output) }} />
      </pre>
      <div className={alertClassName}>{message}</div>
    </div>
  );
};

export default Output;
