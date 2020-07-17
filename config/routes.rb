# frozen_string_literal: true

Rails.application.routes.draw do
  require 'sidekiq/web'
  mount Sidekiq::Web => '/sidekiq', constraints: AdminConstraint.new

  namespace :api do
    resources :lessons, only: [] do
      scope module: :lessons do
        resource :check, only: [:create]
      end
    end
  end

  scope module: :web do
    root 'home#index'

    resource :session, only: %i[new create destroy]
    resources :users, only: %i[new create]

    resources :languages, only: [:show] do
      scope module: :languages do
        resources :lessons, only: [:show]
      end
    end

    namespace :admin do
      root 'home#index'

      resources :languages, only: %i[index new create] do
        scope module: :languages do
          resources :versions, only: %i[index create]
        end
      end
    end
  end
end
