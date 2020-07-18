const host = '';
const prefix = 'api';

export default {
  lessonCheckPath: (lessonId) => [host, prefix, 'lessons', lessonId, 'check'].join('/')
};
