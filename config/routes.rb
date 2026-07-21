# typed: false

Rails.application.routes.draw do
  # Redirects
  draw :redirects

  mount RailsEventStore::Browser => "/res" if Rails.env.development?
  # mount ActionCable.server => "/cable"

  # Feature-flag admin UI, gated to admins by the signed session cookie.
  constraints(AdminConstraint) do
    mount Flipper::UI.app(Flipper) => "/admin/flipper", as: :flipper_ui
  end
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
    namespace :api, defaults: { format: :json } do
      # resource :chatkit, only: [ :create ]
      resources :feeds, only: [] do
        collection do
          get :yandex_courses
        end
      end
      resources :blog_posts, only: [ :show ] do
        member do
          get :next
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

      # use resource
      get "/map", to: "home#sitemap"

      resource :my, only: [ :show ]

      resources :leads

      constraints(RuSuffixConstraint) do
        resource :book do
          member do
            post :create_request
            get :download
          end
        end
      end

      resources :cases, only: %i[index] do
        collection do
          get :for_teachers
        end
      end

      resources :pages, only: %i[show]
      resources :blog_posts, only: %i[index show] do
        scope module: :blog_posts do
          resources :likes, only: [ :create ]
        end
      end
      resources :reviews, only: %i[index]
      resources :language_categories, only: %i[index show]
      resource :session, only: %i[new create destroy]
      resource :phone_auth, only: %i[new create] do
        get :verify
        post :confirm
      end
      resource :passkey_session, only: %i[new create]
      resource :locale, only: [] do
        member do
          get :switch
        end
      end
      resources :users, only: %i[new create]
      resources :magic_links, only: %i[new create show]
      resource :remind_password, only: %i[new create]
      get "/password/:token/edit", to: "passwords#edit", as: :edit_password
      patch "/password/:token", to: "passwords#update", as: :password

      resources :languages, only: %i[index show] do
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
          resources :language_categories, only: [] do
            resources :qna_items, only: %i[index create update destroy], module: :language_categories
          end
          resources :language_landing_pages, only: [] do
            resources :qna_items, only: %i[index create update destroy], module: :language_landing_pages
          end
        end

        resources :reviews
        resources :banners
        resources :leads, only: [ :index ]
        resources :messages, only: [ :index ]
        resources :language_lesson_members, only: [ :index ]
        resources :language_lesson_reviews, only: [ :index ]
        resources :language_categories
        resources :language_lessons, only: [ :index ] do
          member do
            post :review
          end
        end
        resources :blog_posts do
          member do
            post :related_courses
          end
        end
        resources :languages, only: %i[index new edit update create] do
          member do
            post :review
          end
          scope module: :languages do
            resources :versions, only: %i[create]
          end
        end
        resources :language_landing_pages

        namespace :management do
          resources :users, only: %i[index edit update]
          resources :roles
          resources :role_permissions, only: %i[show update]
          resources :staff_members
        end
      end

      namespace :account do
        resource :profile, only: %i[edit update destroy]
        resources :passkeys, only: %i[new create destroy]
      end
    end
  end
end
