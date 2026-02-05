export default {
  account: {
    profiles: {
      destroy: {
        success: 'Account successfully deleted',
      },
      edit: {
        confirm_delete: 'Are you sure you want to delete your account?',
        delete: 'Delete account',
        meta_description: 'Profile',
        title: 'Profile editing',
      },
    },
  },
  admin: {
    blog_posts: {
      edit: {
        header: 'Edit Blog Post',
      },
      index: {
        add_new_blog_post: 'Add New Blog Post',
        header: 'Blog Posts',
      },
      'new': {
        header: 'New Blog Post',
      },
    },
    home: {
      index: {
        dashboard: 'Dashboard',
      },
    },
    languages: {
      edit: {
        header: 'Edit Language',
      },
      index: {
        add_new_language: 'Add new language',
        header: 'Languages',
        languages: 'Languages',
      },
      'new': {
        header: 'New Language',
      },
    },
    management: {
      users: {
        edit: {
          header: 'Edit User',
        },
        filter: {
          from: 'From',
        },
        index: {
          finished_lessons: 'Finished Lessons',
          header: 'Users',
          search_by_email: 'Search by email',
        },
      },
    },
    reviews: {
      edit: {
        header: 'Edit Review',
      },
      index: {
        add_new_review: 'Add New Review',
        header: 'Reviews',
      },
      'new': {
        header: 'New Review',
      },
    },
  },
  blog_posts: {
    index: {
      empty: 'Looks like there is noting yet',
      header: 'Blog',
      meta: {
        description:
          'Explore insightful articles on learning programming languages at Code-basics.com.\nDiscover tips, tutorials, and guides to boost your coding skills,\nfrom beginners to advanced topics in various programming languages.\n',
      },
    },
    show: {
      blog_posts: 'Recommended Posts',
      breadcrumbs: 'Breadcrumbs',
      courses: 'Recommended Courses',
      to_home_title: 'To home',
    },
  },
  cases: {
    for_teachers: 'For school teachers',
    link: 'Go to',
    meta: {
      description: 'Examples of using Code Basics in different areas',
    },
    title: 'Cases',
  },
  errors: {
    show: {
      codes: {
        '403': {
          description:
            'Most likely, this is administrative access, so no entry for you :)',
          header: 'Access Denied',
        },
        '404': {
          description: 'The page may have been deleted',
          header: 'Page Not Found',
        },
        '500': {
          description:
            'Something went wrong. The developers have already received an SMS and are rushing to their computers to deploy a fix as soon as possible',
          header: 'Server Error',
        },
        other: {
          description:
            'Something went wrong. We have received a notification and are investigating the issue carefully',
          header: 'Error',
        },
      },
    },
  },
  for_teachers: {
    description:
      'Free courses with a simulator right in your browser. Practice during each lesson',
    early_career_guidance: 'Stimulate early career guidance:',
    early_career_guidance_list: [
      'We show programming from the inside out to make an informed decision about studying your future profession',
      'We give you the opportunity to quickly immerse yourself in a programming language and see if the student likes it',
      'We help you to see basic examples of real-life tasks to take a closer look at your professional routine',
    ],
    header: 'The fastest way to teach students programming',
    how_to_learn_programming:
      'How do students learn programming in CodeBasics?',
    how_to_learn_programming_cards: [
      {
        title:
          'Integration without integration: to implement CodeBasics in the educational process, you will need a PC with Internet access',
        subtitle:
          'You need to open your browser, register and start taking the course you are interested in',
        img: 'integration-icon',
      },
      {
        title:
          'Theory backed by practice: PHP, Java, JS, Python, Ruby, HTML, CSS, Racket, Elixir and Go simulators are available on the platform now',
        subtitle: 'It is possible to take several courses at the same time',
        img: 'practice-icon',
      },
      {
        title:
          'Interactive format: the system automatically checks the completed tasks, shows errors and correct solution, as well as a detailed output',
        subtitle:
          "If something doesn't work out, you can use the Teacher's Solution",
        img: 'interactive-format-icon',
      },
      {
        title:
          'Discussing complex issues: pupils can clarify unclear details in the “Discussion” section, where they will be answered by tech support or other pupils',
        subtitle:
          'An Artificial Intelligence based Virtual Mentor is also available to them',
        img: 'discussion-icon',
      },
    ],
    integrate_into_education:
      'Why integrate CodeBasics into the education of school children?',
    integrate_now: 'Integrate CodeBasics into your school curriculum now!',
    interactive_approach: 'We focus on practice and an interactive approach:',
    interactive_approach_list: [
      'Accompany theoretical course materials with practical assignments directly in the browser',
      'Give instant feedback: show the correct solution, give a conclusion for each test',
      'Connect an AI tutor who explains theory or the essence of a task to students',
    ],
    lay_programming_foundations:
      'CodeBasics lays the foundations of programming in an interactive way, for which it won 3rd place in the Best Off-Platform Online Course category in the 2020 EdCrunch Award OOC International Open Online Courses Competition',
    meta: {
      description:
        'Programming courses for children and teenagers from Hexlet online school. All classes are designed by experienced teachers and programmers, taking into account the age characteristics of children. Teaching programming to schoolchildren online',
    },
    open_browser_and_sign_up: 'Just open your browser and sign up',
    programming_basic_list: [
      'Explain how systems function in general, rather than the specifics of technology at the startup',
      'We teach novice programmers the best coding standards, such as proper function naming',
      'We cover the major programming languages and continue to add courses on modern technology',
    ],
    programming_competently:
      'Laying the foundations of programming intelligently:',
    quotes_icon: 'Quotes',
    select_course: 'Select a course',
    sign_up: 'Sign up',
    sign_up_and_start_learning:
      'Sign up and start teaching programming to your students now',
    title: 'Online programming courses for high school students - Hexlet',
    try: 'try',
  },
  home: {
    hero: null,
    index: {
      all_blog_posts: 'All Blog Posts',
      all_reviews: 'All Reviews',
      blog_posts: 'Posts',
      categories: 'Categories',
      header: 'Learn to code. Online. For free.',
      hero: {
        ai_count: 'AI Assistant',
        ai_count_description: 'Hints and code analysis',
        community_count: 'Community',
        community_count_description: 'Over 8,000 people',
        courses_count: 'Text Courses',
        courses_count_description: "theory, exercises, teacher's solutions",
        fastest_way_to_start_coding:
          'With practice in the trainer and an assistant powered by ChatGPT',
        free_programming_courses: 'Free Programming Courses from Scratch',
        how_coding_works: 'How Learning Works on Code Basics',
        learn: 'Learn',
        learn_xs: 'Learn to Code',
        source_code: 'Source Code →',
        try: 'Courses →',
      },
      hexlet: 'Hexlet',
      in_development: 'Coming soon',
      in_development_description:
        'Courses below are not in production at the moment. However, you can contribute by sending a pull request to Github repository.',
      join: 'Sign up and start learning. For free. Forever',
      meta: {
        description:
          'Learn the Code Basics of HTML, CSS and JavaScript with Interactive Coding Environment right in your Browser! Perfect for Beginners & 100% for Free | Join our Online Programming Courses\n',
      },
      reviews: 'Reviews',
      sign_up: 'Sign Up',
      start: 'Start',
      subheader: 'For those who start from scratch. From the creators of ',
      title: 'Free Online Programming Courses: HTML, CSS, JavaScript\n',
    },
    languages: {
      courses: 'Courses',
    },
    reviews: {
      course_html: 'Course %{link}',
      reviews: 'Reviews',
    },
    sitemap: {
      home: 'Home',
      title: 'Sitemap',
    },
  },
  language_categories: {
    index: {
      header: 'Course categories',
      link: 'Go to',
      meta: {
        description:
          "Browse free programming courses on Code-basics.com.\nFind categories tailored to your learning needs, whether you're a beginner or looking to deepen your expertise.\nStart coding today with interactive and easy-to-follow lessons.\n",
      },
    },
    show: {
      header: 'Courses in the %{name} category',
      meta: {
        description:
          'Choose a suitable free course in the %{name} category and learn for free right in your browser with hands-on practice and help from an AI assistant',
      },
    },
  },
  languages: {
    complete_language: {
      apps: 'Write the first complete applications that will be part of your developer portfolio',
      certificate:
        'Get the Hexlet Certificate and employment assistance at the end of your study',
      continue: 'Continue Learning',
      enter_profession: 'Join the profession on Hexlet. You will:',
      first_courses_free:
        'First courses in the profession are available for free. All you need to start learning is register and join the profession.',
      first_step:
        'This is the first step into the world of professional development. What’s next?',
      grats: 'Congratulations',
      hexlet_profs: 'Other professions on Hexlet',
      knowledge: 'Get the fundamental programming knowledge',
      likely: {
        share_fb_text: "I've completed this course! #hexlet",
        share_in_fb: 'Share on Facebook',
        share_in_telegram: 'Send via Telegram',
        share_in_twitter: 'Tweet',
        share_in_vk: 'Share on VK',
        share_twitter_text: "I've completed this course! #hexlet",
        share_vk_text: "I've completed this course! #hexlet",
        twitter: 'Hexlet_IO',
      },
      tasks: 'Learn how to solve practical problems',
      title: 'You have successfully completed the basic course on Code Basics.',
    },
    lessons: {
      show: {
        breadcrumb: 'breadcrumb',
        chat: {
          guest:
            "Hi! I'm Tota, and my job is to help you learn. To activate me, please sign up or log in if you already have an account",
          hi: "Hi! My name is Tota, and I'm here to help you with any questions about theory or practice. Just ask — I'm always here for you!\n\nP.S. You can send up to 7 messages per day, so try to ask only the most important questions\n",
          not_available:
            "Our AI assistant will be here soon — we're setting things up right now",
        },
        common_questions: [
          {
            question: "The exercise doesn't pass checking. What to do? 😶",
            answer:
              "If you've reached a deadlock it's time to ask your question in the «Discussions». [How ask a question correctly](https://help.hexlet.io/en/articles/111495):\n\n* Be sure to attach the test output, without it it's almost impossible to figure out what went wrong, even if you show your code. It's complicated for developers to execute code in their heads, but having a mistake before their eyes most probably will be helpful.\n",
          },
          {
            question: 'In my environment the code works, but not here 🤨',
            answer:
              "Tests are designed so that they test the solution in different ways and against different data. Often the solution works with one kind of input data but doesn't work with others. Check the «Tests» tab to figure this out, you can find hints at the error output.\n",
          },
          {
            question: "My code is different from the teacher's one 🤔",
            answer:
              "It's fine. 🙆 One task in programming can be solved in many different ways. If your code passed all tests, it complies with the task conditions.\n\nIn some rare cases, the solution may be adjusted to the tests, but this can be seen immediately.\n",
          },
          {
            question: "I've read the lessons but nothing is clear 🙄",
            answer:
              "It's hard to make educational materials that will suit everyone. We do our best but there is always something to improve. If you see a material that is not clear to you, describe the problem in “Discussions”. It will be great if you'll write unclear points in the question form. Usually, we need a few days for corrections.\n\nBy the way, you can participate in courses improvement. There is a link below to the lessons course code which you can edit right in your browser.\n",
          },
        ],
        controls: {
          body: 'Reset Progress\n',
          header: 'Help',
          run: 'Run',
        },
        definitions: 'Definitions',
        discuss: 'AI Assistent',
        editor: 'Editor',
        instructions: 'Instructions',
        issues:
          'Found a bug? Have something to add? Pull requests are welcome!',
        lesson: 'Lesson',
        navigation: 'Navigation',
        next: 'Next →',
        only_for_signed_in_users: 'Exercise available only for signed users.',
        output: 'Output',
        please_sign_in:
          'Please sign in with your GitHub account, this is necessary to track the progress of the lessons. If you do not have an account yet, now is the time to create an account on GitHub.',
        prev: '← Previous',
        separator: ' ',
        sign_in: 'Sign In',
        sign_up_for_tracking_progress_html:
          'Be sure to <a href="%{link}" class="text-decoration-none" target="_blank">register</a> to ensure you don\'t lose the results you\'ve achieved\n',
        solution: 'Solution',
        tests: 'Tests',
        tips: 'Tips',
        title: '%{lesson_name} | %{language_name}\n',
        to_home_title: 'Home',
      },
    },
    show: {
      about_developer_community:
        "We know how hard it is to start in IT, so we've created a developer community where you're always ready to help. Here you can ask questions, get support, communicate with experienced specialists and get into the profession faster",
      about_learning: 'How the training is organised',
      ai_explanation:
        'AI explains topics, tells you how to solve assignments, and helps you at any time - like a personal tutor 24/7',
      ai_learning:
        'Ready to innovate your learning? Start learning with AI now',
      ai_without_limits: 'AI assistance without limits',
      bash_description: 'Bash is command language',
      blog_posts: 'Blog Posts',
      breadcrumbs: 'Breadcrumbs',
      browser_practice: 'Practice in the browser',
      clang_description:
        'C is a universal programming language with a compact way of writing expressions, modern mechanisms for managing data structures, and a rich set of operators.',
      community_image_preview: 'Developer community',
      continue: 'Continue Learning',
      'convenient format': 'Convenient format',
      course_graduates:
        'Join 74,761 students who have successfully completed courses',
      cover_image: 'Course cover',
      cpp_description: 'C++ is a general-purpose programming language',
      crystal_description:
        'Crystal is a general-purpose object-oriented programming language designed and developed by Arie Borenzweig, Juan Weinerman, and Brian Cardiff.',
      csharp_description:
        'C# is a modern object-oriented and type-safe programming language. Application development language for the Microsoft .NET Framework. It is used in the creation of sites, applications and games',
      css_description:
        'Cascading Style Sheets (CSS) cascading style sheets that allow you to design the content of the page in accordance with the described rules. Text styles, block layout on a page, animation - everything is described using cascading style sheets.',
      dart_description:
        'Dart is a general purpose, strongly typed compiled language. Used to develop mobile and web applications.',
      demo_description:
        'Try a demo lesson without signing up. Practice included',
      demo_start: 'Start',
      dlang_description:
        'D is a multi-paradigm, statically typed, compiled language. Has built-in error prevention - contracts and unittests',
      elixir_description:
        'Elixir is a programming language that runs on top of Erlang. How this happens is a functional language with rigorous computations, unambiguous assignment and dynamic typing, creation to support distributed, fault-tolerant, non-stop, hot-swappable applications.',
      fortran_description:
        'Fortran is a general-purpose, compiled imperative programming language that is especially suited to numeric computation and scientific computing',
      free_course: 'Free course for beginners',
      go_description:
        'Go is a general-purpose language with rich features and clear syntax. Thanks to its multi-platform, reliable, well-documented standard library and focus on convenient approaches to the development itself, Go is an ideal language for the first steps in programming.',
      haskell_description:
        'Haskell is a standardized pure functional general-purpose programming language. The main control structure is function. A distinctive feature is a serious attitude to typing',
      html_description:
        'To standardize the output of text inside the browser, an HTML standard was created that describes the rules for formatting text data for correct output. This course is devoted to the basics of HTML markup, working with typography, outputting media documents and working with forms.',
      java_description:
        'The study of programming is a difficult and lengthy process. Learning the syntax of the language it is impossible to start without it. This The course is devoted to the basics of writing programs for a Java. He is preparing a springboard for writing meaningful programs.',
      javascript_description:
        'Programming learning is an interesting and exciting process. The syntax of the language is the simplest and shortest path, without which it is impossible to start. This free online JavaScript course will introduce you to the basic concepts of the language. Learn the basics and write your first JS programs.',
      join: 'Join',
      kotlin_description:
        'Kotlin is a cross-platform, statically typed, general-purpose programming language with type inference',
      learning_conveniently:
        'Everything you need to master new topics is on one screen. Theory, practice and live examples go hand in hand. The clear structure helps you to learn in a consistent manner and not to miss important details',
      learning_preview: 'Learning preview',
      learning_program: 'Learning program',
      lessons: '%{lessons_count} with practice in the browser',
      more_than_support: 'More than Support',
      no_registration: 'Registration is not required',
      ocaml_description:
        'OCaml is a general-purpose, industrial-strength programming language with an emphasis on expressiveness and safety',
      perl_description:
        'Perl is a highly capable, feature-rich programming language with over 30 years of development',
      php_description:
        'The study of programming is a difficult and lengthy process. Learning the syntax of the language it is impossible to start without it. This The course is devoted to the basics of writing programs for a PHP. He is preparing a springboard for writing meaningful programs.',
      powershell_description:
        'PowerShell is a cross-platform task automation solution that includes a command line shell, a scripting language, and a configuration management platform.',
      prolog_description:
        'Prolog is a logic programming language associated with artificial intelligence and computational linguistics',
      python_description:
        'The study of programming is a difficult and lengthy process. Learning the syntax of the language it is impossible to start without it. This The course is devoted to the basics of writing programs for a racket. He is preparing a springboard for writing meaningful programs.',
      racket_description:
        'The racket belongs to the Lisp language family. These programmers should learn them from the very beginning. This module is about introducing the syntax and concepts underlying any lisp.',
      ready: 'Are you ready?',
      real_life_challenges:
        "You don't need to install anything - all tasks are performed right in your browser. Built-in code editor, console and automatic tests make learning comfortable. And if something fails, you can always see the teacher's solution",
      registration: 'Sign Up',
      registration_description: "Let's get started",
      reviews: 'Reviews',
      ruby_description:
        'Ruby is a language designed to make programmers happy. He took from himself the best of Lisp, Smalltalk and Perl he is object-oriented, but at the same time his declarative nature allows writing in a procedural and functional paradigm.',
      rust_description:
        'Rust is a multi-paradigm, general-purpose programming language. Rust emphasizes performance, type safety, and concurrency',
      see_all_courses_in_category: 'See all courses in %{name}',
      sign_up: 'Sign Up',
      similar_courses: 'Similar courses',
      sort_questions: 'Sorting out the questions',
      start: 'Start Learning',
      start_demo_lesson: 'Demo lesson',
      to_home_title: 'Home',
      try: 'Try It',
      try_without_registration: 'Try it without registering',
      typescript_description:
        'TypeScript is a JavaScript-based programming language. Includes a typing system',
      updated_at: 'updated %{date}',
      without_registration:
        'Start learning right away - first lessons are available without an account. If you like the format, you can register to save your progress',
    },
    success: {
      choose_profession: 'Choose course',
      description:
        "Now you have new knowledge and skills, but that's just the beginning - there are more opportunities ahead. It's time to choose your IT career and move on with your life\n",
      header: 'Congratulations, you completed the course!',
      home: 'Home',
    },
  },
  my: {
    show: {
      finished: 'Finished Courses',
      started: 'Started Courses',
    },
  },
  pages: {
    parts: {
      about: {
        meta: {
          description:
            'Code Basics quality interactive lessons on programming basics for free, for everyone.',
        },
        title: 'About',
      },
      authors: {
        meta: {
          description:
            'Code Basics is open source. Everyone can be an author. Here is some simple steps to start.\n',
        },
        title: 'How to become an author?',
      },
      cookie_policy: {
        meta: {
          description: 'Cookie policy',
        },
        title: 'Cookie Policy',
      },
      privacy: {
        meta: {
          description: 'Code Basics Privacy Policy',
        },
        title: 'Privacy Policy',
      },
      tos: {
        meta: {
          description: 'Terms of Service of the Code Basics Platform',
        },
        title: 'Terms of Service',
      },
    },
  },
  passwords: {
    edit: {
      meta_description: 'Set a new password on Code Basics',
      new_password: 'New password',
      title: 'Change password',
    },
  },
  remind_passwords: {
    'new': {
      dont_have_account: 'New to Code-Basics?',
      forgot_password: 'Forgot password?',
      login: 'Sign in',
      meta_description: 'Restore Code Basics password',
      register: 'Sign up',
      title: 'Remind password',
      trying_to_login: 'Trying to login?',
    },
  },
  reviews: {
    index: {
      course: 'Course %{language}',
      empty: 'Looks like there is noting yet',
      header: 'Reviews',
      read_more: 'Read more',
    },
  },
  sessions: {
    'new': {
      dont_have_account: 'New to Code-Basics?',
      email: 'Email',
      forgot_password: 'Forgot password?',
      meta_description: 'Sign in to Code Basics',
      password: 'Password',
      register: 'Sign up',
      reset_password: 'Reset Password',
      sign_in: 'Sign in',
      sign_in_with_github: 'Log in with GitHub',
      title: 'Sign In',
    },
  },
  shared: {
    languages: {
      course_finished: 'Course finished!',
    },
    via_social_networks: {
      via_social_networks: 'With GitHub',
    },
  },
  users: {
    'new': {
      confirmation_html:
        'By clicking Sign up, you agree to our <a href="%{url}" class="text-decoration-none" target="_blank">service conditions</a>',
      have_account: 'Already have an account?',
      meta_description: 'Sign up on Code Basics',
      sign_in: 'Sign in',
      sign_up: 'Sign up',
      title: 'Registration',
    },
  },
} as const;
