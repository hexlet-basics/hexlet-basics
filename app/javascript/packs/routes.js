const host = '';
const prefix = 'api';

export default {
  checkLessonVersionPath: (languageId, lessonId, versionId) => [host, prefix, 'languages', languageId, 'lessons', lessonId, 'versions', versionId, 'check'].join('/'),
};
