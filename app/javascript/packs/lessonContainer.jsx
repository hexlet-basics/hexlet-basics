import React from 'react';
import ReactDOM from 'react-dom';
import LessonPage from './components/LessonPage';

document.addEventListener('DOMContentLoaded', () => {
  ReactDOM.render(
    <LessonPage />,
    document.querySelector('#basics-lesson-container'),
  );
});
