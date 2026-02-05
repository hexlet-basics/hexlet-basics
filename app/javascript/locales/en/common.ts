export default {
  check: {
    error: {
      headline: 'Oops!',
      message: 'Something went wrong, try one more time, please',
    },
    failed: {
      headline: 'Tests Failed',
      message:
        'Your code contains errors. Read carefully tests output and try to find errors yourself',
    },
    'failed-infinity': {
      headline: 'Tests Failed (Infinity Loop!)',
      message:
        'Maybe the test ran too long or you have an infinite loop in your solution. Make sure your loop has the correct breaking condition and you consider edge cases. If your loop works correctly, send it for checking again in 5 minutes. ',
    },
    passed: {
      headline: 'Tests passed',
      message:
        'Yippee! Well done! Compare your answer with the teacher’s one and move to the next lesson',
    },
  },
  community_url: 'https://t.me/HexletLearningBot\n',
  confirm:
    'You want to reset the exercise progress. The current code version will not be saved. We hope you’ve already copied it. Continue resetting?',
  current_state: 'Current state',
  discuss: 'AI Assistent',
  editor: 'Editor',
  empty: 'Empty',
  errors: {
    network:
      'There was a network problem. Please try again. If it doesn’t work, make sure you have good internet and no blockers.',
    server:
      'Error on server. Maybe it’ll let go soon, but maybe not. Try to find out what happened in https://slack.hexlet.io/',
  },
  export: 'Экспорт',
  hello: 'Привет мир',
  hours: {
    few: '%{count} часа',
    many: '%{count} часов',
    one: '%{count} hour',
  },
  instructions: 'Instructions',
  language_icon: '%{language} icon',
  languages: {
    en: 'English',
    ru: 'Russian',
  },
  lesson: 'Lesson',
  lessons: {
    few: '%{count} урока',
    many: '%{count} уроков',
    one: '%{count} lesson',
  },
  loading: 'Loading...',
  nextLesson: 'Next',
  organization: {
    address: 'The Republic of Kazakhstan, Almaty, Auezova St., 14A',
    description:
      'Автоматизированная система для обучения программированию в браузере. Текстовые курсы + удобный тренажер.',
    email: 'support@hexlet.io',
    legal_name: 'TOO "Hexlet"',
    name: 'Code Basics',
    phone: '+7 (495) 085 21 62',
    phones: ['+7 (495) 085 21 62', '8 800 100 22 47'],
    site: 'https://ru.hexlet.io',
  },
  output: 'Output',
  pages: {
    about: 'О проекте',
  },
  prevLesson: 'Previous',
  reset: 'Reset',
  resetCode: 'Reset code',
  run: 'Run',
  sentryFeedbackWidget: {
    addScreenshotButtonLabel: 'Добавить скриншот',
    cancelButtonLabel: 'Отмена',
    formTitle: 'Сообщить об ошибке',
    isRequiredLabel: '(обязательно)',
    messageLabel: 'Описание',
    messagePlaceholder: 'Опишите ошибку и ожидаемое поведение',
    nameLabel: 'Имя',
    namePlaceholder: 'Ваше имя',
    submitButtonLabel: 'Отправить сообщение',
  },
  showSolution: 'Show solution',
  signInSuggestion:
    '<a href="/users/new">Create an account</a> to save your progress',
  solution: 'Solution',
  solutionInstructions: "Teacher's solution will be available in:",
  solutionNotice:
    "It's best to solve the problem yourself, but if you're stuck for a long time, feel free to check out the solution. But make sure to study it thoroughly to truly understand it.",
  state_events: 'State Events',
  students: {
    few: '%{count} студента',
    many: '%{count} студентов',
    one: '%{count} student',
  },
  teacherSolution: "Teacher's solution:",
  testForExercise: 'Tests',
  testInstructions: 'Your exercise will be checked with these tests:',
  time: {
    minutes_few: '%{count} минуты',
    minutes_many: '%{count} минут',
    minutes_one: '%{count} минута',
    minutes_other: '%{count} минут',
    minutes_two: '%{count} минуты',
    minutes_zero: '%{count} минут',
  },
  tos: 'Table Of Content',
  userCode: 'Your solution:',
  userCodeInstructions:
    "(start writing in Editor, your code will appear here and you'll be able to compare it to the teacher's solution)",
  views: {
    pagination: {
      first: '&laquo;',
      last: '&raquo;',
      next: '&rsaquo;',
      next_text: 'Дальше &rsaquo;',
      previous: '&lsaquo;',
      previous_text: '&lsaquo; Назад',
      truncate: '&hellip;',
    },
  },
} as const;
