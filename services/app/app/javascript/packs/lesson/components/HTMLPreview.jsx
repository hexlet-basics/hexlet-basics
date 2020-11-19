import React from 'react';
import PropTypes from 'prop-types';
import Frame from 'react-frame-component';

const HTMLPreview = (props) => {
  const { html } = props;
  return (
    <div className="mt-2 p-2 h-50 bg-white">
      <Frame className="border-0 h-100 w-100">
        <div className="text-black" dangerouslySetInnerHTML={{ __html: html }} />
      </Frame>
    </div>
  );
};

HTMLPreview.propTypes = {
  html: PropTypes.string.isRequired,
};

export default HTMLPreview;
