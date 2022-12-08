# frozen_string_literal: true

Rails.application.routes.draw do
  # require 'sidekiq/web'
  # mount Sidekiq::Web => '/sidekiq', constraints: AdminConstraint.new

  # for kubernets probe
  get '/health', to: proc { |_env| [200, {}, ['it works!']] }

  scope module: :web do
    post '/auth/:provider', to: 'auth#request', as: :auth_request
    get '/auth/:provider/callback', to: 'auth#callback', as: :callback_auth

    match '/403', to: 'errors#forbidden', via: :all
    match '/404', to: 'errors#not_found', via: :all
    match '/500', to: 'errors#server_error', via: :all
  end

  scope '(:locale)', locale: /en|ru/ do
    namespace :api do
      resources :languages, only: %i[index show] do
        scope module: :languages do
          resources :lessons, only: %i[index show]
        end
      end
      resources :lessons, only: [] do
        member do
          post :check
        end
      end
    end

    scope module: :web do
      root 'home#index'

      get '/robots.:format' => 'home#robots', as: :robots
      resources :pages, only: %i[show]
      resources :blog_posts, only: %i[index show]
      resources :reviews, only: %i[index show]
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

      resources :languages, only: [:show] do
        scope module: :languages do
          resources :lessons, only: [:show] do
            get :next_lesson, on: :member
            get :prev_lesson, on: :member
          end
        end
      end

      namespace :admin do
        root 'home#index'

        namespace :api do
          resources :users do
            collection do
              get :search
            end
          end
        end

        resources :reviews
        resources :blog_posts
        resources :languages, only: %i[index new edit update create] do
          scope module: :languages do
            resources :versions, only: %i[index create]
          end
        end

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
