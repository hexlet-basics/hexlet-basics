export default {
  web: {
    account: {
      profiles: {
        destroy: {
          error: 'Возникла ошибка при удалении аккаунта',
          success: 'Аккаунт успешно удален',
        },
        update: {
          error: 'Ошибка обновления',
          success: 'Данные успешно обновлены',
        },
      },
    },
    admin: {
      blog_posts: {
        create: {
          error: 'Something went wrong',
          success: 'Post successfully created',
        },
        update: {
          error: 'Something went wrong',
          success: 'Post successfully updated!',
        },
      },
      languages: {
        update: {
          error: 'Something went wrong',
          success: 'Language updated successfully',
        },
      },
      reviews: {
        create: {
          error: 'Something went wrong',
          success: 'Testimonial successfully published',
        },
        update: {
          error: 'Something went wrong',
          success: 'Review successfully updated!',
        },
      },
    },
    application: {
      base: {
        error: 'Please check errors',
        success: 'Action completed successfully',
      },
    },
    auth: {
      callback: {
        success: 'You have signed in successfully. You can start learning now.',
      },
    },
    blog_posts: {
      likes: {
        create: {
          notice: 'Ваш лайк уже засчитан :)',
          success: 'Спасибо за лайк!',
        },
      },
    },
    books: {
      create_request: {
        success:
          'Ура, теперь книга доступна для скачивания! Нажмите на кнопку и книга скачается',
      },
    },
    google_auth: {
      one_tap: {
        error:
          'Мы не смогли получить данные из Google! Проверьте настройки вашего аккаунта Google!',
        success: 'You have signed in successfully. You can start learning now.',
      },
    },
    languages: {
      lessons: {
        show: {
          lesson_not_found: 'Lesson not found. Please, try another lesson.',
        },
      },
      show: {
        empty_language_current_version:
          'Language is under development, please check back later or try any other available language',
        language_in_development_html:
          'Language %{language} is under development.\nYou can help send a pull request with new lessons or\nsupplement those that already exist <a href="%{link_to_repo}" target="_blank" rel="noopener">%{link_to_repo}</a>.\nWe have prepared guidelines for writing <a href="%{link_to_recommendations}" target="_blank" rel="noopener">at the link</a>\n',
        warning:
          'The Language has lessons only in russian. Switch locale if you can read it :)',
      },
      success: {
        error: 'В этом курсе есть уроки, которые вы не завершили',
      },
    },
    leads: {
      create: {
        error: 'Проверьте ошибки в форме',
        success:
          'Заявка отправлена! Свяжемся с вами в течение одного-двух рабочих дней. Или напишите нам в <a target="_blank" href="https://t.me/WelcomeCodebasicsBot">телеграм</a> чтобы получить помощь быстрее',
      },
    },
    passwords: {
      update: {
        success:
          'Password has been changed. Please, login with your new password.',
      },
    },
    remind_passwords: {
      create: {
        error: 'В форме есть ошибки',
        success: 'A password recovery instruction was sent to your email.',
      },
    },
    sessions: {
      create: {
        error: 'В форме есть ошибки',
        success: 'You have signed in successfully. You can start learning now.',
      },
    },
    surveys: {
      show: {
        success: 'Вы уже отвечали на этот опрос. Возвращаемся обратно',
      },
    },
    users: {
      create: {
        error: 'Упс, кажется в форме есть ошибки',
        success: 'You have signed in successfully. You can start learning now.',
      },
    },
  },
} as const;
