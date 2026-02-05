export default {
  attributes: {
    remind_password_form: {
      email: 'Email',
    },
    user_password_form: {
      password: 'Новый пароль',
    },
    user_sign_in_form: {
      email: 'Email',
      password: 'Пароль',
    },
    user_sign_up_form: {
      email: 'Email',
      first_name: 'Имя',
      password: 'Пароль',
    },
  },
  errors: {
    models: {
      remind_password_form: {
        attributes: {
          email: {
            blank: 'Email не может быть пустым',
            user_does_not_exist:
              'Пользователь не найден. Попробуйте зарегистрироваться',
          },
        },
      },
      sign_in_form: {
        attributes: {
          email: {
            blank: 'Email не может быть пустым',
            user_does_not_exist_html:
              'Пользователь не найден. Попробуйте зарегистрироваться',
          },
          password: {
            blank: 'Пароль не может быть пустым',
            cannot_sign_in: 'Пароль не подходит. Проверьте раскладку',
          },
        },
      },
    },
  },
} as const;
