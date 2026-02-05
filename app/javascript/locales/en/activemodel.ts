export default {
  attributes: {
    remind_password_form: {
      email: 'Email',
    },
    user_password_form: {
      password: 'Password',
    },
    user_sign_in_form: {
      email: 'Email',
      password: 'Password',
    },
    user_sign_up_form: {
      email: 'Email',
      first_name: 'First Name',
      password: 'Password',
    },
  },
  errors: {
    models: {
      remind_password_form: {
        attributes: {
          email: {
            blank: "Email can't be blank",
            user_does_not_exist: 'User not found. Try signing up.',
          },
        },
      },
      sign_in_form: {
        attributes: {
          email: {
            blank: "Email can't be blank",
            user_does_not_exist_html: 'User not found. Try signing up.',
          },
          password: {
            blank: "Password can't be blank",
            cannot_sign_in: 'Wrong password. Check your keyboard layout.',
          },
        },
      },
    },
  },
} as const;
