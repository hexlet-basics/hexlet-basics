# frozen_string_literal: true

Rails.application.routes.draw do
  require 'sidekiq/web'
  mount Sidekiq::Web => '/sidekiq', constraints: AdminConstraint.new

  namespace :api do
    resources :lessons, only: [] do
      member do
        post :check
      end
    end
  end

  scope module: :web do
    root 'home#index'
    get 'auth/:provider', to: 'auth#request', as: :auth_request
    get 'auth/:provider/callback', to: 'auth#callback', as: :callback_auth

    resource :session, only: %i[new create destroy]
    resources :users, only: %i[new create]

    resources :languages, only: [:show] do
      scope module: :languages do
        resources :lessons, only: [:show] do
          get :next_lesson, on: :member
          get :prev_lesson, on: :member
        end
      end
    end
    resources :pages, only: [] do
      get :about, on: :collection
      get :tos, on: :collection
      get :privacy, on: :collection
    end

    namespace :admin do
      root 'home#index'

      resources :languages, only: %i[index new edit update create] do
        scope module: :languages do
          resources :versions, only: %i[index create]
        end
      end
    end
  end
end
