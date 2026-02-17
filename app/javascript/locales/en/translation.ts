export default {
  translation: {
    account: {
      profiles: {
        edit: {
          delete: "Delete account",
          title: "Profile editing",
        },
      },
    },
    activemodel: {
      attributes: {
        remind_password_form: {
          email: "Email",
        },
        user_password_form: {
          password: "Password",
        },
        user_sign_in_form: {
          email: "Email",
          password: "Password",
        },
        user_sign_up_form: {
          email: "Email",
          first_name: "First Name",
          password: "Password",
        },
      },
      errors: {
        models: {
          remind_password_form: {
            attributes: {
              email: {
                blank: "Email can't be blank",
                user_does_not_exist: "User not found. Try signing up.",
              },
            },
          },
          sign_in_form: {
            attributes: {
              email: {
                blank: "Email can't be blank",
                user_does_not_exist_html: "User not found. Try signing up.",
              },
              password: {
                blank: "Password can't be blank",
                cannot_sign_in: "Wrong password. Check your keyboard layout.",
              },
            },
          },
        },
      },
    },
    activerecord: {
      attributes: {
        base: {
          _destroy: "Удалить",
          answer: "Ответ",
          created_at: "Создание",
          description: "Описание",
          header: "Заголовок",
          name: "Имя",
          question: "Вопрос",
          slug: "Слаг",
          updated_at: "Обновление",
        },
        blog_post: {
          body: "Content",
          cover: "Image",
          creator: "Author",
          state: "State",
          "state/archived": "Archived",
          "state/published": "Published",
        },
        language: {
          category_id: "Category",
          cover: "Cover",
          hexlet_program_landing_page: "Hexlet Program Landing Page",
          learn_as: "Learn as",
          openai_assistant_id: "Assistant ID (OpenAI)",
          progress: "Progress",
          slug: "Slug",
        },
        language_category: {
          description: "Description",
          header: "Header",
          name: "Name",
          slug: "Slug",
        },
        language_category_item: {
          language_landing_page_id: "Лендинг",
        },
        language_landing_page: {
          _destroy: "Destroy",
          answer: "Answer",
          description: "Description",
          footer: "Show in Footer",
          footer_name: "Footer Name",
          header: "Header",
          landing_page_to_redirect_id: "Redirect to Page",
          language_id: "Course",
          listed: "Show in Lists",
          main: "Main",
          meta_description: "Meta Description",
          meta_title: "Meta Title",
          name: "Name",
          order: "Order",
          outcomes_description: "Learning Outcomes (Description)",
          outcomes_header: "Learning Outcomes (Header)",
          outcomes_image: "Learning Outcomes (Image)",
          question: "Question",
          slug: "Slug",
          state: "State",
          used_in_description: "Used In (Description)",
          used_in_header: "Used In (Header)",
        },
        lead: {
          contact_method: "Способ связи",
          contact_value: "Телефон / Имя пользователя",
        },
        qna_item: {
          answer: "Ответ",
          footer: "Выводить в футере",
          footer_name: "Имя в футере",
          listed: "Отображать в списках",
          main: "Основной",
          order: "Order",
          question: "Вопрос",
          slug: "Slug",
        },
        review: {
          body: "Review",
          first_name: "First name",
          language_id: "Course",
          last_name: "Last name",
          locale: "Locale",
          pinned: "Pinned (Shown at the top of the list)",
          state: "State",
          user_id: "User",
        },
        survey: {
          run_always: "Запускать всегда",
          slug: "Слаг",
        },
        survey_scenario: {
          _destroy: "Удалить",
          name: "Название",
          survey_item_id: "Ответ триггер (Начало воронки)",
        },
        user: {
          admin: "Admin?",
          contact_method: "Способ связи",
          contact_value: "Номер для связи (или имя пользователя)",
          email: "Email",
          first_name: "First name",
          last_name: "Last name",
          nickname: "Nickname",
          password: "Password",
        },
      },
      errors: {
        messages: {
          record_invalid: "Validation failed: %{errors}",
          restrict_dependent_destroy: {
            has_one:
              "Cannot delete record because a dependent %{record} exists",
            has_many: "Cannot delete record because dependent %{record} exist",
          },
        },
        models: {
          sign_up_form: {
            attributes: {
              password: {
                too_short: "должен быть не меньше 6 символов",
              },
            },
          },
        },
      },
    },
    admin: {
      blog_posts: {
        edit: {
          header: "Edit Blog Post",
        },
        index: {
          add_new_blog_post: "Add New Blog Post",
          header: "Blog Posts",
        },
        "new": {
          header: "New Blog Post",
        },
      },
      home: {
        index: {
          dashboard: "Dashboard",
        },
      },
      language_categories: {
        edit: {
          header: "Редактирование %{id}",
        },
        form: {
          items: "Items",
          main: "Main",
          qna_items: "Q&A Items",
        },
        index: {
          add_new_language: "Добавить категорию",
          header: "Категории",
        },
        "new": {
          header: "Новая категория",
        },
      },
      language_landing_pages: {
        edit: {
          header: "Редактирование %{id}",
        },
        index: {
          add_new_language: "Добавить лендинг",
          header: "Лендинги",
        },
        "new": {
          header: "Новый лендинг",
        },
      },
      language_lesson_members: {
        index: {
          header: "Прохождение уроков",
        },
      },
      language_lesson_reviews: {
        index: {
          data: "Data",
          header: "Ревью уроков",
        },
      },
      language_lessons: {
        index: {
          header: "Уроки",
        },
      },
      languages: {
        edit: {
          header: "Edit Language",
          id: "ID",
          load_new_version: "Load new version",
          result: "Result",
          versions: "Versions",
        },
        index: {
          add_new_language: "Add new language",
          header: "Languages",
        },
        "new": {
          header: "New Language",
        },
      },
      leads: {
        index: {
          answers: "Answers",
          courses: "Courses",
          data: "Data",
          header: "Лиды",
        },
      },
      management: {
        users: {
          edit: {
            count: "Count",
            header: "Edit User",
            language: "Language",
            progress: "Progress",
          },
          filter: {
            from: "From",
          },
          index: {
            finished_lessons: "Finished Lessons",
            header: "Users",
            search_by_email: "Search by email",
          },
        },
      },
      messages: {
        index: {
          header: "Разговоры с AI",
        },
      },
      reviews: {
        edit: {
          header: "Edit Review",
        },
        index: {
          add_new_review: "Add New Review",
          header: "Reviews",
        },
        "new": {
          header: "New Review",
        },
      },
      survey_answers: {
        index: {
          header: "Ответы на опросы",
        },
      },
      survey_scenarios: {
        edit: {
          header: "Редактирование сценария",
        },
        index: {
          header: "Сценарии",
        },
        "new": {
          header: "Новый сценарий",
        },
      },
      surveys: {
        edit: {
          header: "Редактирование опроса",
        },
        index: {
          header: "Опросы",
        },
        "new": {
          header: "Новый опрос",
        },
      },
      title: "Админка Code-Basics",
    },
    api: {
      partners: {
        yandex_market: {
          languages: {
            index: {
              category: "Freemium",
              company: "Hexlet",
              title: "Education programs",
            },
          },
        },
      },
    },
    blocks: {
      lead_form_block: {
        description1: "Ответим в течение дня",
        description2: "Нам можно написать в <a>телеграм</a>",
      },
    },
    blog_posts: {
      index: {
        empty: "Looks like there is noting yet",
        header: "Blog",
        meta: {
          description:
            "Explore insightful articles on learning programming languages at Code-basics.com.\nDiscover tips, tutorials, and guides to boost your coding skills,\nfrom beginners to advanced topics in various programming languages.\n",
        },
      },
      show: {
        blog_posts: "Recommended Posts",
        breadcrumbs: "Breadcrumbs",
        courses: "Recommended Courses",
        discuss: "Обсуждайте, помогайте и делитесь опытом",
        join_community: "Присоединяйтесь к сообществу Хекслет",
        link: "Перейти в сообщество",
        to_home_title: "To home",
      },
    },
    books: {
      show: {
        chapter: "Глава %{number}",
        description:
          "С нуля до трудоустройства. Этот учебник — ваш путеводитель в профессию программиста. Он создан для тех, кто хочет начать с нуля и шаг за шагом дойти до трудоустройства в одной из самых востребованных и перспективных областей\n",
        download: "Скачать книгу",
        features: {
          direction: "Направление",
          direction_explanation:
            "Узнаете какие бывают виды разработки и языки. Куда стоит идти и почему",
          interview: "Интервью",
          interview_explanation:
            "Подготовитесь к прохождению собеседований и испытательного срока",
          plan: "План",
          plan_explanation:
            "Получите готовый план обучения, с темами, ресурсами для изучения и проектами",
          resume: "Резюме",
          resume_explanation:
            "Научитесь правильно оформлять резюме и находить подходящие вакансии",
        },
        freebook: "Бесплатная электронная книга",
        header: "Книга: Профессия программист. С нуля до трудоустройства\n",
        request: "Запросить книгу",
        toc: "Содержание",
      },
    },
    cases: {
      for_teachers: {
        description:
          "Бесплатные курсы с тренажером прямо в браузере. Много практики и ИИ-ассистент",
        early_career_guidance: "Стимулируем раннюю профориентацию:",
        early_career_guidance_list: [
          "Показываем программирование изнутри, чтобы принять взвешенное решение об изучении будущей профессии",
          "Даем возможность быстро погрузиться в язык программирования и понять, нравится ли он ученику",
          "Помогаем увидеть базовые примеры реальных задач, чтобы присмотреться к профессиональной рутине",
        ],
        header: "Эффективное обучение школьников и студентов программированию",
        how_to_learn_programming:
          "Как пользователи изучают программирование в CodeBasics?",
        how_to_learn_programming_cards: [
          {
            title:
              "Интеграция без интеграции: чтобы внедрить CodeBasics в образовательный процесс, потребуется ПК с доступом в интернет",
            subtitle:
              "Нужно открыть браузер, зарегистрироваться и начать проходить интересующий курс",
            img: "integration-icon",
          },
          {
            title:
              "Теория, подкрепленная практикой: на платформе уже сейчас доступны тренажеры по PHP, Java, JS, Python, Ruby, HTML, CSS, Racket, Elixir и Go",
            subtitle: "Можно проходить несколько курсов одновременно",
            img: "practice-icon",
          },
          {
            title:
              "Интерактивный формат: система автоматически проверяет выполненные задания, показывает ошибки и правильное решение, а также детальный вывод",
            subtitle:
              "Если что-то не получается, можно воспользоваться Решением учителя",
            img: "interactive-format-icon",
          },
          {
            title:
              "Пользователям доступен Виртуальный наставник на основе Искусственного интеллекта",
            subtitle: "",
            img: "discussion-icon",
          },
        ],
        integrate_into_education: "Зачем интегрировать CodeBasics в обучение?",
        integrate_now: "Интегрируйте CodeBasics в учебную программу сейчас!",
        interactive_approach:
          "Ориентируемся на практику и интерактивный подход:",
        interactive_approach_list: [
          "Сопровождаем теоретические материалы курса практическими заданиями прямо в браузере",
          "Даем мгновенную обратную связь: показываем правильное решение, даем вывод по каждому тесту",
          "Подключаем ИИ-наставника, который объясняет ученикам теорию или суть задания",
        ],
        lay_programming_foundations:
          "CodeBasics закладывает основы программирования в интерактивной форме, за что получил 3-е место в номинации «Лучший внеплатформенный онлайн-курс» в 2020 году на Международном конкурсе открытых онлайн-курсов EdCrunch Award OOC",
        meta: {
          description:
            "Курсы программирования для детей и подростков от онлайн-школы Хекслет. Все занятия разработаны опытными педагогами и программистами, с учетом возрастных особенностей детей. Обучение программированию онлайн",
        },
        open_browser_and_sign_up: "Просто откройте браузер и зарегистрируйтесь",
        programming_basic_list: [
          "Объясняем, как функционируют системы в общем, а не специфику технологий на старте",
          "Обучаем начинающих программистов лучшим стандартам кодирования, например правильному именованию функций",
          "Охватываем основные языки программирования и продолжаем добавлять курсы по современным технологиям",
        ],
        programming_competently:
          "Закладываем основы программирования грамотно:",
        quotes_icon: "Ковычки",
        select_course: "Выбрать курс",
        sign_up: "Зарегистрироваться",
        sign_up_and_start_learning:
          "Зарегистрируйтесь и начните обучать ваших учеников программированию прямо сейчас",
        title: "Обучение программированию школьников и студентов",
        try: "Попробовать",
      },
      index: {
        for_teachers: "Для учителей и преподавателей",
        link: "Перейти",
        meta: {
          description: "Примеры использования Code Basics в разных сферах",
        },
        title: "Кейсы",
      },
    },
    common: {
      boolean: {
        no: "No",
        yes: "Yes",
      },
      check: {
        error: {
          headline: "Oops!",
          message: "Something went wrong, try one more time, please",
        },
        failed: {
          headline: "Tests Failed",
          message:
            "Your code contains errors. Read carefully tests output and try to find errors yourself",
        },
        "failed-infinity": {
          headline: "Tests Failed (Infinity Loop!)",
          message:
            "Maybe the test ran too long or you have an infinite loop in your solution. Make sure your loop has the correct breaking condition and you consider edge cases. If your loop works correctly, send it for checking again in 5 minutes. ",
        },
        passed: {
          headline: "Tests passed",
          message:
            "Yippee! Well done! Compare your answer with the teacher’s one and move to the next lesson",
        },
      },
      community_url: "https://t.me/HexletLearningBot\n",
      confirm:
        "You want to reset the exercise progress. The current code version will not be saved. We hope you’ve already copied it. Continue resetting?",
      current_state: "Current state",
      discuss: "AI Assistent",
      editor: "Editor",
      empty: "Empty",
      errors: {
        network:
          "There was a network problem. Please try again. If it doesn’t work, make sure you have good internet and no blockers.",
        server:
          "Error on server. Maybe it’ll let go soon, but maybe not. Try to find out what happened in https://slack.hexlet.io/",
      },
      export: "Экспорт",
      hello: "Привет мир",
      hours: {
        few: "%{count} часа",
        many: "%{count} часов",
        one: "%{count} hour",
      },
      instructions: "Instructions",
      language_icon: "%{language} icon",
      languages: {
        en: "English",
        ru: "Russian",
      },
      lesson: "Lesson",
      lessons: {
        few: "%{count} урока",
        many: "%{count} уроков",
        one: "%{count} lesson",
      },
      loading: "Loading...",
      nextLesson: "Next",
      organization: {
        address: "The Republic of Kazakhstan, Almaty, Auezova St., 14A",
        description:
          "Автоматизированная система для обучения программированию в браузере. Текстовые курсы + удобный тренажер.",
        email: "support@hexlet.io",
        legal_name: 'TOO "Hexlet"',
        name: "Code Basics",
        phone: "+7 (495) 085 21 62",
        phones: ["+7 (495) 085 21 62", "8 800 100 22 47"],
        site: "https://ru.hexlet.io",
      },
      output: "Output",
      pages: {
        about: "О проекте",
      },
      prevLesson: "Previous",
      reset: "Reset",
      resetCode: "Reset code",
      run: "Run",
      sentryFeedbackWidget: {
        addScreenshotButtonLabel: "Добавить скриншот",
        cancelButtonLabel: "Отмена",
        formTitle: "Сообщить об ошибке",
        isRequiredLabel: "(обязательно)",
        messageLabel: "Описание",
        messagePlaceholder: "Опишите ошибку и ожидаемое поведение",
        nameLabel: "Имя",
        namePlaceholder: "Ваше имя",
        submitButtonLabel: "Отправить сообщение",
      },
      showSolution: "Show solution",
      signInSuggestion:
        '<a href="/users/new">Create an account</a> to save your progress',
      solution: "Solution",
      solutionInstructions: "Teacher's solution will be available in:",
      solutionNotice:
        "It's best to solve the problem yourself, but if you're stuck for a long time, feel free to check out the solution. But make sure to study it thoroughly to truly understand it.",
      state_events: "State Events",
      students: {
        few: "%{count} студента",
        many: "%{count} студентов",
        one: "%{count} student",
      },
      teacherSolution: "Teacher's solution:",
      testForExercise: "Tests",
      testInstructions: "Your exercise will be checked with these tests:",
      time: {
        minutes_zero: "%{count} минут",
        minutes_one: "%{count} минута",
        minutes_two: "%{count} минуты",
        minutes_few: "%{count} минуты",
        minutes_many: "%{count} минут",
        minutes_other: "%{count} минут",
      },
      tos: "Table Of Content",
      userCode: "Your solution:",
      userCodeInstructions:
        "(start writing in Editor, your code will appear here and you'll be able to compare it to the teacher's solution)",
      views: {
        pagination: {
          first: "&laquo;",
          last: "&raquo;",
          next: "&rsaquo;",
          next_text: "Дальше &rsaquo;",
          previous: "&lsaquo;",
          previous_text: "&lsaquo; Назад",
          truncate: "&hellip;",
        },
      },
    },
    date: {
      abbr_day_names: ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"],
      abbr_month_names: [
        "",
        "Jan",
        "Feb",
        "Mar",
        "Apr",
        "May",
        "Jun",
        "Jul",
        "Aug",
        "Sep",
        "Oct",
        "Nov",
        "Dec",
      ],
      day_names: [
        "Sunday",
        "Monday",
        "Tuesday",
        "Wednesday",
        "Thursday",
        "Friday",
        "Saturday",
      ],
      formats: {
        default: "%Y-%m-%d",
        long: "%B %d, %Y",
        short: "%b %d",
      },
      month_names: [
        "",
        "January",
        "February",
        "March",
        "April",
        "May",
        "June",
        "July",
        "August",
        "September",
        "October",
        "November",
        "December",
      ],
      order: ["year", "month", "day"],
    },
    datetime: {
      distance_in_words: {
        about_x_hours: {
          few: "около %{count} часов",
          many: "около %{count} часов",
          one: "about %{count} hour",
          other: "about %{count} hours",
        },
        about_x_months: {
          few: "около %{count} месяцев",
          many: "около %{count} месяцев",
          one: "about %{count} month",
          other: "about %{count} months",
        },
        about_x_years: {
          few: "около %{count} лет",
          many: "около %{count} лет",
          one: "about %{count} year",
          other: "about %{count} years",
        },
        almost_x_years: {
          few: "почти %{count} года",
          many: "почти %{count} лет",
          one: "almost %{count} year",
          other: "almost %{count} years",
        },
        half_a_minute: "half a minute",
        less_than_x_minutes: {
          few: "меньше %{count} минут",
          many: "меньше %{count} минут",
          one: "less than a minute",
          other: "less than %{count} minutes",
        },
        less_than_x_seconds: {
          few: "меньше %{count} секунд",
          many: "меньше %{count} секунд",
          one: "less than %{count} second",
          other: "less than %{count} seconds",
        },
        over_x_years: {
          few: "больше %{count} лет",
          many: "больше %{count} лет",
          one: "over %{count} year",
          other: "over %{count} years",
        },
        x_days: {
          few: "%{count} дня",
          many: "%{count} дней",
          one: "%{count} day",
          other: "%{count} days",
        },
        x_minutes: {
          few: "%{count} минуты",
          many: "%{count} минут",
          one: "%{count} minute",
          other: "%{count} minutes",
        },
        x_months: {
          few: "%{count} месяца",
          many: "%{count} месяцев",
          one: "%{count} month",
          other: "%{count} months",
        },
        x_seconds: {
          few: "%{count} секунды",
          many: "%{count} секунд",
          one: "%{count} second",
          other: "%{count} seconds",
        },
        x_years: {
          few: "%{count} года",
          many: "%{count} лет",
          one: "%{count} year",
          other: "%{count} years",
        },
      },
      prompts: {
        day: "Day",
        hour: "Hour",
        minute: "Minute",
        month: "Month",
        second: "Seconds",
        year: "Year",
      },
    },
    enums: {
      language: {
        learn_as: {
          first_language: "basics for beginners",
          second_language: "as a second language",
        },
        progress: {
          completed: "Completed",
          draft: "Draft",
          in_development: "In Development",
        },
      },
    },
    errors: {
      show: {
        codes: {
          "403": {
            description:
              "Most likely, this is administrative access, so no entry for you :)",
            header: "Access Denied",
          },
          "404": {
            description: "The page may have been deleted",
            header: "Page Not Found",
          },
          "500": {
            description:
              "Something went wrong. The developers have already received an SMS and are rushing to their computers to deploy a fix as soon as possible",
            header: "Server Error",
          },
          other: {
            description:
              "Something went wrong. We have received a notification and are investigating the issue carefully",
            header: "Error",
          },
        },
      },
    },
    faq: {
      header: "Frequently Asked Questions",
      main: {
        "about-project": {
          answer:
            "It’s a completely free platform for learning IT from scratch. Code Basics was created by the developers and community of the Hexlet programming school to help anyone try programming and build a solid foundation for a new profession—based not on memorization of specifics, but on understanding the system as a whole. One of the key features of the project is its programming simulator.\n",
          question: "What is Code Basics?",
        },
        "for-whom": {
          answer:
            "Code Basics courses are suitable for anyone starting to learn programming from scratch and aiming to become a professional developer. Age or profession doesn't matter. Adults and teenagers, tech-savvy learners and pure humanities majors all study on Code Basics.\n",
          question: "Who are the courses for?",
        },
        "how-it-works": {
          answer:
            "Code Basics is a fully automated online platform, so you can study whenever it’s convenient for you.\n\nProgramming courses are divided into thematic blocks, each containing several lessons. Go through them in order: first study the theory, then complete a hands-on exercise directly in the browser—your solution will be checked automatically.\n",
          question: "How does programming education work?",
        },
        "how-much": {
          answer:
            "Code Basics was created as a free project to teach programming from scratch. That’s how it was, is, and will remain. Moreover, Code Basics is an open-source project—you can find its code on GitHub and even contribute to its development.\n",
          question: "How much do the courses cost?",
        },
        "how-to-start": {
          answer:
            "Just register. Registration is free, and after signing up, you’ll get access to all the courses on Code Basics. Choose any course and start learning.\n",
          question: "How do I start using the platform?",
        },
        "if-something-goes-wrong": {
          answer:
            "Just ask! Each lesson includes an \"AI Assistant\" section where you can chat with our assistant, Tota AI. You can ask questions or clarify the assignment if something seems unclear.\n\nWhen working on practical tasks, we also recommend carefully reading the test output—it contains helpful information about why your solution didn't pass. You can also compare your solution with the correct one. However, we suggest only doing this after successfully completing the task on your own. Remember: you're here to learn, not just copy.\n",
          question: "What if I get stuck?",
        },
        possibilities: {
          answer:
            "Code Basics is an open-source project. Any practicing developer can contribute to its growth. The source code is available on [GitHub](https://github.com/hexlet-basics). You can help by creating new courses and lessons, improving existing ones, fixing typos and inaccuracies, or translating courses into English.\n",
          question: "What opportunities are there for experienced developers?",
        },
        "what-type-of-courses": {
          answer:
            "Code Basics was originally focused solely on programming, but we are considering expanding the platform to include a broader range of IT courses. Eventually, it will become a full-fledged free IT education platform for beginners.\n",
          question: "Do you only teach programming?",
        },
      },
    },
    flash: {
      account: {
        profiles: {
          destroy: {
            error: "Возникла ошибка при удалении аккаунта",
            success: "Аккаунт успешно удален",
          },
          update: {
            error: "Ошибка обновления",
            success: "Данные успешно обновлены",
          },
        },
      },
      admin: {
        blog_posts: {
          create: {
            error: "Something went wrong",
            success: "Post successfully created",
          },
          update: {
            error: "Something went wrong",
            success: "Post successfully updated!",
          },
        },
        languages: {
          update: {
            error: "Something went wrong",
            success: "Language updated successfully",
          },
        },
        reviews: {
          create: {
            error: "Something went wrong",
            success: "Testimonial successfully published",
          },
          update: {
            error: "Something went wrong",
            success: "Review successfully updated!",
          },
        },
      },
      application: {
        base: {
          error: "Please check errors",
          success: "Action completed successfully",
        },
      },
      auth: {
        callback: {
          success:
            "You have signed in successfully. You can start learning now.",
        },
      },
      blog_posts: {
        likes: {
          create: {
            notice: "Ваш лайк уже засчитан :)",
            success: "Спасибо за лайк!",
          },
        },
      },
      books: {
        create_request: {
          success:
            "Ура, теперь книга доступна для скачивания! Нажмите на кнопку и книга скачается",
        },
      },
      google_auth: {
        one_tap: {
          error:
            "Мы не смогли получить данные из Google! Проверьте настройки вашего аккаунта Google!",
          success:
            "You have signed in successfully. You can start learning now.",
        },
      },
      languages: {
        lessons: {
          show: {
            lesson_not_found: "Lesson not found. Please, try another lesson.",
          },
        },
        show: {
          empty_language_current_version:
            "Language is under development, please check back later or try any other available language",
          language_in_development_html:
            'Language %{language} is under development.\nYou can help send a pull request with new lessons or\nsupplement those that already exist <a href="%{link_to_repo}" target="_blank" rel="noopener">%{link_to_repo}</a>.\nWe have prepared guidelines for writing <a href="%{link_to_recommendations}" target="_blank" rel="noopener">at the link</a>\n',
          warning:
            "The Language has lessons only in russian. Switch locale if you can read it :)",
        },
        success: {
          error: "В этом курсе есть уроки, которые вы не завершили",
        },
      },
      leads: {
        create: {
          error: "Проверьте ошибки в форме",
          success:
            "Заявка отправлена! Свяжемся с вами в течение одного-двух рабочих дней. Или напишите нам в <a>телеграм</a> чтобы получить помощь быстрее",
        },
      },
      passwords: {
        update: {
          success:
            "Password has been changed. Please, login with your new password.",
        },
      },
      remind_passwords: {
        create: {
          error: "В форме есть ошибки",
          success: "A password recovery instruction was sent to your email.",
        },
      },
      sessions: {
        create: {
          error: "В форме есть ошибки",
          success:
            "You have signed in successfully. You can start learning now.",
        },
      },
      surveys: {
        show: {
          success: "Вы уже отвечали на этот опрос. Возвращаемся обратно",
        },
      },
      users: {
        create: {
          error: "Упс, кажется в форме есть ошибки",
          success:
            "You have signed in successfully. You can start learning now.",
        },
      },
    },
    helpers: {
      continue: "Continue",
      crud: {
        add: "Add",
        adding: "Addmin",
        editing: "Editing",
        list: "List",
        remove: "Удалить",
      },
      read: "Read",
      reset: "Reset",
      select: {
        prompt: "Please select",
      },
      send: "Отправить",
      submit: {
        create: "Create",
        q: {
          create: "Search",
        },
        remind_password_form: {
          create: "Reset password",
        },
        replace: "Change",
        save: "Save",
        submit: "Save %{model}",
        update: "Update",
        user_sign_in_form: {
          create: "Sign In",
        },
        user_sign_up_form: {
          create: "Create an Account",
        },
      },
    },
    home: {
      index: {
        all_blog_posts: "All Blog Posts",
        all_reviews: "All Reviews",
        blog_posts: "Posts",
        categories: "Categories",
        consultation: "Нужна помощь? Оставьте заявку, мы поможем",
        hero: {
          ai_count: "AI Assistant",
          ai_count_description: "Hints and code analysis",
          community_count: "Community",
          community_count_description: "Over 8,000 people",
          courses_count: "Text Courses",
          courses_count_description: "theory, exercises, teacher's solutions",
          fastest_way_to_start_coding:
            "With practice in the trainer and an assistant powered by ChatGPT",
          free_programming_courses: "Free Programming Courses from Scratch",
          how_coding_works: "How Learning Works on Code Basics",
          learn: "Learn",
          learn_xs: "Learn to Code",
          source_code: "Source Code →",
          try: "Courses →",
        },
        hexlet: "Hexlet",
        in_development: "Coming soon",
        in_development_description:
          "Courses below are not in production at the moment. However, you can contribute by sending a pull request to Github repository.",
        join: "Sign up and start learning. For free. Forever",
        meta: {
          description:
            "Learn the Code Basics of HTML, CSS and JavaScript with Interactive Coding Environment right in your Browser! Perfect for Beginners & 100% for Free | Join our Online Programming Courses\n",
        },
        reviews: "Reviews",
        sign_up: "Sign Up",
        subheader: "For those who start from scratch. From the creators of ",
        title: "Free Online Programming Courses: HTML, CSS, JavaScript\n",
      },
      languages: {
        courses: "Courses",
      },
      reviews: {
        course_html: "Course %{link}",
        reviews: "Reviews",
      },
      sitemap: {
        home: "Home",
        title: "Sitemap",
      },
      stats: {
        trusted_by:
          "Присоединись к %{count} изующих программирование с помощью Code Basics\n",
      },
    },
    i18n: {
      plural: {
        keys: ["one", "other", "many", "other"],
      },
    },
    js: {
      hour_zero: "",
      hour_one: "",
      hour_two: "",
      hour_few: "",
      hour_many: "",
      hour_other: "",
    },
    language_categories: {
      index: {
        header: "Course categories",
        link: "Go to",
        meta: {
          description:
            "Browse free programming courses on Code-basics.com.\nFind categories tailored to your learning needs, whether you're a beginner or looking to deepen your expertise.\nStart coding today with interactive and easy-to-follow lessons.\n",
        },
      },
      show: {
        header: "Courses in the %{name} category",
        meta: {
          description:
            "Choose a suitable free course in the %{name} category and learn for free right in your browser with hands-on practice and help from an AI assistant",
        },
      },
    },
    languages: {
      lessons: {
        show: {
          breadcrumb: "breadcrumb",
          chat: {
            community: "Живые люди",
            disabled_html:
              'Чат временно отключен, так как вы достигли суточного лимита. А наше <a href="https://t.me/HexletLearningBot" traget="_blank">телеграм-сообещство</a> работает круглосуточно, подключайтесь :)',
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
              question: "In my environment the code works, but not here 🤨",
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
            {
              question: "Прочитал урок — ничего не понятно 🙄",
              answer:
                "Создавать обучающие материалы, понятные для всех без исключения, довольно сложно. Мы очень стараемся, но всегда есть что улучшать. Если вы встретили материал, который вам непонятен, опишите проблему в обратной связи нашего [сообщества](https://t.me/HexletLearningBot)\n",
            },
          ],
          confirm:
            "Вы хотите сбросить прогресс упражнения. Текущая версия кода не сохранится — надеемся, вы уже скопировали его. Продолжаем сброс?",
          controls: {
            body: "Reset Progress\n",
            header: "Help",
            run: "Run",
          },
          definitions: "Definitions",
          discuss: "AI Assistent",
          editor: "Editor",
          finish: "Завершить",
          if_stuck_html:
            "Если вы столкнулись с трудностями и не знаете, что делать, задайте вопрос в нашем большом и дружном <a>телеграм-сообществе</a>\n",
          instructions: "Instructions",
          issues:
            "Found a bug? Have something to add? Pull requests are welcome!",
          lesson: "Lesson",
          navigation: "Navigation",
          next: "Next →",
          only_for_signed_in_users: "Exercise available only for signed users.",
          output: "Output",
          please_sign_in:
            "Please sign in with your GitHub account, this is necessary to track the progress of the lessons. If you do not have an account yet, now is the time to create an account on GitHub.",
          prev: "← Previous",
          profession_description: "Коммерческий опыт и Трудоустройство",
          separator: " ",
          show_full_version: "Реактор кода доступен в основной версии →",
          sign_in: "Sign In",
          sign_up_for_tracking_progress_html:
            'Be sure to <a href="%{link}" class="text-decoration-none" target="_blank">register</a> to ensure you don\'t lose the results you\'ve achieved\n',
          solution: "Solution",
          tests: "Tests",
          tips: "Tips",
          title: "%{lesson_name} | %{language_name}\n",
          to_home_title: "Home",
        },
      },
      show: {
        about_developer_community:
          "We know how hard it is to start in IT, so we've created a developer community where you're always ready to help. Here you can ask questions, get support, communicate with experienced specialists and get into the profession faster",
        about_learning: "How the training is organised",
        ai_explanation:
          "AI explains topics, tells you how to solve assignments, and helps you at any time - like a personal tutor 24/7",
        ai_without_limits: "AI assistance without limits",
        blog_posts: "Blog Posts",
        breadcrumbs: "Breadcrumbs",
        browser_practice: "Practice in the browser",
        community_image_preview: "Developer community",
        completed_html:
          'Поздравляем! Вы успешно завершили базовый курс на Code Basics. Это первый шаг в мир <b>профессиональной разработки</b>.\nЧто дальше? Посмотрите <a class="link-body-emphasis" href="https://ru.hexlet.io/courses_for_beginners?utm_source=code-basics&utm_medium=referral&utm_campaign=courses_for_beginners&utm_content=course_landing_page" target="_blank">продолжение на Хекслете</a>\n',
        continue: "Continue Learning",
        "convenient format": "Convenient format",
        course_graduates:
          "Join 74,761 students who have successfully completed courses",
        cover_image: "Course cover",
        demo_description:
          "Try a demo lesson without signing up. Practice included",
        demo_start: "Start",
        free_course: "Free course for beginners",
        hexlet_program_link: "Профессия и Трудоустройство",
        join: "Join",
        learning_conveniently:
          "Everything you need to master new topics is on one screen. Theory, practice and live examples go hand in hand. The clear structure helps you to learn in a consistent manner and not to miss important details",
        learning_preview: "Learning preview",
        learning_program: "Learning program",
        lessons: "%{lessons_count} with practice in the browser",
        more_than_support: "More than Support",
        no_registration: "Registration is not required",
        ready: "Are you ready?",
        real_life_challenges:
          "You don't need to install anything - all tasks are performed right in your browser. Built-in code editor, console and automatic tests make learning comfortable. And if something fails, you can always see the teacher's solution",
        registration: "Sign Up",
        registration_description: "Let's get started",
        restart: "Начать заново",
        reviews: "Reviews",
        see_all_courses_in_category: "See all courses in %{name}",
        sign_up: "Sign Up",
        similar_courses: "Similar courses",
        sort_questions: "Sorting out the questions",
        start: "Start Learning",
        start_demo_lesson: "Demo lesson",
        to_home_title: "Home",
        try: "Try It",
        try_without_registration: "Try it without registering",
        updated_at: "updated %{date}",
        without_registration:
          "Start learning right away - first lessons are available without an account. If you like the format, you can register to save your progress",
      },
      success: {
        add_review:
          "Не забудьте оставить отзыв, авторам курсов будет приятно (или не очень хехе). Это можно сделать <a>тут</a>",
        changing_career_html:
          "<a>Карьера</a>: если вы хотите обучиться новой профессии\n",
        choose_your_path: "🔍 Выберите свой путь:\n",
        description:
          "Now you have new knowledge and skills, but that's just the beginning - there are more opportunities ahead. It's time to choose your IT career and move on with your life\n",
        getting_new_skill_html:
          "<a>Навыки и Инструменты</a>: если вы хотите прокачаться в новых для вас технологиях \n",
        header: "Congratulations, you completed the course!",
        home: "Home",
        leave_request:
          "Оставьте заявку на бесплатную консультацию — наш специалист поможет вам определиться и ответит на все вопросы. 👉",
        struggle_choosing: "💬 Не уверены, что выбрать?",
      },
    },
    layouts: {
      admin: {
        application: {
          authentication: "Authentication",
          blog_posts: "Blog Posts",
          dashboard: "Dashboard",
          language_categories: "Categories",
          language_landing_pages: "Landing Pages",
          language_lesson_members: "Прохождение уроков",
          language_lesson_reviews: "Ревью уроков",
          language_lessons: "Уроки",
          languages: "Courses",
          leads: "Лиды",
          menu: "Menu",
          messages: "Conversations with AI",
          reviews: "Reviews",
          survey_answers: "Ответы на опросы",
          survey_scenarios: "Сценарии",
          surveys: "Опросы",
          surveys_analytics: "Аналитика опросов",
          users: "Users",
        },
      },
      root: {
        fallback: {
          description:
            "Oops, something went wrong on the frontend. Try refreshing the page.\nIf that doesn’t help, see if it works in a different browser (sometimes browser extensions cause issues), or even with a different internet connection, like mobile data.\nAnd as a last resort, try clearing your cookies — if you know what that is :)\nWe’ve already received the error and are checking it. If it’s on our side, the developers won’t get their morning coffee.\n",
          header: "Something went wrong",
        },
      },
      shared: {
        contact_method_requesting: {
          description_html:
            "Чтобы мы могли связаться с вами, укажите предпочительный способ связи\n",
          go: "Заполнить&nbsp;→",
        },
        footer: {
          about: "About Platform",
          all_courses: "All Courses",
          authors: "For the author",
          blog: "Blog",
          categories: "Categories",
          codebasics: "Code Basics (by Hexlet)",
          community: "Сообщество",
          cookie_policy: "Cookie Policy",
          courses: "Courses",
          english_links: "English links",
          free_call: "Calls within Russia are toll free",
          hexlet: "Hexlet",
          hexlet_address:
            'TOO "Hexlet", The Republic of Kazakhstan, Almaty, Auezova St., 14A',
          hexlet_blog: "Blog",
          hexlet_number: "BIN 230340043714",
          hexlet_title_html:
            '© <a href="%{link}" class="text-decoration-none text-light" target="_blank" rel="noopener">Hexlet</a>, %{year}',
          "hexlet-cv": "Хекслет.Карьера",
          "hexlet-sicp": "Хекслет.SICP",
          information: "Information",
          language_categories: "Categories",
          matrix: "Matrix",
          privacy: "Privacy",
          recommended_books: "Recommended Books",
          resources: "Resources",
          reviews: "Reviews",
          sitemap: "Sitemap",
          social_networks: "Social",
          source_code: "Source code",
          tos: "Terms of Service",
          useful: "Useful",
        },
        language_menu: {
          in_development: "In development",
        },
        nav: {
          about: "About",
          about_code_basics: "About Code Basics",
          about_platform: "About Platform",
          admin: "Admins",
          blog: "Blog",
          book: "Книга для начинающих",
          business: "Обучение для бизнеса",
          business_description: "Прокачаем ваших сотрудников и дадим аналитику",
          career: "Центр Карьеры",
          career_description: "Доведем до собеседования и научим их проходить",
          cases: "Cases",
          categories: "Categories",
          community: "Community",
          courses: "Courses",
          courses_with_employement: "Курсы с трудоустройством",
          courses_with_employement_description:
            "Обучим программированию, аналитике и тестированию",
          for_authors: "For authors",
          for_teachers: "For school teachers",
          for_teachers_description:
            "Организуем обучение школьников и студентов",
          for_whom: "Школа программирования Хекслет",
          hexly: "Хекслет Колледж",
          hexly_description: "Дадим качественное образование (СПО)",
          my: "My learning",
          profile: "Settings",
          registration: "Sign up",
          reviews: "Reviews",
          sign_in: "Sign in",
          sign_out: "Sign out",
          success_stories: "Why CodeBasics?",
          toggle_navigation: "Toggle navigation",
          toggle_user_menu: "Toggle the display of user menu",
          upskilling: "Повышение квалификации",
          upskilling_description: "Обучим новым навыкам в программировании",
        },
        unsupported_browser_warning: {
          unsupported_browser_warning:
            "Browser is not supported. Some features may be disabled.",
        },
      },
    },
    leads: {
      "new": {
        description:
          "Понимаем, что выбрать направление и план обучения непросто. Заполните форму — мы подскажем, с чего начать и как двигаться дальше, учитывая ваши интересы и цели",
        do_it:
          "📝 <b>Оставьте заявку</b> — наш специалист свяжется с вами, выслушает ваши потребности и поможет подобрать оптимальный план обучения.",
        header: "Консультация со специалистом",
        help_items: [
          "Обсудить ваши цели и интересы, чтобы понять, в каком направлении двигаться дальше",
          "Подобрать подходящие материалы и ресурсы, которые помогут вам разобраться в теме",
          "Поделиться вариантами развития, которые будут актуальны и интересны в вашей ситуации",
        ],
        how_can_we_help: "💡 Чем мы можем помочь?",
        return: "Вернуться",
      },
    },
    links: {
      hexlet:
        "https://hexlet.io/?promo_name=main&promo_position=body&promo_type=link",
      hexlet_author_school_course:
        "https://ru.hexlet.io/courses/author-school-hexlet?promo_name=course-author-school-hexlet&promo_position=body&promo_type=link",
      hexlet_authors_school:
        "https://making.hexlet.io/authors?promo_name=school-of-authors&promo_position=body&promo_type=link",
      hexlet_awesome_text:
        "https://ru.hexlet.io/courses/awesome-text?promo_name=course-awesome-text&promo_position=body&promo_type=link",
      hexlet_awesome_text_structure:
        "https://ru.hexlet.io/courses/awesome-text-structure?promo_name=course-awesome-text-structure&promo_position=body&promo_type=link",
      hexlet_blog:
        "https://hexlet.io/blog?promo_name=blog&promo_position=body&promo_type=link",
      hexlet_courses: "https://hexlet.io/courses",
      hexlet_facebook: "https://www.facebook.com/hexlethq",
      hexlet_frontend:
        "https://hexlet.io/programs/frontend?promo_name=prof-frontend&promo_position=body&promo_type=link",
      hexlet_instagram: "https://instagram.com/hello_hexlet",
      hexlet_java:
        "https://hexlet.io/programs/java?promo_name=prof-java&promo_position=body&promo_type=link",
      hexlet_layout_designer:
        "https://hexlet.io/programs/layout-designer?promo_name=prof-layout-designer&promo_position=body&promo_type=link",
      hexlet_matrix:
        "https://docs.google.com/spreadsheets/d/e/2PACX-1vSr58Xg4DVX2WdmAkv3hE2ITJ3fPeNUmRFe0Ekro53U-ACFrhcUkV8PlUm4ddcn53Uh-5UIezZtZZgc/pubhtml",
      hexlet_php:
        "https://hexlet.io/programs/php?promo_name=prof-php&promo_position=body&promo_type=link",
      hexlet_profession:
        "https://hexlet.io/programs?promo_name=programs&promo_position=body&promo_type=link",
      hexlet_python:
        "https://hexlet.io/programs/python?promo_name=prof-python&promo_position=body&promo_type=link",
      hexlet_rails:
        "https://hexlet.io/programs/rails?promo_name=prof-rails&promo_position=body&promo_type=link",
      hexlet_recommended_books:
        "https://hexlet.io/pages/recommended-books?promo_name=books&promo_position=body&promo_type=link",
      hexlet_success_stories:
        "https://hexlet.io/blog/categories/success?promo_name=blog-success&promo_position=body&promo_type=link",
      hexlet_telegram: "https://ttttt.me/hexlet_ru",
      hexlet_telegram_volunteers: "https://ttttt.me/hexletcommunity/12",
      hexlet_twitter: "https://twitter.com/hexlet_io",
      hexlet_youtube:
        "https://www.youtube.com/channel/UCMGJqXOa4m37IdmDLO-n-eQ",
    },
    models: {
      attributes: {
        base: {
          _destroy: "Удалить",
          answer: "Ответ",
          created_at: "Создание",
          description: "Описание",
          header: "Заголовок",
          name: "Имя",
          question: "Вопрос",
          slug: "Слаг",
          updated_at: "Обновление",
        },
        blog_post: {
          body: "Content",
          cover: "Image",
          creator: "Author",
          state: "State",
          "state/archived": "Archived",
          "state/published": "Published",
        },
        language: {
          category_id: "Category",
          cover: "Cover",
          hexlet_program_landing_page: "Hexlet Program Landing Page",
          learn_as: "Learn as",
          openai_assistant_id: "Assistant ID (OpenAI)",
          progress: "Progress",
          slug: "Slug",
        },
        language_category: {
          description: "Description",
          header: "Header",
          name: "Name",
          slug: "Slug",
        },
        language_category_item: {
          language_landing_page_id: "Лендинг",
        },
        language_landing_page: {
          _destroy: "Destroy",
          answer: "Answer",
          description: "Description",
          footer: "Show in Footer",
          footer_name: "Footer Name",
          header: "Header",
          landing_page_to_redirect_id: "Redirect to Page",
          language_id: "Course",
          listed: "Show in Lists",
          main: "Main",
          meta_description: "Meta Description",
          meta_title: "Meta Title",
          name: "Name",
          order: "Order",
          outcomes_description: "Learning Outcomes (Description)",
          outcomes_header: "Learning Outcomes (Header)",
          outcomes_image: "Learning Outcomes (Image)",
          question: "Question",
          slug: "Slug",
          state: "State",
          used_in_description: "Used In (Description)",
          used_in_header: "Used In (Header)",
        },
        lead: {
          contact_method: "Способ связи",
          contact_value: "Телефон / Имя пользователя",
        },
        qna_item: {
          answer: "Ответ",
          footer: "Выводить в футере",
          footer_name: "Имя в футере",
          listed: "Отображать в списках",
          main: "Основной",
          order: "Order",
          question: "Вопрос",
          slug: "Slug",
        },
        remind_password_form: {
          email: "Email",
        },
        review: {
          body: "Review",
          first_name: "First name",
          language_id: "Course",
          last_name: "Last name",
          locale: "Locale",
          pinned: "Pinned (Shown at the top of the list)",
          state: "State",
          user_id: "User",
        },
        survey: {
          run_always: "Запускать всегда",
          slug: "Слаг",
        },
        survey_scenario: {
          _destroy: "Удалить",
          name: "Название",
          survey_item_id: "Ответ триггер (Начало воронки)",
        },
        user: {
          admin: "Admin?",
          contact_method: "Способ связи",
          "contact_method/values": {
            phone: "Телефон",
            telegram: "Telegram",
            whatsapp: "WhatsApp",
          },
          contact_value: "Номер для связи (или имя пользователя)",
          email: "Email",
          first_name: "First name",
          last_name: "Last name",
          nickname: "Nickname",
          password: "Password",
        },
        user_password_form: {
          password: "Password",
        },
        user_sign_in_form: {
          email: "Email",
          password: "Password",
        },
        user_sign_up_form: {
          email: "Email",
          first_name: "First Name",
          password: "Password",
        },
      },
    },
    my: {
      show: {
        add_review:
          "Мы будем очень признательны, если вы оставите отзывы по пройденным курсам. Это можно сделать <a>тут</a>",
        finished: "Finished Courses",
        started: "Started Courses",
      },
    },
    number: {
      currency: {
        format: {
          delimiter: ",",
          format: "%u%n",
          precision: 2,
          separator: ".",
          significant: false,
          strip_insignificant_zeros: false,
          unit: "$",
        },
      },
      format: {
        delimiter: ",",
        precision: 3,
        round_mode: "default",
        separator: ".",
        significant: false,
        strip_insignificant_zeros: false,
      },
      human: {
        decimal_units: {
          format: "%n %u",
          units: {
            billion: {
              few: "",
              many: "",
              one: "",
              other: "",
            },
            million: {
              few: "",
              many: "",
              one: "",
              other: "",
            },
            quadrillion: {
              few: "",
              many: "",
              one: "",
              other: "",
            },
            thousand: {
              few: "",
              many: "",
              one: "",
              other: "",
            },
            trillion: {
              few: "",
              many: "",
              one: "",
              other: "",
            },
            unit: "",
          },
        },
        format: {
          delimiter: "",
          precision: 3,
          significant: true,
          strip_insignificant_zeros: true,
        },
        storage_units: {
          format: "%n %u",
          units: {
            byte: {
              few: "байта",
              many: "байт",
              one: "Byte",
              other: "Bytes",
            },
            eb: "EB",
            gb: "GB",
            kb: "KB",
            mb: "MB",
            pb: "PB",
            tb: "TB",
          },
        },
      },
      percentage: {
        format: {
          delimiter: "",
          format: "%n%",
        },
      },
      precision: {
        format: {
          delimiter: "",
        },
      },
    },
    pages: {
      account: {
        profiles: {
          edit: {
            delete: "Delete account",
            title: "Profile editing",
          },
        },
      },
      admin: {
        blog_posts: {
          edit: {
            header: "Edit Blog Post",
          },
          index: {
            header: "Blog Posts",
          },
          "new": {
            header: "New Blog Post",
          },
        },
        home: {
          index: {
            dashboard: "Dashboard",
          },
        },
        language_categories: {
          edit: {
            header: "Редактирование %{id}",
          },
          form: {
            items: "Items",
            main: "Main",
            qna_items: "Q&A Items",
          },
          index: {
            header: "Категории",
          },
          "new": {
            header: "Новая категория",
          },
        },
        language_landing_pages: {
          edit: {
            header: "Редактирование %{id}",
          },
          index: {
            header: "Лендинги",
          },
          "new": {
            header: "Новый лендинг",
          },
        },
        language_lesson_members: {
          index: {
            header: "Прохождение уроков",
          },
        },
        language_lesson_reviews: {
          index: {
            data: "Data",
            header: "Ревью уроков",
          },
        },
        language_lessons: {
          index: {
            header: "Уроки",
          },
        },
        languages: {
          edit: {
            header: "Edit Language",
            id: "ID",
            load_new_version: "Load new version",
            result: "Result",
            versions: "Versions",
          },
          index: {
            header: "Languages",
          },
          "new": {
            header: "New Language",
          },
        },
        leads: {
          index: {
            answers: "Answers",
            courses: "Courses",
            data: "Data",
            header: "Лиды",
          },
        },
        management: {
          users: {
            edit: {
              count: "Count",
              header: "Edit User",
              language: "Language",
              progress: "Progress",
            },
            index: {
              header: "Users",
              search_by_email: "Search by email",
            },
          },
        },
        messages: {
          index: {
            header: "Разговоры с AI",
          },
        },
        reviews: {
          edit: {
            header: "Edit Review",
          },
          index: {
            header: "Reviews",
          },
          "new": {
            header: "New Review",
          },
        },
        survey_answers: {
          index: {
            header: "Ответы на опросы",
          },
        },
        survey_scenarios: {
          edit: {
            header: "Редактирование сценария",
          },
          index: {
            header: "Сценарии",
          },
          "new": {
            header: "Новый сценарий",
          },
        },
        surveys: {
          edit: {
            header: "Редактирование опроса",
          },
          index: {
            header: "Опросы",
          },
          "new": {
            header: "Новый опрос",
          },
        },
      },
      blog_posts: {
        index: {
          header: "Blog",
        },
        show: {
          discuss: "Обсуждайте, помогайте и делитесь опытом",
          join_community: "Присоединяйтесь к сообществу Хекслет",
          link: "Перейти в сообщество",
        },
      },
      books: {
        show: {
          chapter: "Глава %{number}",
          description:
            "С нуля до трудоустройства. Этот учебник — ваш путеводитель в профессию программиста. Он создан для тех, кто хочет начать с нуля и шаг за шагом дойти до трудоустройства в одной из самых востребованных и перспективных областей\n",
          download: "Скачать книгу",
          freebook: "Бесплатная электронная книга",
          header: "Книга: Профессия программист. С нуля до трудоустройства\n",
          request: "Запросить книгу",
          toc: "Содержание",
        },
      },
      cases: {
        for_teachers: {
          description:
            "Бесплатные курсы с тренажером прямо в браузере. Много практики и ИИ-ассистент",
          early_career_guidance: "Стимулируем раннюю профориентацию:",
          early_career_guidance_list: [
            "Показываем программирование изнутри, чтобы принять взвешенное решение об изучении будущей профессии",
            "Даем возможность быстро погрузиться в язык программирования и понять, нравится ли он ученику",
            "Помогаем увидеть базовые примеры реальных задач, чтобы присмотреться к профессиональной рутине",
          ],
          header:
            "Эффективное обучение школьников и студентов программированию",
          how_to_learn_programming:
            "Как пользователи изучают программирование в CodeBasics?",
          how_to_learn_programming_cards: [
            {
              title:
                "Интеграция без интеграции: чтобы внедрить CodeBasics в образовательный процесс, потребуется ПК с доступом в интернет",
              subtitle:
                "Нужно открыть браузер, зарегистрироваться и начать проходить интересующий курс",
              img: "integration-icon",
            },
            {
              title:
                "Теория, подкрепленная практикой: на платформе уже сейчас доступны тренажеры по PHP, Java, JS, Python, Ruby, HTML, CSS, Racket, Elixir и Go",
              subtitle: "Можно проходить несколько курсов одновременно",
              img: "practice-icon",
            },
            {
              title:
                "Интерактивный формат: система автоматически проверяет выполненные задания, показывает ошибки и правильное решение, а также детальный вывод",
              subtitle:
                "Если что-то не получается, можно воспользоваться Решением учителя",
              img: "interactive-format-icon",
            },
            {
              title:
                "Пользователям доступен Виртуальный наставник на основе Искусственного интеллекта",
              subtitle: "",
              img: "discussion-icon",
            },
          ],
          integrate_into_education:
            "Зачем интегрировать CodeBasics в обучение?",
          integrate_now: "Интегрируйте CodeBasics в учебную программу сейчас!",
          interactive_approach:
            "Ориентируемся на практику и интерактивный подход:",
          interactive_approach_list: [
            "Сопровождаем теоретические материалы курса практическими заданиями прямо в браузере",
            "Даем мгновенную обратную связь: показываем правильное решение, даем вывод по каждому тесту",
            "Подключаем ИИ-наставника, который объясняет ученикам теорию или суть задания",
          ],
          lay_programming_foundations:
            "CodeBasics закладывает основы программирования в интерактивной форме, за что получил 3-е место в номинации «Лучший внеплатформенный онлайн-курс» в 2020 году на Международном конкурсе открытых онлайн-курсов EdCrunch Award OOC",
          open_browser_and_sign_up:
            "Просто откройте браузер и зарегистрируйтесь",
          select_course: "Выбрать курс",
          sign_up: "Зарегистрироваться",
          try: "Попробовать",
        },
        index: {
          for_teachers: "Для учителей и преподавателей",
          link: "Перейти",
        },
      },
      home: {
        index: {
          blog_posts: "Posts",
          consultation: "Нужна помощь? Оставьте заявку, мы поможем",
          hero: {
            ai_count: "AI Assistant",
            ai_count_description: "Hints and code analysis",
            community_count: "Community",
            community_count_description: "Over 8,000 people",
            courses_count: "Text Courses",
            courses_count_description: "theory, exercises, teacher's solutions",
            fastest_way_to_start_coding:
              "With practice in the trainer and an assistant powered by ChatGPT",
            free_programming_courses: "Free Programming Courses from Scratch",
            source_code: "Source Code →",
            try: "Courses →",
          },
          join: "Sign up and start learning. For free. Forever",
          reviews: "Reviews",
          sign_up: "Sign Up",
        },
        languages: {
          courses: "Courses",
        },
        sitemap: {
          home: "Home",
        },
      },
      language_categories: {
        index: {
          header: "Course categories",
          link: "Go to",
        },
      },
      languages: {
        lessons: {
          show: {
            chat: {
              community: "Живые люди",
              disabled_html:
                'Чат временно отключен, так как вы достигли суточного лимита. А наше <a href="https://t.me/HexletLearningBot" traget="_blank">телеграм-сообещство</a> работает круглосуточно, подключайтесь :)',
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
                question: "In my environment the code works, but not here 🤨",
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
              {
                question: "Прочитал урок — ничего не понятно 🙄",
                answer:
                  "Создавать обучающие материалы, понятные для всех без исключения, довольно сложно. Мы очень стараемся, но всегда есть что улучшать. Если вы встретили материал, который вам непонятен, опишите проблему в обратной связи нашего [сообщества](https://t.me/HexletLearningBot)\n",
              },
            ],
            controls: {
              body: "Reset Progress\n",
              header: "Help",
              run: "Run",
            },
            discuss: "AI Assistent",
            editor: "Editor",
            finish: "Завершить",
            if_stuck_html:
              "Если вы столкнулись с трудностями и не знаете, что делать, задайте вопрос в нашем большом и дружном <a>телеграм-сообществе</a>\n",
            instructions: "Instructions",
            issues:
              "Found a bug? Have something to add? Pull requests are welcome!",
            lesson: "Lesson",
            navigation: "Navigation",
            next: "Next →",
            output: "Output",
            prev: "← Previous",
            profession_description: "Коммерческий опыт и Трудоустройство",
            solution: "Solution",
            tests: "Tests",
            tips: "Tips",
          },
        },
        show: {
          about_developer_community:
            "We know how hard it is to start in IT, so we've created a developer community where you're always ready to help. Here you can ask questions, get support, communicate with experienced specialists and get into the profession faster",
          about_learning: "How the training is organised",
          ai_explanation:
            "AI explains topics, tells you how to solve assignments, and helps you at any time - like a personal tutor 24/7",
          ai_without_limits: "AI assistance without limits",
          browser_practice: "Practice in the browser",
          continue: "Continue Learning",
          "convenient format": "Convenient format",
          course_graduates:
            "Join 74,761 students who have successfully completed courses",
          cover_image: "Course cover",
          demo_description:
            "Try a demo lesson without signing up. Practice included",
          demo_start: "Start",
          free_course: "Free course for beginners",
          hexlet_program_link: "Профессия и Трудоустройство",
          join: "Join",
          learning_conveniently:
            "Everything you need to master new topics is on one screen. Theory, practice and live examples go hand in hand. The clear structure helps you to learn in a consistent manner and not to miss important details",
          learning_preview: "Learning preview",
          learning_program: "Learning program",
          more_than_support: "More than Support",
          real_life_challenges:
            "You don't need to install anything - all tasks are performed right in your browser. Built-in code editor, console and automatic tests make learning comfortable. And if something fails, you can always see the teacher's solution",
          restart: "Начать заново",
          sort_questions: "Sorting out the questions",
          start: "Start Learning",
          try: "Try It",
          updated_at: "updated %{date}",
        },
        success: {
          add_review:
            "Не забудьте оставить отзыв, авторам курсов будет приятно (или не очень хехе). Это можно сделать <a>тут</a>",
          changing_career_html:
            "<a>Карьера</a>: если вы хотите обучиться новой профессии\n",
          choose_your_path: "🔍 Выберите свой путь:\n",
          description:
            "Now you have new knowledge and skills, but that's just the beginning - there are more opportunities ahead. It's time to choose your IT career and move on with your life\n",
          getting_new_skill_html:
            "<a>Навыки и Инструменты</a>: если вы хотите прокачаться в новых для вас технологиях \n",
          header: "Congratulations, you completed the course!",
          leave_request:
            "Оставьте заявку на бесплатную консультацию — наш специалист поможет вам определиться и ответит на все вопросы. 👉",
          struggle_choosing: "💬 Не уверены, что выбрать?",
        },
      },
      leads: {
        "new": {
          description:
            "Понимаем, что выбрать направление и план обучения непросто. Заполните форму — мы подскажем, с чего начать и как двигаться дальше, учитывая ваши интересы и цели",
          do_it:
            "📝 <b>Оставьте заявку</b> — наш специалист свяжется с вами, выслушает ваши потребности и поможет подобрать оптимальный план обучения.",
          header: "Консультация со специалистом",
          help_items: [
            "Обсудить ваши цели и интересы, чтобы понять, в каком направлении двигаться дальше",
            "Подобрать подходящие материалы и ресурсы, которые помогут вам разобраться в теме",
            "Поделиться вариантами развития, которые будут актуальны и интересны в вашей ситуации",
          ],
          how_can_we_help: "💡 Чем мы можем помочь?",
          return: "Вернуться",
        },
      },
      my: {
        show: {
          add_review:
            "Мы будем очень признательны, если вы оставите отзывы по пройденным курсам. Это можно сделать <a>тут</a>",
          finished: "Finished Courses",
          started: "Started Courses",
        },
      },
      parts: {
        about: {
          meta: {
            description:
              "Code Basics quality interactive lessons on programming basics for free, for everyone.",
          },
          title: "About",
        },
        authors: {
          meta: {
            description:
              "Code Basics is open source. Everyone can be an author. Here is some simple steps to start.\n",
          },
          title: "How to become an author?",
        },
        cookie_policy: {
          meta: {
            description: "Cookie policy",
          },
          title: "Cookie Policy",
        },
        privacy: {
          meta: {
            description: "Code Basics Privacy Policy",
          },
          title: "Privacy Policy",
        },
        tos: {
          meta: {
            description: "Terms of Service of the Code Basics Platform",
          },
          title: "Terms of Service",
        },
      },
      passwords: {
        edit: {
          title: "Change password",
        },
      },
      remind_passwords: {
        "new": {
          title: "Remind password",
        },
      },
      reviews: {
        index: {
          add_review:
            "Попробовали наши курсы и хотите оставить отзыв? Это можно сделать <a>тут</a>",
          header: "Reviews",
        },
      },
      sessions: {
        "new": {
          dont_have_account: "New to Code-Basics?",
          forgot_password: "Forgot password?",
          register: "Sign up",
          reset_password: "Reset Password",
          title: "Sign In",
        },
      },
      users: {
        "new": {
          confirmation_html:
            "By clicking Sign up, you agree to our <a>service conditions</a>",
          demo_html:
            "<b>Зарегистрируйтесь</b> для продолжения бесплатного обучения. Так мы сможем сохранить ваш прогресс и открыть доступ к нашему ассистенту",
          have_account: "Already have an account?",
          sign_in: "Sign in",
          sign_up: "Sign up",
        },
      },
    },
    passwords: {
      edit: {
        meta_description: "Set a new password on Code Basics",
        new_password: "New password",
        title: "Change password",
      },
    },
    ransack: {
      all: "all",
      and: "and",
      any: "any",
      asc: "ascending",
      attribute: "attribute",
      combinator: "combinator",
      condition: "condition",
      desc: "descending",
      or: "or",
      predicate: "predicate",
      predicates: {
        blank: "is blank",
        cont: "contains",
        cont_all: "contains all",
        cont_any: "contains any",
        does_not_match: "doesn't match",
        does_not_match_all: "doesn't match all",
        does_not_match_any: "doesn't match any",
        end: "ends with",
        end_all: "ends with all",
        end_any: "ends with any",
        eq: "equals",
        eq_all: "equals all",
        eq_any: "equals any",
        false: "is false",
        gt: "greater than",
        gt_all: "greater than all",
        gt_any: "greater than any",
        gteq: "greater than or equal to",
        gteq_all: "greater than or equal to all",
        gteq_any: "greater than or equal to any",
        in: "in",
        in_all: "in all",
        in_any: "in any",
        lt: "less than",
        lt_all: "less than all",
        lt_any: "less than any",
        lteq: "less than or equal to",
        lteq_all: "less than or equal to all",
        lteq_any: "less than or equal to any",
        matches: "matches",
        matches_all: "matches all",
        matches_any: "matches any",
        not_cont: "doesn't contain",
        not_cont_all: "doesn't contain all",
        not_cont_any: "doesn't contain any",
        not_end: "doesn't end with",
        not_end_all: "doesn't end with all",
        not_end_any: "doesn't end with any",
        not_eq: "not equal to",
        not_eq_all: "not equal to all",
        not_eq_any: "not equal to any",
        not_in: "not in",
        not_in_all: "not in all",
        not_in_any: "not in any",
        not_null: "is not null",
        not_start: "doesn't start with",
        not_start_all: "doesn't start with all",
        not_start_any: "doesn't start with any",
        null: "is null",
        present: "is present",
        start: "starts with",
        start_all: "starts with all",
        start_any: "starts with any",
        true: "is true",
      },
      search: "search",
      sort: "sort",
      value: "value",
    },
    remind_passwords: {
      "new": {
        dont_have_account: "New to Code-Basics?",
        forgot_password: "Forgot password?",
        login: "Sign in",
        meta_description: "Restore Code Basics password",
        register: "Sign up",
        title: "Remind password",
        trying_to_login: "Trying to login?",
      },
    },
    reviews: {
      index: {
        add_review:
          "Попробовали наши курсы и хотите оставить отзыв? Это можно сделать <a>тут</a>",
        course: "Course %{language}",
        empty: "Looks like there is noting yet",
        header: "Reviews",
        meta: {
          description:
            "Честные отзывы наших студентов о курсах, процессе обучения и платформе\n",
        },
        read_more: "Read more",
      },
    },
    sessions: {
      "new": {
        dont_have_account: "New to Code-Basics?",
        email: "Email",
        forgot_password: "Forgot password?",
        meta: {
          description: "Войти на Code Basics",
        },
        password: "Password",
        register: "Sign up",
        reset_password: "Reset Password",
        sign_in: "Sign in",
        sign_in_with_github: "Log in with GitHub",
        title: "Sign In",
      },
    },
    shared: {
      languages: {
        course_finished: "Course finished!",
      },
      via_social_networks: {
        via_social_networks: "With GitHub",
      },
    },
    support: {
      array: {
        last_word_connector: ", and ",
        two_words_connector: " and ",
        words_connector: ", ",
      },
    },
    time: {
      am: "am",
      formats: {
        date_time: "%Y-%m-%d %H:%M",
        default: "%a, %d %b %Y %H:%M:%S %z",
        long: "%B %d, %Y %H:%M",
        short: "%d %b %H:%M",
      },
      pm: "pm",
    },
    user_mailer: {
      reset_password: {
        change_password: "Сменить пароль",
        link: "Ссылка для изменения пароля на сайте code-basics.com",
        requested_to_change_password:
          "Вы запросили ссылку для изменения вашего пароля на code-basics.com. Чтобы изменить пароль – перейдите по этой ссылке или нажмите на кнопку ниже.",
        subject: "Смена пароля",
        you_dont_want_to_change:
          "Если вы не запрашивали изменения пароля, просто проигнорируйте это письмо. Ваш пароль не изменится пока вы не перейдете по ссылке и не введете новый пароль.",
      },
    },
    users: {
      "new": {
        confirmation_html:
          "By clicking Sign up, you agree to our <a>service conditions</a>",
        demo_html:
          "<b>Зарегистрируйтесь</b> для продолжения бесплатного обучения. Так мы сможем сохранить ваш прогресс и открыть доступ к нашему ассистенту",
        have_account: "Already have an account?",
        meta: {
          description: "Зарегистрироваться на Code Basics",
        },
        sign_in: "Sign in",
        sign_up: "Sign up",
        title: "Registration",
      },
    },
  },
} as const;
