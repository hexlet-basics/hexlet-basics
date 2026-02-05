export default {
  actions: {
    create: {
      notice: '%{resource_name} was successfully created.',
    },
    destroy: {
      alert: '%{resource_name} could not be destroyed.',
      notice: '%{resource_name} was successfully destroyed.',
    },
    update: {
      notice: '%{resource_name} was successfully updated.',
    },
  },
  web: {
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
    google_auth: {
      one_tap: {
        failure:
          "We couldn't get data from Google! Requires setting up your Google account!",
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
    },
    passwords: {
      update: {
        success:
          'Password has been changed. Please, login with your new password.',
      },
    },
    remind_passwords: {
      create: {
        success: 'A password recovery instruction was sent to your email.',
      },
    },
    sessions: {
      create: {
        success: 'You have signed in successfully. You can start learning now.',
      },
    },
    users: {
      create: {
        success: 'You have signed in successfully. You can start learning now.',
      },
    },
  },
} as const;
