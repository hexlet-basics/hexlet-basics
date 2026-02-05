export default {
  attributes: {
    remind_password_form: {
      email: '',
    },
    user_password_form: {
      password: '',
    },
    user_sign_in_form: {
      email: '',
      password: '',
    },
    user_sign_up_form: {
      email: '',
      first_name: '',
      password: '',
    },
  },
  errors: {
    models: {
      remind_password_form: {
        attributes: {
          email: {
            blank: '',
            user_does_not_exist: '',
          },
        },
      },
      sign_in_form: {
        attributes: {
          email: {
            blank: '',
            user_does_not_exist_html: '',
          },
          password: {
            blank: '',
            cannot_sign_in: '',
          },
        },
      },
    },
  },
} as const;
