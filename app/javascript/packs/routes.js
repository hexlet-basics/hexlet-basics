const host = '';
const prefix = 'api';

export default {
  checkLessonPath: (lessonId) => [host, prefix, 'lessons', lessonId, 'check'].join('/'),
  languageLessonPath: (languageSlug, lessonSlug) => [host, 'languages', languageSlug, 'lessons', lessonSlug]
    .join('/'),
  languagePath: (languageSlug) => [host, 'languages', languageSlug].join('/'),
};
