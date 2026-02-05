export default {
  web: {
    account: {
      profiles: {
        destroy: {
          error: '',
          success: '',
        },
        update: {
          error: '',
          success: '',
        },
      },
    },
    admin: {
      blog_posts: {
        create: {
          error: '',
          success: '',
        },
        update: {
          error: '',
          success: '',
        },
      },
      languages: {
        update: {
          error: '',
          success: '',
        },
      },
      reviews: {
        create: {
          error: '',
          success: '',
        },
        update: {
          error: '',
          success: '',
        },
      },
    },
    application: {
      base: {
        error: '',
        success: '',
      },
    },
    auth: {
      callback: {
        success: '',
      },
    },
    blog_posts: {
      likes: {
        create: {
          notice: '',
          success: '',
        },
      },
    },
    books: {
      create_request: {
        success: '',
      },
    },
    google_auth: {
      one_tap: {
        error: '',
        success: '',
      },
    },
    languages: {
      lessons: {
        show: {
          lesson_not_found: '',
        },
      },
      show: {
        empty_language_current_version: '',
        language_in_development_html: '',
        warning: '',
      },
      success: {
        error: '',
      },
    },
    leads: {
      create: {
        error: '',
        success: '',
      },
    },
    passwords: {
      update: {
        success: '',
      },
    },
    remind_passwords: {
      create: {
        error: '',
        success: '',
      },
    },
    sessions: {
      create: {
        error: '',
        success: '',
      },
    },
    surveys: {
      show: {
        success: '',
      },
    },
    users: {
      create: {
        error: '',
        success: '',
      },
    },
  },
} as const;
