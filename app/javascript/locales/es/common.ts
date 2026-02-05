export default {
  check: {
    error: {
      headline: 'Упс!',
      message: 'Что-то пошло не так. Попробуйте ещё раз',
    },
    failed: {
      headline: 'Тесты не прошли',
      message:
        'В вашем коде есть ошибки. Прочитайте внимательно вывод тестов, найдите и попробуйте исправить их. <a href="https://help.hexlet.io/article/20623" target="_blank">Как искать ошибки в коде</a>. Если не получается, поговорите с нашим ассистентом Тота во вкладке "ИИ-помощник".\n',
    },
    'failed-infinity': {
      headline: 'Долгое выполнение',
      message:
        'Возможно, тесты работали слишком долго или в вашем решении есть бесконечный цикл. Убедитесь, что в цикле верное условие остановки, и вы учитываете пограничные случаи. Если с циклом всё в порядке, попробуйте отправить код на проверку снова через 5 минут',
    },
    passed: {
      headline: 'Тесты пройдены',
      message:
        'Ура! Всё получилось! Сравните своё решение с решением учителя и переходите к следующему уроку',
    },
  },
  community_url: 'https://t.me/HexletLearningBot\n',
  confirm: 'Убеждаемся, что мы делаем то что действительно хотим. Хотим?',
  current_state: 'Текущее состояние',
  discuss: 'ИИ-помощник',
  editor: 'Редактор',
  empty: 'Пусто',
  errors: {
    network:
      'Возникла сетевая проблема. Попробуйте повторить ещё раз. Если не получилось, убедитесь, что с интернетом всё в порядке, и отключите блокировщики рекламы',
    server:
      'Ошибка на сервере. Возможно, скоро отпустит, а возможно — нет. Попробуйте узнать, что произошло в https://ttttt.me/hexletcommunity/12',
  },
  export: 'Экспорт',
  hello: 'Привет мир',
  hours: {
    few: '%{count} часа',
    many: '%{count} часов',
    one: '%{count} час',
  },
  instructions: 'Инструкции',
  language_icon: 'Иконка %{language}',
  languages: {
    en: 'Английский',
    ru: 'Русский',
  },
  lesson: 'Урок',
  lessons: {
    few: '%{count} урока',
    many: '%{count} уроков',
    one: '%{count} урок',
  },
  loading: 'Загрузка...',
  nextLesson: 'Следующий',
  organization: {
    address:
      '108813, г. Москва, вн.тер.г. поселение Московский, г. Московский, ул. Солнечная, д. 3А, стр. 1, помещ. 10/3',
    description:
      'Автоматизированная система для обучения программированию в браузере. Текстовые курсы + удобный тренажер.',
    email: 'support@hexlet.io',
    legal_name: 'ООО «Хекслет Рус»',
    name: 'Code Basics',
    phone: '+7 (495) 085 21 62',
    phones: ['+7 (495) 085 21 62', '8 800 100 22 47'],
    site: 'https://ru.hexlet.io',
  },
  output: 'Вывод',
  pages: {
    about: 'О проекте',
  },
  prevLesson: 'Предыдущий',
  reset: 'Сброс',
  resetCode: 'Сбросить код',
  run: 'Проверить',
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
  showSolution: 'Показать решение',
  signInSuggestion:
    '<a href="%{url}">Создайте аккаунт</a> для сохранения прогресса',
  solution: 'Решение',
  solutionInstructions: 'Решение учителя откроется через:',
  solutionNotice:
    'Желательно решить задачу самостоятельно, но если вы застряли и долгое время ничего не получается, посмотрите решение учителя. Но обязательно разберитесь в нём и повторите по памяти',
  state_events: 'Варианты действий',
  students: {
    few: '%{count} студента',
    many: '%{count} студентов',
    one: '%{count} студент',
  },
  teacherSolution: 'Решение учителя:',
  testForExercise: 'Тесты',
  testInstructions: 'Ваше упражнение проверяется по этим тестам',
  time: {
    minutes_few: '%{count} минуты',
    minutes_many: '%{count} минут',
    minutes_one: '%{count} минута',
    minutes_other: '%{count} минут',
    minutes_two: '%{count} минуты',
    minutes_zero: '%{count} минут',
  },
  tos: 'Содержание',
  userCode: 'Ваше решение:',
  userCodeInstructions:
    'Когда вы начнёте писать решение в Редакторе, оно появится тут для сравнения с учительским',
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
