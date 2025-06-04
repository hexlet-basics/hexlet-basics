Rails.application.routes.draw do
  get "*path.amp", to: redirect("/%{path}", status: 301)

  # Redirects (Old routes/SEO)
  get "/:locale/languages/:lang/modules/:module/lessons/:lesson", to: redirect("/%{locale}/languages/%{lang}/lessons/%{lesson}")

  get "/ru/languages/pre-course-python", to: redirect("/ru/languages/python")
  get "/ru/languages/pre-course-python/*", to: redirect("/ru/languages/python")

  get "/ru/languages/pre-course-java", to: redirect("/ru/languages/java")
  get "/ru/languages/pre-course-java/*", to: redirect("/ru/languages/java")

  get "/ru/languages/pre-course-javascript", to: redirect("/ru/languages/javascript")
  get "/ru/languages/pre-course-javascript/*", to: redirect("/ru/languages/javascript")


  mount RailsEventStore::Browser => "/res" if Rails.env.development?
  mount ActionCable.server => "/cable"
  # Reveal health status on /up that returns 200 if the app boots with no exceptions, otherwise 500.
  # Can be used by load balancers and uptime monitors to verify that the app is live.
  get "up" => "rails/health#show", as: :rails_health_check

  scope module: :web do
    # post "/auth/:provider", to: "auth#request", as: :auth_request
    # get "/auth/:provider/callback", to: "auth#callback", as: :callback_auth
    # post "/google/callback", to: "google_auth#one_tap", as: :google_onetap_callback

    # match "/403", to: "errors#forbidden", via: :all
    # match "/404", to: "errors#not_found", via: :all
    # match "/500", to: "errors#server_error", via: :all
    match "/:code",
      to: "errors#show",
      via: :all,
      constraints: {
        code: Regexp.new(
          Rack::Utils::HTTP_STATUS_CODES.keys.join("|")
        )
      }
  end

  scope "(:suffix)", suffix: /es|ru/ do
    namespace :ai do
      resources :lessons, only: [] do
        scope module: :lessons do
          resources :messages, only: [ :index, :create ]
        end
      end
    end
    namespace :api do
      resources :feeds, only: [] do
        collection do
          get :yandex_courses
        end
      end
      resources :lessons, only: [] do
        member do
          post :check
        end
      end
    end

    scope module: :web, format: false do
      root "home#index"

      get "/robots.:format", to: "home#robots", as: :robots
      # get "/map", to: "home#sitemap"

      resource :my, only: [ :show ]

      resources :leads
      resources :scenarios, only: [ :show ] do
        scope module: :scenarios do
          resources :surveys, only: [ :show ] do
            scope module: :surveys do
              resources :answers, only: [ :create ]
            end
          end
        end
      end

      resource :book do
        member do
          post :create_request
          get :download
        end
      end

      resources :cases, only: %i[index] do
        collection do
          get :for_teachers
        end
      end

      resources :pages, only: %i[show]
      resources :blog_posts, only: %i[index show]
      resources :reviews, only: %i[index]
      resources :language_categories, only: %i[index show]
      resource :session, only: %i[new create destroy]
      resource :locale, only: [] do
        member do
          get :switch
        end
      end
      resources :users, only: %i[new create]
      resource :remind_password, only: %i[new create]
      resource :password, only: %i[edit update]

      resources :languages, only: [ :show ] do
        member do
          get :success
        end
        scope module: :languages do
          resources :lessons, only: [ :show ]
          # do
          #   get :next_lesson, on: :member
          #   get :prev_lesson, on: :member
          # end
        end
      end

      namespace :admin do
        root "home#index"

        namespace :api do
          resources :users do
            collection do
              get :search
            end
          end
        end

        resources :surveys
        resources :survey_scenarios
        resources :survey_answers, only: [ :index ]

        resources :reviews
        resources :leads, only: [ :index ]
        resources :messages, only: [ :index ]
        resources :language_lesson_members, only: [ :index ]
        resources :language_categories
        resources :blog_posts
        resources :languages, only: %i[index new edit update create] do
          scope module: :languages do
            resources :versions, only: %i[create]
          end
        end
        resources :language_landing_pages

        namespace :management do
          resources :users, only: %i[index edit update]
        end
      end

      namespace :account do
        resource :profile, only: %i[edit update destroy]
      end
    end
  end
end
