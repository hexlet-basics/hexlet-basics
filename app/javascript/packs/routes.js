const host = '';
const prefix = 'api';

export default {
  checkLessonPath: (lessonId) => [host, prefix, 'lessons', lessonId, 'check'].join('/'),
};
