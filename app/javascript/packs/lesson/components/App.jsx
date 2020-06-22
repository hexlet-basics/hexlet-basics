import React from 'react';
import TabsBox from './TabsBox.jsx';
import ControlBox from './ControlBox.jsx';

const App = () => {
  const vdom = (
    <>
      <TabsBox />
      <ControlBox />
    </>
  );
  return vdom;
};

export default App;
