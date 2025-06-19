export default {
  en: {
    consentModal: {
      title: 'We use cookies',
      description: 'Cookie modal description',
      acceptAllBtn: 'Accept all',
      acceptNecessaryBtn: 'Reject all',
      showPreferencesBtn: 'Manage Individual preferences'
    },
    preferencesModal: {
      title: 'Manage cookie preferences',
      acceptAllBtn: 'Accept all',
      acceptNecessaryBtn: 'Reject all',
      savePreferencesBtn: 'Accept current selection',
      closeIconLabel: 'Close modal',
      sections: [
        {
          title: 'Somebody said ... cookies?',
          description: 'I want one!'
        },
        {
          title: 'Strictly Necessary cookies',
          description: 'These cookies are essential for the proper functioning of the website and cannot be disabled.',
          linkedCategory: 'necessary'
        },
        {
          title: 'Performance and Analytics',
          description: 'These cookies collect information about how you use our website. All of the data is anonymized and cannot be used to identify you.',
          linkedCategory: 'analytics'
        },
        {
          title: 'More information',
          description: 'For any queries in relation to my policy on cookies and your choices, please <a href="#contact-page">contact us</a>'
        }
      ]
    }
  },
  ru: {
    consentModal: {
      title: 'Мы используем файлы cookie',
      description: 'Пользуясь нашим сайтом, вы соглашаетесь с тем, что мы <a href="/ru/pages/cookie_policy">используем cookies</a>',
      acceptAllBtn: 'Принять все',
      acceptNecessaryBtn: 'Отклонить все',
      showPreferencesBtn: 'Управлять предпочтениями'
    },
    preferencesModal: {
      title: 'Управление настройками cookie',
      acceptAllBtn: 'Принять все',
      acceptNecessaryBtn: 'Отклонить все',
      savePreferencesBtn: 'Принять текущий выбор',
      closeIconLabel: 'Закрыть окно',
      sections: [
        {
          title: 'Кто-то сказал ... печенье?',
          description: 'Я хочу одно!'
        },
        {
          title: 'Строго необходимые файлы cookie',
          description: 'Эти файлы cookie необходимы для правильной работы сайта и не могут быть отключены.',
          linkedCategory: 'necessary'
        },
        {
          title: 'Производительность и аналитика',
          description: 'Эти файлы cookie собирают информацию о том, как вы используете наш сайт. Все данные анонимны и не могут быть использованы для вашей идентификации.',
          linkedCategory: 'analytics'
        },
        {
          title: 'Дополнительная информация',
          description: 'Если у вас есть вопросы относительно нашей политики в отношении cookie и ваших выборов, пожалуйста, <a href="mailto:support@hexlet.io">свяжитесь с нами</a>'
        }
      ]
    }
  }
}
